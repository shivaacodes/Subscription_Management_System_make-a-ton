package reporting

import (
	"net/http"

	"github.com/subflow/backend/internal/platform"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	metrics, err := h.service.GetDashboard(r.Context())
	if err != nil {
		platform.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to calculate revenue metrics", nil)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"metrics": metrics})
}
