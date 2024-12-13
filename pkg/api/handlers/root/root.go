package root

import (
	"encoding/json"
	"net/http"
)

type RootResponse struct {
	Message string `json:"message"`
}

type RootHandler struct{}

func New() *RootHandler {
	return &RootHandler{}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := RootResponse{Message: "Welcome to Go Server"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
