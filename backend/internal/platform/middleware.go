package platform

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/subflow/backend/internal/users"
)

type contextKey string

const (
	UserContextKey    contextKey = "subflow_user"
	SessionCookieName            = "subflow_session"
)

type SessionValidator interface {
	ValidateSessionToken(ctx context.Context, rawToken string) (*users.User, error)
}

func GetUserFromContext(ctx context.Context) *users.User {
	val := ctx.Value(UserContextKey)
	if u, ok := val.(*users.User); ok {
		return u
	}
	return nil
}

func SessionAuthMiddleware(validator SessionValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			user, err := validator.ValidateSessionToken(r.Context(), cookie.Value)
			if err != nil {
				// Cookie invalid or expired
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRole(roles ...users.Role) func(http.Handler) http.Handler {
	allowed := make(map[users.Role]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
				return
			}

			if !allowed[user.Role] {
				WriteError(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource", map[string]string{
					"role": string(user.Role),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("recovered from panic in http handler", "error", rec)
				WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error occurred", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
