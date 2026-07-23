package health

import (
	"context"
	"log"
	"net/http"
	"taskflow/internal/response"
	"time"
)

type DatabasePinger interface {
	PingContext(ctx context.Context) error
}

type Handler struct {
	database DatabasePinger
}

func NewHandler(database DatabasePinger) *Handler {
	return &Handler{
		database: database,
	}
}

func (h *Handler) HealthCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := response.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	); err != nil {
		log.Printf("failed to write health response: %v", err)
	}
}

func (h *Handler) ReadinessCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		2*time.Second,
	)
	defer cancel()

	if err := h.database.PingContext(ctx); err != nil {
		log.Printf("database readiness check failed: %v", err)

		if writeErr := response.WriteJSON(
			w,
			http.StatusServiceUnavailable,
			map[string]string{
				"status": "not ready",
			},
		); writeErr != nil {
			log.Printf(
				"failed to write readiness response: %v",
				writeErr,
			)
		}
		return
	}

	if err := response.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ready",
		},
	); err != nil {
		log.Printf("failed to write readiness response: %v", err)
	}
}
