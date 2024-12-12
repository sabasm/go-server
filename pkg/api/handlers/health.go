package handlers

import (
	"net/http"

	"github.com/sabasm/go-server/pkg/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"status": "OK"}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
