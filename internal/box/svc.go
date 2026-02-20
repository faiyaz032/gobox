package box

import (
	"context"
	"time"

	"github.com/faiyaz032/gobox/internal/domain"
	"go.uber.org/zap"
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
	logger      *zap.Logger
	connEventCh chan connEvent
	shutdownCh  chan shutdownRequest
}

func NewSvc(repo Repo, dockerSvc DockerSvc, logger *zap.Logger) *Svc {
	svc := &Svc{
		repo:        repo,
		dockerSvc:   dockerSvc,
		logger:      logger,
		connEventCh: make(chan connEvent),
		shutdownCh:  make(chan shutdownRequest),
	}

	go svc.manageConnections()
	go svc.cleanupExpiredBoxes()

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
					s.logger.Info("Cancelled shutdown timer for fingerprint (new connection)",
						zap.String("fingerprint", event.fingerprint))
				}

				s.logger.Info("Connection established",
					zap.String("fingerprint", event.fingerprint),
					zap.Int("active_connections", activeConns[event.fingerprint]))
			} else {

				activeConns[event.fingerprint]--
				remaining := activeConns[event.fingerprint]

				if remaining <= 0 {
					delete(activeConns, event.fingerprint)
					s.logger.Info("Last connection closed, will shutdown in 5 seconds",
						zap.String("fingerprint", event.fingerprint))
				} else {
					s.logger.Info("Connection closed",
						zap.String("fingerprint", event.fingerprint),
						zap.Int("remaining_connections", remaining))
				}
			}

			close(event.responseCh)

		case req := <-s.shutdownCh:

			if activeConns[req.fingerprint] <= 0 {
				s.logger.Info("Executing shutdown for container",
					zap.String("container_id", req.containerID))

				if err := s.dockerSvc.StopContainer(context.Background(), req.containerID); err != nil {
					s.logger.Error("Error stopping container",
						zap.String("container_id", req.containerID),
						zap.Error(err))
				} else {
					s.logger.Info("Container stopped successfully",
						zap.String("container_id", req.containerID))
				}

				_, err := s.repo.UpdateStatus(context.Background(), req.fingerprint, string(domain.StatusPaused))
				if err != nil {
					s.logger.Error("Error updating box status to paused",
						zap.String("fingerprint", req.fingerprint),
						zap.Error(err))
				} else {
					s.logger.Info("Box status updated to paused",
						zap.String("fingerprint", req.fingerprint))
				}

				delete(shutdownTimers, req.fingerprint)
			} else {
				s.logger.Info("Shutdown cancelled (new connections arrived)",
					zap.String("fingerprint", req.fingerprint))
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

func (s *Svc) cleanupExpiredBoxes() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		expiredAfter := time.Now().Add(-24 * time.Hour)
		boxes, err := s.repo.GetExpiredBoxes(context.Background(), expiredAfter)
		if err != nil {
			s.logger.Error("failed to get expired boxes", zap.Error(err))
			continue
		}

		for _, b := range boxes {
			err := s.dockerSvc.RemoveContainer(context.Background(), b.ContainerID)
			if err != nil {
				s.logger.Error("failed to remove container", zap.String("container_id", b.ContainerID), zap.Error(err))
			}
			
			err = s.repo.Delete(context.Background(), b.FingerprintID)
			if err != nil {
				s.logger.Error("failed to delete box from db", zap.String("fingerprint", b.FingerprintID), zap.Error(err))
			}
			
			s.logger.Info("removed expired container", zap.String("container_id", b.ContainerID), zap.String("fingerprint", b.FingerprintID))
		}
	}
}
