package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/digikeys/backend/internal/application"
	"github.com/digikeys/backend/internal/domain"
)

type contextKey string

const (
	UserIDKey    contextKey = "userID"
	UserRoleKey  contextKey = "userRole"
	EmbassyIDKey contextKey = "embassyID"
)

func Auth(authService *application.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := authService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			ctx = context.WithValue(ctx, EmbassyIDKey, claims.EmbassyID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

func GetUserRole(ctx context.Context) domain.UserRole {
	if v, ok := ctx.Value(UserRoleKey).(domain.UserRole); ok {
		return v
	}
	return ""
}

func GetEmbassyID(ctx context.Context) string {
	if v, ok := ctx.Value(EmbassyIDKey).(string); ok {
		return v
	}
	return ""
}
