package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

func (h *Handler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAuthError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	user, err := h.service.Register(r.Context(), request)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(
		w,
		http.StatusCreated,
		user,
	); err != nil {
		log.Printf("failed to write registration response: %v", err)
	}
}

func (h *Handler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeAuthError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	authResponse, err := h.service.Login(r.Context(), request)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	if err := response.WriteJSON(
		w,
		http.StatusOK,
		authResponse,
	); err != nil {
		log.Printf("failed to write login response: %v", err)
	}
}

func (h *Handler) writeServiceError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(err, ErrNameRequired),
		errors.Is(err, ErrEmailRequired),
		errors.Is(err, ErrInvalidEmail),
		errors.Is(err, ErrPasswordTooShort):
		writeAuthError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrEmailAlreadyExists):
		writeAuthError(w, http.StatusConflict, err.Error())

	case errors.Is(err, ErrInvalidCredentials):
		writeAuthError(w, http.StatusUnauthorized, err.Error())

	default:
		log.Printf("auth service error: %v", err)
		writeAuthError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}

func writeAuthError(
	w http.ResponseWriter,
	statusCode int,
	message string,
) {
	if err := response.WriteError(
		w,
		statusCode,
		message,
	); err != nil {
		log.Printf("failed to write auth error response: %v", err)
	}
}
