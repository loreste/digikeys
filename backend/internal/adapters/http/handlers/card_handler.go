package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/digikeys/backend/internal/adapters/http/middleware"
	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
)

type CardHandler struct {
	cardService *application.CardService
}

func NewCardHandler(cardService *application.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

type RequestCardRequest struct {
	EnrollmentID string `json:"enrollmentId" validate:"required"`
}

func (h *CardHandler) RequestCard(w http.ResponseWriter, r *http.Request) {
	var req RequestCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	issuedBy := middleware.GetUserID(r.Context())

	card, err := h.cardService.RequestCard(r.Context(), req.EnrollmentID, issuedBy)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, card)
}

func (h *CardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	card, err := h.cardService.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, card)
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPageParams(r)
	embassyID := middleware.GetEmbassyID(r.Context())

	role := middleware.GetUserRole(r.Context())
	if role == domain.UserRoleSuperAdmin && r.URL.Query().Get("embassyId") != "" {
		embassyID = r.URL.Query().Get("embassyId")
	}

	filter := domain.CardFilter{
		EmbassyID: embassyID,
		Status:    r.URL.Query().Get("status"),
		CitizenID: r.URL.Query().Get("citizenId"),
		Page:      page,
		PageSize:  pageSize,
	}

	cards, total, err := h.cardService.List(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":       cards,
		"pagination": newPagination(page, pageSize, total),
	})
}

func (h *CardHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	approvedBy := middleware.GetUserID(r.Context())

	if err := h.cardService.ApproveCard(r.Context(), id, approvedBy); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "approved"})
}

type QueuePrintRequest struct {
	BatchID string `json:"batchId"`
}

func (h *CardHandler) QueueForPrinting(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req QueuePrintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.cardService.QueueForPrinting(r.Context(), id, req.BatchID); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "printing"})
}

func (h *CardHandler) MarkPrinted(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.cardService.MarkPrinted(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "printed"})
}

func (h *CardHandler) MarkDelivered(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.cardService.MarkDelivered(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "delivered"})
}

func (h *CardHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.cardService.Activate(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "active"})
}

func (h *CardHandler) Suspend(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.cardService.Suspend(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "suspended"})
}

func (h *CardHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.cardService.Revoke(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

type RenewRequest struct {
	Reason string `json:"reason" validate:"required"`
}

func (h *CardHandler) Renew(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req RenewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	issuedBy := middleware.GetUserID(r.Context())

	newCard, err := h.cardService.Renew(r.Context(), id, req.Reason, issuedBy)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newCard)
}
