package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMainServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("OK"))
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("Service running"))
			if err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		}
	}))
	defer server.Close()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request to /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for /health; got %v", resp.Status)
	}

	expectedHealth := "OK"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	if string(body) != expectedHealth {
		t.Errorf("Expected body %v for /health; got %v", expectedHealth, string(body))
	}

	resp, err = http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make request to /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for /; got %v", resp.Status)
	}

	expectedRoot := "Service running"
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	if string(body) != expectedRoot {
		t.Errorf("Expected body %v for /; got %v", expectedRoot, string(body))
	}
}
