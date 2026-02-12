package boxhandler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	svc Svc
}

func NewHandler(svc Svc) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	fingerprint := r.URL.Query().Get("fingerprint")
	if fingerprint == "" {
		http.Error(w, "fingerprint is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for fingerprint=%s: %v", fingerprint, err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established for fingerprint=%s", fingerprint)

	if err := h.svc.Connect(r.Context(), conn, fingerprint); err != nil {
		log.Printf("Connection error for fingerprint=%s: %v", fingerprint, err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		return
	}

	log.Printf("WebSocket connection closed for fingerprint=%s", fingerprint)
}
