package main

import (
	"net/http"
	"testing"
	"time"
)

func TestMainServer(t *testing.T) {
	go main()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
}
