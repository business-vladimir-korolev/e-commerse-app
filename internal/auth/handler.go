package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	user, err := h.service.Register(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidEmail):
			writeError(w, http.StatusBadRequest, "invalid email")

		case errors.Is(err, ErrInvalidPassword):
			writeError(w, http.StatusBadRequest, "invalid password")

		case errors.Is(err, ErrEmailAlreadyExists):
			writeError(w, http.StatusConflict, "email already exists")

		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	writeJSON(w, http.StatusCreated, RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	token, err := h.service.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
