package health

import (
	"net/http"

	"github.com/sabasm/go-server/pkg/utils"
)

type HealthHandler struct{}

func New() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"status": "OK"})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
