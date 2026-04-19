package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/digikeys/backend/internal/application"
)

type VerificationHandler struct {
	verifyService *application.VerificationService
}

func NewVerificationHandler(verifyService *application.VerificationService) *VerificationHandler {
	return &VerificationHandler{verifyService: verifyService}
}

func (h *VerificationHandler) VerifyCard(w http.ResponseWriter, r *http.Request) {
	cardNumber := chi.URLParam(r, "cardNumber")
	if cardNumber == "" {
		cardNumber = r.URL.Query().Get("cardNumber")
	}

	if cardNumber == "" {
		writeError(w, http.StatusBadRequest, "card number is required")
		return
	}

	result, err := h.verifyService.VerifyCard(r.Context(), cardNumber)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
