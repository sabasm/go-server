package root

import (
	"net/http"

	"github.com/sabasm/go-server/internal/api/handlers"
)

type RootHandler struct{}

func New() *RootHandler {
	return &RootHandler{}
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bw := handlers.NewBufferedResponseWriter(w)
	_, err := bw.Write([]byte("Service running"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	bw.WriteHeader(http.StatusOK)
	if err := bw.Flush(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
