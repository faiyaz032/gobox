package box

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/faiyaz032/gobox/internal/domain"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *Svc) Connect(ctx context.Context, conn *websocket.Conn, fingerprint string) error {
	if strings.TrimSpace(fingerprint) == "" {
		return domain.NewValidationError("fingerprint cannot be empty")
	}

	s.incrementConnection(fingerprint)

	box, err := s.repo.GetByFingerprint(ctx, fingerprint)
	if err != nil {

		if appErr, ok := domain.IsAppError(err); !ok || !appErr.IsType(domain.ErrorTypeNotFound) {
			s.decrementConnection(fingerprint, "")
			return err
		}
		box = nil
	}

	if box == nil {
		containerID, err := s.dockerSvc.CreateContainer(ctx)
		if err != nil {
			s.decrementConnection(fingerprint, "")
			return err
		}
		if err := s.dockerSvc.StartIfNotRunning(ctx, containerID); err != nil {
			s.decrementConnection(fingerprint, containerID)
			return err
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
			return err
		}
		s.logger.Info("Created new box with container",
			zap.String("container_id", containerID),
			zap.String("fingerprint", fingerprint))
	} else {
		if err := s.dockerSvc.StartIfNotRunning(ctx, box.ContainerID); err != nil {
			s.decrementConnection(fingerprint, box.ContainerID)
			return err
		}

		_, err := s.repo.UpdateStatus(ctx, fingerprint, string(domain.StatusRunning))
		if err != nil {
			s.logger.Warn("Failed to update box status to running",
				zap.String("fingerprint", fingerprint),
				zap.Error(err))
		}
		s.logger.Info("Reconnected to existing box",
			zap.String("container_id", box.ContainerID),
			zap.String("fingerprint", fingerprint))
	}

	attachResp, err := s.dockerSvc.AttachContainer(ctx, box.ContainerID)
	if err != nil {
		s.decrementConnection(fingerprint, box.ContainerID)
		return err
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
						s.logger.Error("Container read error", zap.Error(err))
					}
					return
				}
				if n > 0 {
					if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
						s.logger.Error("Error writing to websocket", zap.Error(err))
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
				s.logger.Info("WebSocket closed normally")
				return nil
			}
			s.logger.Error("Error reading from websocket", zap.Error(err))
			return domain.NewInternalError("websocket read error", err)
		}

		if len(msg) > 0 {
			_, err := attachResp.Conn.Write(msg)
			if err != nil {
				s.logger.Error("Error writing to container stdin", zap.Error(err))
				return domain.NewInternalError("container write error", err)
			}
		}
	}
}
