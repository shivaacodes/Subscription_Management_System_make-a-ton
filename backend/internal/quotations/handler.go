package quotations

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/subflow/backend/internal/platform"
	"github.com/subflow/backend/internal/users"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateQuotationRequest struct {
	CustomerID   string          `json:"customerId"`
	PlanID       string          `json:"planId"`
	ValidityDays int             `json:"validityDays"`
	Lines        []QuotationLine `json:"lines"`
}

func (h *Handler) ListQuotations(w http.ResponseWriter, r *http.Request) {
	user := platform.GetUserFromContext(r.Context())
	var customerID *string

	if user != nil && user.Role == users.RoleCustomer {
		customerID = &user.ID
	} else if qCust := r.URL.Query().Get("customerId"); qCust != "" {
		customerID = &qCust
	}

	quotations, err := h.service.ListQuotations(r.Context(), customerID)
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve quotations", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"quotations": quotations})
}

func (h *Handler) CreateQuotation(w http.ResponseWriter, r *http.Request) {
	var req CreateQuotationRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse quotation payload", nil)
		return
	}

	user := platform.GetUserFromContext(r.Context())
	if user != nil && user.Role == users.RoleCustomer {
		req.CustomerID = user.ID
	}

	q, err := h.service.CreateQuotation(r.Context(), req.CustomerID, req.PlanID, req.ValidityDays, req.Lines)
	if err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{"quotation": q})
}

func (h *Handler) GetQuotation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	q, lines, err := h.service.GetQuotation(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrQuotationNotFound) {
			platform.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Quotation not found", nil)
			return
		}
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed retrieving quotation", nil)
		return
	}

	user := platform.GetUserFromContext(r.Context())
	if user != nil && user.Role == users.RoleCustomer && q.CustomerID != user.ID {
		platform.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		return
	}

	q.Lines = lines
	platform.WriteJSON(w, http.StatusOK, map[string]any{"quotation": q})
}

func (h *Handler) ConvertQuotation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subID, err := h.service.ConvertQuotation(r.Context(), id)
	if err != nil {
		platform.WriteError(w, http.StatusBadRequest, "CONVERSION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusOK, map[string]any{
		"message":        "Quotation converted to active subscription successfully",
		"subscriptionId": subID,
	})
}

func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.service.ListTemplates(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve templates", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"templates": templates})
}
