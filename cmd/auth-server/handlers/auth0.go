package handlers

import (
	"net/http"

	"github.com/sabasm/go-server/internal/utils"
)

type ValidationResponse struct {
	Valid bool `json:"valid"`
}

type Auth0Handler struct{}

func NewAuth0Handler() *Auth0Handler {
	return &Auth0Handler{}
}

func (h *Auth0Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := ValidationResponse{Valid: true}
	if err := utils.WriteJSONResponse(w, http.StatusOK, response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
