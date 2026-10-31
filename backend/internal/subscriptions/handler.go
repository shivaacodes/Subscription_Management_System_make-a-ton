package subscriptions

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

type CreateSubscriptionRequest struct {
	CustomerID   string             `json:"customerId"`
	PlanID       string             `json:"planId"`
	PaymentTerms string             `json:"paymentTerms"`
	Lines        []SubscriptionLine `json:"lines"`
}

type ToggleAddonRequest struct {
	ProductID string `json:"productId"`
	Active    bool   `json:"active"`
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
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

	subs, err := h.service.ListSubscriptions(r.Context(), customerID, status)
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve subscriptions", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"subscriptions": subs})
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req CreateSubscriptionRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse subscription payload", nil)
		return
	}

	user := platform.GetUserFromContext(r.Context())
	var actorID *string
	if user != nil {
		actorID = &user.ID
		if user.Role == users.RoleCustomer {
			req.CustomerID = user.ID
		}
	}

	sub := &Subscription{
		CustomerID:   req.CustomerID,
		PlanID:       req.PlanID,
		PaymentTerms: req.PaymentTerms,
	}

	created, err := h.service.CreateSubscription(r.Context(), sub, req.Lines, actorID)
	if err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{"subscription": created})
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sub, lines, err := h.service.GetSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			platform.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Subscription not found", nil)
			return
		}
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed retrieving subscription", nil)
		return
	}

	user := platform.GetUserFromContext(r.Context())
	if user != nil && user.Role == users.RoleCustomer && sub.CustomerID != user.ID {
		platform.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Access denied", nil)
		return
	}

	sub.Lines = lines
	platform.WriteJSON(w, http.StatusOK, map[string]any{"subscription": sub})
}

func (h *Handler) ActivateSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := platform.GetUserFromContext(r.Context())
	var actorID *string
	if user != nil {
		actorID = &user.ID
	}

	if err := h.service.TransitionStatus(r.Context(), id, StatusActive, actorID); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "TRANSITION_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]string{"message": "Subscription activated successfully"})
}

func (h *Handler) PauseSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := platform.GetUserFromContext(r.Context())
	var actorID *string
	if user != nil {
		actorID = &user.ID
	}

	if err := h.service.PauseSubscription(r.Context(), id, actorID); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "PAUSE_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]string{"message": "Subscription paused successfully"})
}

func (h *Handler) ResumeSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := platform.GetUserFromContext(r.Context())
	var actorID *string
	if user != nil {
		actorID = &user.ID
	}

	if err := h.service.ResumeSubscription(r.Context(), id, actorID); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "RESUME_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]string{"message": "Subscription resumed successfully"})
}

func (h *Handler) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := platform.GetUserFromContext(r.Context())
	var actorID *string
	if user != nil {
		actorID = &user.ID
	}

	if err := h.service.CancelSubscription(r.Context(), id, actorID); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "CANCEL_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]string{"message": "Subscription cancelled successfully"})
}

func (h *Handler) ToggleAddon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ToggleAddonRequest
	if err := platform.ReadJSON(r, &req); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse addon payload", nil)
		return
	}

	user := platform.GetUserFromContext(r.Context())
	var actorID *string
	if user != nil {
		actorID = &user.ID
	}

	if err := h.service.ToggleAddon(r.Context(), id, req.ProductID, req.Active, actorID); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "ADDON_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]string{"message": "Addon preference updated for next billing cycle"})
}

func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	events, err := h.service.GetTimeline(r.Context(), id)
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve timeline", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"events": events})
}
