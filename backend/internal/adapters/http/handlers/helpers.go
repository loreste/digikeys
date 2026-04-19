package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/pkg/pagination"
	"github.com/digikeys/backend/pkg/validator"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIError struct {
	Code    string   `json:"code"`
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func newPagination(page, pageSize, total int) pagination.Result {
	return pagination.NewResult(page, pageSize, total)
}

func getPageParams(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	p := pagination.NewParams(page, pageSize)
	return p.Page, p.PageSize
}

func writeDomainError(w http.ResponseWriter, err error) {
	var validationErr *validator.ValidationError
	if errors.As(err, &validationErr) {
		writeJSON(w, http.StatusBadRequest, APIError{
			Code:    "VALIDATION_ERROR",
			Error:   "Données invalides",
			Details: validationErr.Messages,
		})
		return
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, APIError{
			Code:  "NOT_FOUND",
			Error: "Ressource non trouvée",
		})

	case errors.Is(err, domain.ErrAlreadyExists):
		writeJSON(w, http.StatusConflict, APIError{
			Code:  "ALREADY_EXISTS",
			Error: "Cette ressource existe déjà",
		})

	case errors.Is(err, domain.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, APIError{
			Code:  "UNAUTHORIZED",
			Error: "Non autorisé",
		})

	case errors.Is(err, domain.ErrForbidden):
		writeJSON(w, http.StatusForbidden, APIError{
			Code:  "FORBIDDEN",
			Error: "Accès refusé",
		})

	case errors.Is(err, domain.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, APIError{
			Code:  "INVALID_INPUT",
			Error: err.Error(),
		})

	case errors.Is(err, domain.ErrBiometricFailed):
		writeJSON(w, http.StatusServiceUnavailable, APIError{
			Code:    "BIOMETRIC_ERROR",
			Error:   "Erreur biométrique",
			Details: []string{err.Error()},
		})

	case errors.Is(err, domain.ErrCardExpired):
		writeJSON(w, http.StatusConflict, APIError{
			Code:  "CARD_EXPIRED",
			Error: "La carte consulaire a expiré",
		})

	case errors.Is(err, domain.ErrInvalidMRZ):
		writeJSON(w, http.StatusBadRequest, APIError{
			Code:    "INVALID_MRZ",
			Error:   "Données MRZ invalides",
			Details: []string{err.Error()},
		})

	case errors.Is(err, domain.ErrSyncFailed):
		writeJSON(w, http.StatusServiceUnavailable, APIError{
			Code:  "SYNC_FAILED",
			Error: "Échec de la synchronisation",
		})

	default:
		writeJSON(w, http.StatusInternalServerError, APIError{
			Code:  "INTERNAL_ERROR",
			Error: "Erreur interne du serveur",
		})
	}
}
