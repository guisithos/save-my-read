package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/guisithos/save-my-read/internal/application"
	"github.com/guisithos/save-my-read/internal/domain/auth"
)

type AuthHandler struct {
	authService *application.AuthService
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email    string   `json:"email"`
		Password string   `json:"password"`
		Name     string   `json:"name"`
		Genres   []string `json:"genres"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if err := validateRegistration(req.Email, req.Password, req.Name); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.Name, req.Genres)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailAlreadyExists):
			respondError(w, http.StatusConflict, "Email already registered")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	// Generate token for auto-login
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Registration successful but failed to login")
		return
	}

	respondJSON(w, http.StatusCreated, Response{
		Success: true,
		Data: map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    token,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}

func validateRegistration(email, password, name string) error {
	if len(email) < 5 || !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	return nil
}
