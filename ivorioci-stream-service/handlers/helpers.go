package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ivorioci-stream-service/models"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, models.NewSuccess(data))
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, models.NewError(code, message))
}

// writeAppError maps an AppError (or any sentinel) to the correct HTTP status.
func writeAppError(w http.ResponseWriter, err error) {
	var appErr *models.AppError
	if errors.As(err, &appErr) {
		writeError(w, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
}

// queryInt reads an integer query param with a fallback default.
func queryInt(r *http.Request, key string, def int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return def
	}
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return def
	}
	return n
}
