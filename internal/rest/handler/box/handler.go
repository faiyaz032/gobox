package boxhandler

import (
	"encoding/json"
	"net/http"

	"github.com/faiyaz032/gobox/internal/domain"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	svc    Svc
	logger *zap.Logger
}

func NewHandler(svc Svc, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
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
		h.logger.Error("Failed to upgrade connection",
			zap.String("fingerprint", fingerprint),
			zap.Error(err))
		h.writeError(w, domain.NewInternalError("failed to upgrade websocket connection", err))
		return
	}
	defer conn.Close()

	h.logger.Info("WebSocket connection established",
		zap.String("fingerprint", fingerprint))

	if err := h.svc.Connect(r.Context(), conn, fingerprint); err != nil {
		h.logger.Error("Connection error",
			zap.String("fingerprint", fingerprint),
			zap.Error(err))
		
		// Extract error message for websocket close
		errorMsg := domain.GetErrorMessage(err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, errorMsg))
		return
	}

	h.logger.Info("WebSocket connection closed",
		zap.String("fingerprint", fingerprint))
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
		h.logger.Error("Failed to encode error response", zap.Error(err))
	}
}
