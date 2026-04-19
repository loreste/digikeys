package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/digikeys/backend/internal/adapters/http/middleware"
	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
	"github.com/digikeys/backend/pkg/validator"
)

type CitizenHandler struct {
	citizenService *application.CitizenService
}

func NewCitizenHandler(citizenService *application.CitizenService) *CitizenHandler {
	return &CitizenHandler{citizenService: citizenService}
}

func (h *CitizenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var citizen domain.Citizen
	if err := json.NewDecoder(r.Body).Decode(&citizen); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	citizen.EmbassyID = middleware.GetEmbassyID(r.Context())
	citizen.RegisteredBy = middleware.GetUserID(r.Context())

	if err := h.citizenService.Create(r.Context(), &citizen); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, citizen)
}

func (h *CitizenHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	citizen, err := h.citizenService.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, citizen)
}

func (h *CitizenHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var citizen domain.Citizen
	if err := json.NewDecoder(r.Body).Decode(&citizen); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	citizen.ID = id
	if err := h.citizenService.Update(r.Context(), &citizen); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, citizen)
}

type SearchRequest struct {
	Query string `json:"query" validate:"required,min=2"`
}

func (h *CitizenHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		var req SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			query = req.Query
		}
	}

	if len(query) < 2 {
		if err := validator.Validate(SearchRequest{Query: query}); err != nil {
			writeDomainError(w, err)
			return
		}
	}

	page, pageSize := getPageParams(r)
	embassyID := middleware.GetEmbassyID(r.Context())

	filter := domain.CitizenFilter{
		Query:     query,
		EmbassyID: embassyID,
		Page:      page,
		PageSize:  pageSize,
	}

	citizens, total, err := h.citizenService.Search(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":       citizens,
		"pagination": newPagination(page, pageSize, total),
	})
}

func (h *CitizenHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPageParams(r)
	embassyID := middleware.GetEmbassyID(r.Context())

	// Super admin can filter by embassy
	if r.URL.Query().Get("embassyId") != "" {
		role := middleware.GetUserRole(r.Context())
		if role == domain.UserRoleSuperAdmin {
			embassyID = r.URL.Query().Get("embassyId")
		}
	}

	filter := domain.CitizenFilter{
		EmbassyID:          embassyID,
		CountryOfResidence: r.URL.Query().Get("country"),
		Status:             r.URL.Query().Get("status"),
		Page:               page,
		PageSize:           pageSize,
	}

	citizens, total, err := h.citizenService.List(r.Context(), filter)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":       citizens,
		"pagination": newPagination(page, pageSize, total),
	})
}
