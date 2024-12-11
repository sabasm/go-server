package docs

type SwaggerInfo struct {
	Title       string
	Description string
	Version     string
	BasePath    string
}

func HealthCheck() SwaggerInfo {
	return SwaggerInfo{
		Title:       "Health Check Endpoint",
		Description: "Returns service health status",
		Version:     "1.0",
		BasePath:    "/health",
	}
}

func ServiceStatus() SwaggerInfo {
	return SwaggerInfo{
		Title:       "Service Status Endpoint",
		Description: "Returns current service operational status",
		Version:     "1.0",
		BasePath:    "/",
	}
}
