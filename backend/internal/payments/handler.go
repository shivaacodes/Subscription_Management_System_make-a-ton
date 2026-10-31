package payments

import (
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

type RecordPaymentRequest struct {
	AmountCents int64   `json:"amountCents"`
	Method      string  `json:"method"`
	Notes       *string `json:"notes,omitempty"`
}

func (h *Handler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoiceId")
	var req RecordPaymentRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse payment payload", nil)
		return
	}

	p, err := h.service.RecordPayment(r.Context(), invoiceID, req.AmountCents, req.Method, req.Notes)
	if err != nil {
		platform.WriteError(w, http.StatusBadRequest, "PAYMENT_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{"payment": p})
}

func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	user := platform.GetUserFromContext(r.Context())
	var customerID *string

	if user != nil && user.Role == users.RoleCustomer {
		customerID = &user.ID
	} else if qCust := r.URL.Query().Get("customerId"); qCust != "" {
		customerID = &qCust
	}

	var invoiceID *string
	if qInv := r.URL.Query().Get("invoiceId"); qInv != "" {
		invoiceID = &qInv
	}

	payments, err := h.service.ListPayments(r.Context(), invoiceID, customerID)
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve payments", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"payments": payments})
}
