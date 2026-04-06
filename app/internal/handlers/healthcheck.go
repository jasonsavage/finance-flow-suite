package handlers

import (
	"net/http"

	"github.com/jasonsavage/financeflow/internal/repository"
)

type HealthcheckHandler struct {
	Repo repository.DatabaseRepo
}

func NewHealthcheckHandler(repo repository.DatabaseRepo) *HealthcheckHandler {
	return &HealthcheckHandler{Repo: repo}
}

// GET /healthcheck
func (h *HealthcheckHandler) Check(w http.ResponseWriter, r *http.Request) {
	err := h.Repo.Ping(r.Context())
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":  "error",
			"message": "database disconnected",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "database connected",
	})
}
