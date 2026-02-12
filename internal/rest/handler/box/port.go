package boxhandler

import (
	"context"

	"github.com/gorilla/websocket"
)

type Svc interface {
	Connect(ctx context.Context, conn *websocket.Conn, fingerprint string) error
}
