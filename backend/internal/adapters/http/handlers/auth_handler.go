package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/pkg/validator"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Phone     string `json:"phone"`
	Role      string `json:"role" validate:"required,oneof=super_admin embassy_admin enrollment_agent print_operator bank_agent verifier readonly"`
	EmbassyID string `json:"embassyId"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresIn    int          `json:"expiresIn"`
	User         *domain.User `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		writeDomainError(w, err)
		return
	}

	tokens, user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeDomainError(w, domain.ErrUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User:         user,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		writeDomainError(w, err)
		return
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      domain.UserRole(req.Role),
		EmbassyID: req.EmbassyID,
	}

	if err := h.authService.Register(r.Context(), user, req.Password); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Validate(body); err != nil {
		writeDomainError(w, err)
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		writeDomainError(w, domain.ErrUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}
