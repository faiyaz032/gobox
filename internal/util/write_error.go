package util

import (
	"encoding/json"
	"net/http"

	"github.com/faiyaz032/gobox/internal/errors"
)

func WriteError(w http.ResponseWriter, err *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
