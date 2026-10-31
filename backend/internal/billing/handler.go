package billing

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

type PreviewRequest struct {
	BasePriceCents    int64 `json:"basePriceCents"`
	VariantAdjustment int64 `json:"variantAdjustmentCents"`
	Quantity          int   `json:"quantity"`
	DiscountCents     int64 `json:"discountCents"`
	TaxRateBps        int   `json:"taxRateBps"`
}

func (h *Handler) PreviewCalculation(w http.ResponseWriter, r *http.Request) {
	var req PreviewRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse preview payload", nil)
		return
	}

	breakdown := h.service.PreviewCalculation(r.Context(), req.BasePriceCents, req.VariantAdjustment, req.Quantity, req.DiscountCents, req.TaxRateBps)
	platform.WriteJSON(w, http.StatusOK, map[string]any{"breakdown": breakdown})
}

func (h *Handler) GetReliabilityCenter(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.service.GetReliabilityMetrics(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to load reliability metrics", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"metrics": metrics})
}

func (h *Handler) TriggerBillingRun(w http.ResponseWriter, r *http.Request) {
	run, err := h.service.RunAutomatedBilling(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "RUN_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Automated billing run completed",
		"run":     run,
	})
}

func (h *Handler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	user := platform.GetUserFromContext(r.Context())
	var customerID *string

	if user != nil && user.Role == users.RoleCustomer {
		customerID = &user.ID
	} else if qCust := r.URL.Query().Get("customerId"); qCust != "" {
		customerID = &qCust
	}

	var status *string
	if qStatus := r.URL.Query().Get("status"); qStatus != "" {
		status = &qStatus
	}

	invoices, err := h.service.ListInvoices(r.Context(), customerID, status)
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve invoices", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"invoices": invoices})
}

func (h *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	inv, lines, err := h.service.GetInvoice(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvoiceNotFound) {
			platform.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Invoice not found", nil)
			return
		}
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed retrieving invoice", nil)
		return
	}

	user := platform.GetUserFromContext(r.Context())
	if user != nil && user.Role == users.RoleCustomer && inv.CustomerID != user.ID {
		platform.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		return
	}

	inv.Lines = lines
	platform.WriteJSON(w, http.StatusOK, map[string]any{"invoice": inv})
}

func (h *Handler) ConfirmInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.ConfirmInvoice(r.Context(), id); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "CONFIRM_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]string{"message": "Invoice confirmed successfully"})
}
