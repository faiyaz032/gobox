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
	return nil
}
