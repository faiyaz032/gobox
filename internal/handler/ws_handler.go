package handler

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/faiyaz032/gobox/internal/service"
	"github.com/gorilla/websocket"
)

type Handler struct {
	apiClient *client.Client
	svc       *service.Service
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewHandler(svc *service.Service, apiClient *client.Client) *Handler {
	return &Handler{
		svc:       svc,
		apiClient: apiClient,
	}
}

func (h *Handler) HandleWS(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id missing", http.StatusBadRequest)
		return nil
	}

	resp, err := h.svc.Start(ctx, sessionID)
	if err != nil {
		return err
	}

	// read from container and write to WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := resp.Reader.Read(buf)
			if err != nil {
				break
			}
			conn.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	// read from WebSocket and write to container
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.svc.Pause(ctx, resp.ContainerID)
			break
		}
		if _, err := resp.Conn.Write(msg); err != nil {
			return err
		}
	}

	return nil
}
