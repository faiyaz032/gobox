package box

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/faiyaz032/gobox/internal/domain"
	"github.com/gorilla/websocket"
)

type Svc struct {
	repo      Repo
	dockerSvc DockerSvc
}

func NewSvc(repo Repo, dockerSvc DockerSvc) *Svc {
	return &Svc{
		repo:      repo,
		dockerSvc: dockerSvc,
	}
}

func (s *Svc) Connect(ctx context.Context, conn *websocket.Conn, fingerprint string) error {
	box, err := s.repo.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to get box by fingerprint: %w", err)
	}

	if box == nil {
		containerID, err := s.dockerSvc.CreateContainer(ctx)
		if err != nil {
			return fmt.Errorf("failed to create container: %w", err)
		}

		if err := s.dockerSvc.StartIfNotRunning(ctx, containerID); err != nil {
			return fmt.Errorf("failed to start new container: %w", err)
		}

		newBox := domain.Box{
			FingerprintID: fingerprint,
			ContainerID:   containerID,
			Status:        domain.StatusRunning,
			LastActive:    time.Now(),
		}

		box, err = s.repo.Create(ctx, newBox)
		if err != nil {
			return fmt.Errorf("failed to create box in repo: %w", err)
		}
	} else {
		if err := s.dockerSvc.StartIfNotRunning(ctx, box.ContainerID); err != nil {
			return fmt.Errorf("failed to start container: %w", err)
		}
	}

	attachResp, err := s.dockerSvc.AttachContainer(ctx, box.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to attach to container: %w", err)
	}
	defer attachResp.Close()

	// container → websocket
	go func() {
		buf := make([]byte, 8192)
		for {
			n, err := attachResp.Reader.Read(buf)
			if err != nil {
				if err != io.EOF {
				}
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
		}
	}()

	// websocket input → container stdin
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			return fmt.Errorf("read message error: %w", err)
		}

		if len(msg) > 0 {
			_, err := attachResp.Conn.Write(msg)
			if err != nil {
				return fmt.Errorf("write to container error: %w", err)
			}
		}
	}
}
