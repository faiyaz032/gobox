package box

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/faiyaz032/gobox/internal/domain"
	"github.com/gorilla/websocket"
)

type connEvent struct {
	fingerprint string
	increment   bool
	responseCh  chan struct{}
}

type shutdownRequest struct {
	fingerprint string
	containerID string
}

type Svc struct {
	repo        Repo
	dockerSvc   DockerSvc
	connEventCh chan connEvent
	shutdownCh  chan shutdownRequest
}

func NewSvc(repo Repo, dockerSvc DockerSvc) *Svc {
	svc := &Svc{
		repo:        repo,
		dockerSvc:   dockerSvc,
		connEventCh: make(chan connEvent),
		shutdownCh:  make(chan shutdownRequest),
	}

	// Start connection manager goroutine
	go svc.manageConnections()

	return svc
}

func (s *Svc) manageConnections() {
	activeConns := make(map[string]int)
	shutdownTimers := make(map[string]*time.Timer)

	for {
		select {
		case event := <-s.connEventCh:
			if event.increment {

				activeConns[event.fingerprint]++

				if timer, exists := shutdownTimers[event.fingerprint]; exists {
					timer.Stop()
					delete(shutdownTimers, event.fingerprint)
					fmt.Printf("[BoxSvc] Cancelled shutdown timer for %s (new connection)\n", event.fingerprint)
				}

				fmt.Printf("[BoxSvc] Connection established (active: %d) for fingerprint: %s\n",
					activeConns[event.fingerprint], event.fingerprint)
			} else {

				activeConns[event.fingerprint]--
				remaining := activeConns[event.fingerprint]

				if remaining <= 0 {
					delete(activeConns, event.fingerprint)
					fmt.Printf("[BoxSvc] Last connection closed for %s, will shutdown in 5 seconds...\n", event.fingerprint)
				} else {
					fmt.Printf("[BoxSvc] Connection closed (%d remaining) for %s\n", remaining, event.fingerprint)
				}
			}

			close(event.responseCh)

		case req := <-s.shutdownCh:

			if activeConns[req.fingerprint] <= 0 {
				fmt.Printf("[BoxSvc] Executing shutdown for container: %s\n", req.containerID)

				if err := s.dockerSvc.StopContainer(context.Background(), req.containerID); err != nil {
					fmt.Printf("[BoxSvc] Error stopping container %s: %v\n", req.containerID, err)
				} else {
					fmt.Printf("[BoxSvc] Container %s stopped successfully\n", req.containerID)
				}

				_, err := s.repo.UpdateStatus(context.Background(), req.fingerprint, string(domain.StatusPaused))
				if err != nil {
					fmt.Printf("[BoxSvc] Error updating box status to paused: %v\n", err)
				} else {
					fmt.Printf("[BoxSvc] Box status updated to paused\n")
				}

				delete(shutdownTimers, req.fingerprint)
			} else {
				fmt.Printf("[BoxSvc] Shutdown cancelled for %s (new connections arrived)\n", req.fingerprint)
				delete(shutdownTimers, req.fingerprint)
			}
		}
	}
}

func (s *Svc) incrementConnection(fingerprint string) {
	responseCh := make(chan struct{})
	s.connEventCh <- connEvent{
		fingerprint: fingerprint,
		increment:   true,
		responseCh:  responseCh,
	}
	<-responseCh
}

func (s *Svc) decrementConnection(fingerprint, containerID string) {
	responseCh := make(chan struct{})
	s.connEventCh <- connEvent{
		fingerprint: fingerprint,
		increment:   false,
		responseCh:  responseCh,
	}
	<-responseCh

	time.AfterFunc(5*time.Second, func() {
		s.shutdownCh <- shutdownRequest{
			fingerprint: fingerprint,
			containerID: containerID,
		}
	})
}

func (s *Svc) Connect(ctx context.Context, conn *websocket.Conn, fingerprint string) error {
	// Validate fingerprint
	if strings.TrimSpace(fingerprint) == "" {
		return domain.NewValidationError("fingerprint cannot be empty")
	}

	s.incrementConnection(fingerprint)

	box, err := s.repo.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		// Check if it's a NotFound error - this is expected for new boxes
		if appErr, ok := domain.IsAppError(err); !ok || !appErr.IsType(domain.ErrorTypeNotFound) {
			s.decrementConnection(fingerprint, "")
			return err // Return AppError as-is
		}
		// box not found is expected, box will be nil and we'll create a new one
		box = nil
	}

	if box == nil {
		containerID, err := s.dockerSvc.CreateContainer(ctx)
		if err != nil {
			s.decrementConnection(fingerprint, "")
			return err // Already an AppError from dockerSvc
		}
		if err := s.dockerSvc.StartIfNotRunning(ctx, containerID); err != nil {
			s.decrementConnection(fingerprint, containerID)
			return err // Already an AppError from dockerSvc
		}
		newBox := domain.Box{
			FingerprintID: fingerprint,
			ContainerID:   containerID,
			Status:        domain.StatusRunning,
			LastActive:    time.Now(),
		}
		box, err = s.repo.Create(ctx, newBox)
		if err != nil {
			s.decrementConnection(fingerprint, containerID)
			return err // Already an AppError from repo
		}
		fmt.Printf("[BoxSvc] Created new box with container: %s\n", containerID)
	} else {
		if err := s.dockerSvc.StartIfNotRunning(ctx, box.ContainerID); err != nil {
			s.decrementConnection(fingerprint, box.ContainerID)
			return err // Already an AppError from dockerSvc
		}

		_, err := s.repo.UpdateStatus(ctx, fingerprint, string(domain.StatusRunning))
		if err != nil {
			fmt.Printf("[BoxSvc] Warning: failed to update box status to running: %v\n", err)
		}
		fmt.Printf("[BoxSvc] Reconnected to existing box with container: %s\n", box.ContainerID)
	}

	attachResp, err := s.dockerSvc.AttachContainer(ctx, box.ContainerID)
	if err != nil {
		s.decrementConnection(fingerprint, box.ContainerID)
		return err // Already an AppError from dockerSvc
	}
	defer attachResp.Close()

	done := make(chan struct{})

	// container → websocket
	go func() {
		buf := make([]byte, 8192)
		for {
			select {
			case <-done:
				return
			default:
				n, err := attachResp.Reader.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Println("[BoxSvc] Container read error:", err)
					}
					return
				}
				if n > 0 {
					if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
						fmt.Println("[BoxSvc] Error writing to websocket:", err)
						return
					}
				}
			}
		}
	}()

	defer func() {
		close(done)
		s.decrementConnection(fingerprint, box.ContainerID)
	}()

	// websocket input → container stdin
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("[BoxSvc] WebSocket closed normally")
				return nil
			}
			fmt.Println("[BoxSvc] Error reading from websocket:", err)
			return domain.NewInternalError("websocket read error", err)
		}

		if len(msg) > 0 {
			_, err := attachResp.Conn.Write(msg)
			if err != nil {
				fmt.Println("[BoxSvc] Error writing to container stdin:", err)
				return domain.NewInternalError("container write error", err)
			}
		}
	}
}
