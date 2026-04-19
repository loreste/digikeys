package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/digikeys/backend/internal/adapters/http/handlers"
	"github.com/digikeys/backend/internal/adapters/http/middleware"
	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
)

type RouterDeps struct {
	AuthService       *application.AuthService
	CitizenService    *application.CitizenService
	EnrollmentService *application.EnrollmentService
	CardService       *application.CardService
	VerifyService     *application.VerificationService
	TransferService   *application.TransferService
	FSBService        *application.FSBService
	StatsService      *application.StatisticsService
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"carteconsulaire"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// ── Public ───────────────────────────────────────
		authHandler := handlers.NewAuthHandler(deps.AuthService)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)
			r.Post("/refresh", authHandler.RefreshToken)
		})

		// Public card verification
		if deps.VerifyService != nil {
			vh := handlers.NewVerificationHandler(deps.VerifyService)
			r.Route("/verify", func(r chi.Router) {
				r.Get("/{cardNumber}", vh.VerifyCard)
			})
		}

		// ── Authenticated ────────────────────────────────
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.AuthService))

			// ── Citizens (embassy_admin, enrollment_agent, super_admin) ──
			if deps.CitizenService != nil {
				ch := handlers.NewCitizenHandler(deps.CitizenService)
				r.Route("/citizens", func(r chi.Router) {
					r.Use(middleware.RequireRole(
						domain.UserRoleSuperAdmin,
						domain.UserRoleEmbassyAdmin,
						domain.UserRoleEnrollmentAgent,
					))
					r.Get("/", ch.List)
					r.Post("/", ch.Create)
					r.Get("/search", ch.Search)
					r.Get("/{id}", ch.Get)
					r.Put("/{id}", ch.Update)
				})
			}

			// ── Enrollments ──────────────────────────────
			if deps.EnrollmentService != nil {
				eh := handlers.NewEnrollmentHandler(deps.EnrollmentService)
				r.Route("/enrollments", func(r chi.Router) {
					r.Use(middleware.RequireRole(
						domain.UserRoleSuperAdmin,
						domain.UserRoleEmbassyAdmin,
						domain.UserRoleEnrollmentAgent,
					))
					r.Get("/", eh.List)
					r.Post("/", eh.Create)
					r.Post("/sync", eh.Sync)
					r.Get("/{id}", eh.Get)
					r.Post("/{id}/review", eh.Review)
				})
			}

			// ── Cards ────────────────────────────────────
			if deps.CardService != nil {
				cardH := handlers.NewCardHandler(deps.CardService)
				r.Route("/cards", func(r chi.Router) {
					r.Use(middleware.RequireRole(
						domain.UserRoleSuperAdmin,
						domain.UserRoleEmbassyAdmin,
						domain.UserRolePrintOperator,
					))
					r.Get("/", cardH.List)
					r.Post("/", cardH.RequestCard)
					r.Get("/{id}", cardH.Get)
					r.Post("/{id}/approve", cardH.Approve)
					r.Post("/{id}/print", cardH.QueueForPrinting)
					r.Post("/{id}/printed", cardH.MarkPrinted)
					r.Post("/{id}/delivered", cardH.MarkDelivered)
					r.Post("/{id}/activate", cardH.Activate)
					r.Post("/{id}/suspend", cardH.Suspend)
					r.Post("/{id}/revoke", cardH.Revoke)
					r.Post("/{id}/renew", cardH.Renew)
				})
			}

			// ── Admin ────────────────────────────────────
			r.Route("/admin", func(r chi.Router) {
				r.Use(middleware.RequireRole(domain.UserRoleSuperAdmin))

				if deps.StatsService != nil {
					r.Get("/statistics", func(w http.ResponseWriter, r *http.Request) {
						embassyID := r.URL.Query().Get("embassyId")
						stats, err := deps.StatsService.GetDashboard(r.Context(), embassyID)
						if err != nil {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(`{"error":"failed to load statistics"}`))
							return
						}
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						data, _ := json.Marshal(stats)
						w.Write(data)
					})
				}
			})
		})
	})

	return r
}
