package box

import (
	"context"

	"github.com/gorilla/websocket"
)

type Svc struct {
	repo Repo
}

func NewSvc(repo Repo) *Svc {
	return &Svc{
		repo: repo,
	}
}

func (s *Svc) Connect(ctx context.Context, conn *websocket.Conn, fingerprint string) error {
	box, err := s.repo.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		return err
	}

	if box == nil {
		/*
		* spin a new container
		* attach the container with conn for io forwarding
		* save the fingerprint - container mapping in the database
		*
		 */
	} else {
		/*
		* resume the container
		* attach the container with conn for io forwarding
		 */
	}

	return nil
}
