package server

import (
	"net/http"
	"testing"
)

func TestBuilderFunctionality(t *testing.T) {
	cfg := &Config{Port: 8080, Host: "localhost", BasePath: "/"}

	t.Run("builder configuration", func(t *testing.T) {
		srv := NewBuilder(cfg).
			WithRoute("/test", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}).
			Build()

		if srv == nil {
			t.Error("builder returned nil server")
		}
	})
}
