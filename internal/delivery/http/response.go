package http

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, status, message string, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    code,
		"status":  status,
		"message": message,
		"data":    data,
	})
}

func WriteError(w http.ResponseWriter, code int, status, message string) {
	WriteJSON(w, code, status, message, nil)
}
