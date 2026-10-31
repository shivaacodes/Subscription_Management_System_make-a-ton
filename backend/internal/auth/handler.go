package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/subflow/backend/internal/config"
	"github.com/subflow/backend/internal/platform"
	"github.com/subflow/backend/internal/users"
)

type RegisterRequest struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	CompanyName *string `json:"companyName,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User    *users.User `json:"user"`
	Message string      `json:"message,omitempty"`
}

type Handler struct {
	service Service
	cfg     *config.Config
}

func NewHandler(service Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse JSON body", nil)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := extractIP(r)

	user, token, err := h.service.RegisterCustomer(r.Context(), req.Email, req.Password, req.CompanyName, &userAgent, &ipAddress)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrDuplicateEmail):
			platform.WriteError(w, http.StatusConflict, "DUPLICATE_EMAIL", "An account with this email address already exists", nil)
		case errors.Is(err, ErrInvalidEmail):
			platform.WriteError(w, http.StatusBadRequest, "INVALID_EMAIL", "The provided email address is invalid", nil)
		case errors.Is(err, ErrPasswordTooShort),
			errors.Is(err, ErrPasswordNoUpper),
			errors.Is(err, ErrPasswordNoLower),
			errors.Is(err, ErrPasswordNoSpecial):
			platform.WriteError(w, http.StatusBadRequest, "WEAK_PASSWORD", err.Error(), nil)
		default:
			platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Registration failed", nil)
		}
		return
	}

	h.setSessionCookie(w, token, 7*24*time.Hour)
	platform.WriteJSON(w, http.StatusCreated, AuthResponse{
		User:    user,
		Message: "Registration successful",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse JSON body", nil)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := extractIP(r)

	user, token, err := h.service.Login(r.Context(), req.Email, req.Password, &userAgent, &ipAddress)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			platform.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Login failed", nil)
		return
	}

	h.setSessionCookie(w, token, 7*24*time.Hour)
	platform.WriteJSON(w, http.StatusOK, AuthResponse{
		User:    user,
		Message: "Login successful",
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(platform.SessionCookieName)
	if err == nil && cookie.Value != "" {
		_ = h.service.Logout(r.Context(), cookie.Value)
	}

	h.clearSessionCookie(w)
	platform.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := platform.GetUserFromContext(r.Context())
	if user == nil {
		platform.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated", nil)
		return
	}

	platform.WriteJSON(w, http.StatusOK, map[string]any{
		"user": user,
	})
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, rawToken string, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     platform.SessionCookieName,
		Value:    rawToken,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	if h.cfg.CookieDomain != "" {
		cookie.Domain = h.cfg.CookieDomain
	}
	http.SetCookie(w, cookie)
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     platform.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	if h.cfg.CookieDomain != "" {
		cookie.Domain = h.cfg.CookieDomain
	}
	http.SetCookie(w, cookie)
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return r.RemoteAddr
}
