package response

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) error {
	return WriteJSON(w, statusCode, map[string]string{
		"error": message,
	})
}
