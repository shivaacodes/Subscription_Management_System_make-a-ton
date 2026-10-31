package plans

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/subflow/backend/internal/platform"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.service.ListPlans(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve plans", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"plans": plans})
}

func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var p Plan
	if err := platform.ReadJSON(r, &p); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse plan payload", nil)
		return
	}

	if err := h.service.CreatePlan(r.Context(), &p); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{"plan": p})
}

func (h *Handler) GetPlan(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.service.GetPlan(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPlanNotFound) {
			platform.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Plan not found", nil)
			return
		}
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve plan", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"plan": p})
}

func (h *Handler) LinkProduct(w http.ResponseWriter, r *http.Request) {
	planID := chi.URLParam(r, "planId")
	var pp PlanProduct
	if err := platform.ReadJSON(r, &pp); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse payload", nil)
		return
	}
	pp.PlanID = planID

	if err := h.service.LinkProduct(r.Context(), &pp); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusCreated, map[string]any{"planProduct": pp})
}
