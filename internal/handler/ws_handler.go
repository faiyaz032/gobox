package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	appErr "github.com/faiyaz032/gobox/internal/errors"
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

type WSMessage struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeWSError(ws *websocket.Conn, err error, context string) {
	log.Printf("[WebSocket error] %s: %v", context, err)

	msg := WSMessage{Type: "error"}

	var appError *appErr.AppError
	if appErr.IsAppError(err) {
		appError = err.(*appErr.AppError)
		msg.Error = appError.Message
	} else {
		msg.Error = context
	}

	ws.WriteJSON(msg)
}

func (h *Handler) HandleWS(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeWSError(ws, appErr.NewAppError(400, "session_id is required"), "invalid request")
		return nil
	}

	container, err := h.svc.Start(ctx, sessionID)
	if err != nil {
		writeWSError(ws, err, "failed to start session")
		return nil
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := container.Reader.Read(buf)
			if err != nil {
				writeWSError(ws, err, "container read error")
				break
			}
			ws.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if pauseErr := h.svc.Pause(ctx, container.ContainerID); pauseErr != nil {
				writeWSError(ws, pauseErr, "failed to pause container")
			}
			ws.WriteJSON(WSMessage{Type: "status", Data: "container paused"})
			break
		}

		if _, err := container.Conn.Write(msg); err != nil {
			writeWSError(ws, err, "failed to write to container")
			break
		}
	}

	return nil
}
