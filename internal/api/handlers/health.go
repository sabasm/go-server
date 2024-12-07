package handlers

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
