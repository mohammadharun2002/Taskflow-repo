package health

import (
	"log"
	"net/http"
	"taskflow/internal/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	}); err != nil {
		log.Printf("failed to write health response: %v", err)
	}
}

func ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	if err := response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	}); err != nil {
		log.Printf("failed to write health response: %v", err)
	}
}
