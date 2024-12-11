package health

import (
	"net/http"

	"github.com/sabasm/go-server/internal/api/handlers"
)

type HealthHandler struct{}

func New() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bw := handlers.NewBufferedResponseWriter(w)
	_, err := bw.Write([]byte("OK"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	bw.WriteHeader(http.StatusOK)
	if err := bw.Flush(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
