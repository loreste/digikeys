package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/digikeys/backend/internal/adapters/http/middleware"
	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
)

type EnrollmentHandler struct {
	enrollmentService *application.EnrollmentService
}

func NewEnrollmentHandler(enrollmentService *application.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{enrollmentService: enrollmentService}
}

func (h *EnrollmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var enrollment domain.Enrollment
	if err := json.NewDecoder(r.Body).Decode(&enrollment); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	enrollment.AgentID = middleware.GetUserID(r.Context())
	enrollment.EmbassyID = middleware.GetEmbassyID(r.Context())

	if err := h.enrollmentService.CreateEnrollment(r.Context(), &enrollment); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, enrollment)
}

func (h *EnrollmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	enrollment, err := h.enrollmentService.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, enrollment)
}

func (h *EnrollmentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPageParams(r)
	embassyID := middleware.GetEmbassyID(r.Context())

	role := middleware.GetUserRole(r.Context())
	if role == domain.UserRoleSuperAdmin && r.URL.Query().Get("embassyId") != "" {
		embassyID = r.URL.Query().Get("embassyId")
	}

	filter := domain.EnrollmentFilter{
		EmbassyID:    embassyID,
		AgentID:      r.URL.Query().Get("agentId"),
		SyncStatus:   r.URL.Query().Get("syncStatus"),
		ReviewStatus: r.URL.Query().Get("reviewStatus"),
		Page:         page,
		PageSize:     pageSize,
	}

	enrollments, total, err := h.enrollmentService.List(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":       enrollments,
		"pagination": newPagination(page, pageSize, total),
	})
}

func (h *EnrollmentHandler) Sync(w http.ResponseWriter, r *http.Request) {
	var req application.MobileSyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	agentID := middleware.GetUserID(r.Context())
	embassyID := middleware.GetEmbassyID(r.Context())

	for i := range req.Enrollments {
		req.Enrollments[i].AgentID = agentID
		req.Enrollments[i].EmbassyID = embassyID
	}

	result, err := h.enrollmentService.SyncFromMobile(r.Context(), req.Enrollments)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type ReviewRequest struct {
	Status string `json:"status" validate:"required,oneof=approved rejected needs_correction"`
	Notes  string `json:"notes"`
}

func (h *EnrollmentHandler) Review(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reviewerID := middleware.GetUserID(r.Context())

	if err := h.enrollmentService.ReviewEnrollment(r.Context(), id, req.Status, reviewerID, req.Notes); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "reviewed"})
}
