package server

import (
	"net/http"
	"testing"
	"time"
)

func TestSetHandler(t *testing.T) {
	srv := &Server{
		srv: &http.Server{
			ReadHeaderTimeout: 5 * time.Second, // Fix for G112 Slowloris attack protection
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv.SetHandler(handler)

	if srv.srv.Handler == nil {
		t.Fatal("expected handler to be set")
	}
}
