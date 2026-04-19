package middleware

import (
	"net/http"

	"github.com/digikeys/backend/internal/domain"
)

// RequireRole ensures the authenticated user has one of the specified roles.
func RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())

			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, `{"error":"forbidden: insufficient permissions"}`, http.StatusForbidden)
		})
	}
}

// RequireSuperAdmin restricts access to super_admin only.
func RequireSuperAdmin(next http.Handler) http.Handler {
	return RequireRole(domain.UserRoleSuperAdmin)(next)
}

// RequireEmbassyAdmin restricts access to embassy_admin and super_admin.
func RequireEmbassyAdmin(next http.Handler) http.Handler {
	return RequireRole(domain.UserRoleSuperAdmin, domain.UserRoleEmbassyAdmin)(next)
}

// RequireEnrollmentAgent restricts access to enrollment agents.
func RequireEnrollmentAgent(next http.Handler) http.Handler {
	return RequireRole(domain.UserRoleSuperAdmin, domain.UserRoleEmbassyAdmin, domain.UserRoleEnrollmentAgent)(next)
}
