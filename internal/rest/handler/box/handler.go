package boxhandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/faiyaz032/gobox/internal/domain"
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
		h.writeError(w, domain.NewValidationError("fingerprint query parameter is required"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for fingerprint=%s: %v", fingerprint, err)
		h.writeError(w, domain.NewInternalError("failed to upgrade websocket connection", err))
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established for fingerprint=%s", fingerprint)

	if err := h.svc.Connect(r.Context(), conn, fingerprint); err != nil {
		log.Printf("Connection error for fingerprint=%s: %v", fingerprint, err)
		
		// Extract error message for websocket close
		errorMsg := domain.GetErrorMessage(err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, errorMsg))
		return
	}

	log.Printf("WebSocket connection closed for fingerprint=%s", fingerprint)
}

// writeError writes an error response using AppError
func (h *Handler) writeError(w http.ResponseWriter, err error) {
	statusCode := domain.GetStatusCode(err)
	errorType := domain.GetErrorType(err)
	message := domain.GetErrorMessage(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"type":    errorType,
			"message": message,
		},
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}
