package catalog

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

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"
	products, err := h.service.ListProducts(r.Context(), activeOnly)
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve products", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"products": products})
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := platform.ReadJSON(r, &p); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse product payload", nil)
		return
	}

	if err := h.service.CreateProduct(r.Context(), &p); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{"product": p})
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			platform.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Product not found", nil)
			return
		}
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed retrieving product", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"product": p})
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p Product
	if err := platform.ReadJSON(r, &p); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse product payload", nil)
		return
	}
	p.ID = id

	if err := h.service.UpdateProduct(r.Context(), &p); err != nil {
		if errors.Is(err, ErrProductNotFound) {
			platform.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Product not found", nil)
			return
		}
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusOK, map[string]any{"product": p})
}

func (h *Handler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
	var v ProductVariant
	if err := platform.ReadJSON(r, &v); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse variant payload", nil)
		return
	}
	v.ProductID = productID

	if err := h.service.CreateVariant(r.Context(), &v); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{"variant": v})
}

func (h *Handler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var v ProductVariant
	if err := platform.ReadJSON(r, &v); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse variant payload", nil)
		return
	}
	v.ID = id

	if err := h.service.UpdateVariant(r.Context(), &v); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	platform.WriteJSON(w, http.StatusOK, map[string]any{"variant": v})
}

func (h *Handler) ListTaxes(w http.ResponseWriter, r *http.Request) {
	taxes, err := h.service.ListTaxRules(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve taxes", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"taxes": taxes})
}

func (h *Handler) CreateTax(w http.ResponseWriter, r *http.Request) {
	var t TaxRule
	if err := platform.ReadJSON(r, &t); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse tax payload", nil)
		return
	}

	if err := h.service.CreateTaxRule(r.Context(), &t); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusCreated, map[string]any{"tax": t})
}

func (h *Handler) ListDiscounts(w http.ResponseWriter, r *http.Request) {
	discounts, err := h.service.ListDiscountRules(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve discounts", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"discounts": discounts})
}

func (h *Handler) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var d DiscountRule
	if err := platform.ReadJSON(r, &d); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse discount payload", nil)
		return
	}

	if err := h.service.CreateDiscountRule(r.Context(), &d); err != nil {
		platform.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	platform.WriteJSON(w, http.StatusCreated, map[string]any{"discount": d})
}
