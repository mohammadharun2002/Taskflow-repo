package project

import (
	"encoding/json"
	"errors"
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
	var request CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		if writeErr := response.WriteError(
			w,
			http.StatusBadRequest,
			"Invalid Request Body",
		); writeErr != nil {
			log.Printf("failed to write error response : %v", writeErr)
		}
		return
	}

	createdProject, err := h.service.Create(request)
	if errors.Is(err, ErrNameRequired) {
		if writeErr := response.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}
	if err != nil {
		log.Printf("failed to create project: %v", err)
		if writeErr := response.WriteError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	if err := response.WriteJSON(
		w,
		http.StatusCreated,
		createdProject,
	); err != nil {
		log.Printf("failed to write project response: %v", err)
	}
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	projects := h.service.FindAll()
	if err := response.WriteJSON(w, http.StatusOK, projects); err != nil {
		log.Printf("failed to write projects response: %v", err)
	}
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		if writeErr := response.WriteError(
			w,
			http.StatusBadRequest,
			"invalid project ID",
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	foundProject, err := h.service.FindByID(id)
	if errors.Is(err, ErrNotFound) {
		if writeErr := response.WriteError(
			w,
			http.StatusNotFound,
			err.Error(),
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	if err != nil {
		log.Printf("failed to find project: %v", err)

		if writeErr := response.WriteError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	if err := response.WriteJSON(
		w,
		http.StatusOK,
		foundProject,
	); err != nil {
		log.Printf("failed to write project response: %v", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		if writeErr := response.WriteError(
			w,
			http.StatusBadRequest,
			"invalid project ID",
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	err = h.service.Delete(id)
	if errors.Is(err, ErrNotFound) {
		if writeErr := response.WriteError(
			w,
			http.StatusNotFound,
			err.Error(),
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	if err != nil {
		log.Printf("failed to delete project: %v", err)

		if writeErr := response.WriteError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
