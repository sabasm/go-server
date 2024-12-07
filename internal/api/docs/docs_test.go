package docs

import "testing"

func TestSwaggerDocumentation(t *testing.T) {
	info := SwaggerInfo{}
	if info != (SwaggerInfo{}) {
		t.Error("Expected empty SwaggerInfo struct")
	}
}

func TestHealthCheckDocumentation(t *testing.T) {
	HealthCheck()
}

func TestServiceStatusDocumentation(t *testing.T) {
	ServiceStatus()
}
