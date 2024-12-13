package health

import (
        "encoding/json"
        "net/http"
)

type HealthHandler struct{}

func New() *HealthHandler {
        return &HealthHandler{}
}

func (h *HealthHandler) ValidateRequest(r *http.Request) error {
        return nil
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        response := struct {
                Status string `json:"status"`
        }{
                Status: "OK",
        }
        if err := json.NewEncoder(w).Encode(response); err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
}


