package docs

import "testing"

func TestHealthCheckDocumentation(t *testing.T) {
	info := HealthCheck()
	if info.BasePath != "/health" {
		t.Errorf("Expected /health base path, got %s", info.BasePath)
	}
}

func TestServiceStatusDocumentation(t *testing.T) {
	info := ServiceStatus()
	if info.BasePath != "/" {
		t.Errorf("Expected / base path, got %s", info.BasePath)
	}
}
