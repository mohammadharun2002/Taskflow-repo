package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"taskflow/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := parsePositiveID(r, "projectID")
	if err != nil {
		writeBadRequest(w, "invalid project ID")
		return
	}

	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	createdTask, err := h.service.Create(r.Context(), projectID, request)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusCreated, createdTask); err != nil {
		log.Printf("failed to write task response: %v", err)
	}
}

func (h *Handler) FindByProjectID(w http.ResponseWriter, r *http.Request) {
	projectID, err := parsePositiveID(r, "projectID")
	if err != nil {
		writeBadRequest(w, "invalid project ID")
		return
	}

	projectTasks, err := h.service.FindByProjectID(r.Context(), projectID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, projectTasks); err != nil {
		log.Printf("failed to write tasks response: %v", err)
	}
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := parsePositiveID(r, "id")
	if err != nil {
		writeBadRequest(w, "invalid task ID")
		return
	}

	foundTask, err := h.service.FindByID(r.Context(), taskID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, foundTask); err != nil {
		log.Printf("failed to write task response: %v", err)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := parsePositiveID(r, "id")
	if err != nil {
		writeBadRequest(w, "invalid task ID")
		return
	}

	var request CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeBadRequest(w, "invalid request body")
		return
	}

	updatedTask, err := h.service.Update(r.Context(), taskID, request)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, updatedTask); err != nil {
		log.Printf("failed to write task response: %v", err)
	}
}

func (h *Handler) AdvanceStatus(w http.ResponseWriter, r *http.Request) {
	taskID, err := parsePositiveID(r, "id")
	if err != nil {
		writeBadRequest(w, "invalid task ID")
		return
	}

	updatedTask, err := h.service.AdvanceStatus(r.Context(), taskID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(w, http.StatusOK, updatedTask); err != nil {
		log.Printf("failed to write task response: %v", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := parsePositiveID(r, "id")
	if err != nil {
		writeBadRequest(w, "invalid task ID")
		return
	}

	if err := h.service.Delete(r.Context(), taskID); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parsePositiveID(r *http.Request, parameter string) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue(parameter), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s", parameter)
	}
	return id, nil
}

func writeBadRequest(w http.ResponseWriter, message string) {
	if err := response.WriteError(w, http.StatusBadRequest, message); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}

func writeServiceError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, ErrNameRequired), errors.Is(err, ErrInvalidStatus):
		statusCode = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrProjectNotFound):
		statusCode = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, ErrAlreadyCompleted):
		statusCode = http.StatusConflict
		message = err.Error()
	default:
		log.Printf("task service error: %v", err)
	}

	if writeErr := response.WriteError(w, statusCode, message); writeErr != nil {
		log.Printf("failed to write error response: %v", writeErr)
	}
}

func (h *Handler) GetProjectSummary(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID, err := parsePositiveID(r, "id")
	if err != nil {
		writeBadRequest(w, "invalid project ID")
		return
	}

	summary, err := h.service.GetProjectSummary(r.Context(), projectID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(
		w,
		http.StatusOK,
		summary,
	); err != nil {
		log.Printf("failed to write project summary response: %v", err)
	}
}
