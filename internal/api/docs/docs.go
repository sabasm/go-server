package docs

type SwaggerInfo struct{}

// @title Hello World Go API
// @version 1.0
// @description A simple Hello World API built with Go
// @host localhost:8080
// @BasePath /

// @tag.name health
// @tag.description Health check endpoint

// @tag.name root
// @tag.description Root endpoint returning service status

// @Summary Health check endpoint
// @Description Returns OK if the service is healthy
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func HealthCheck() {}

// @Summary Service status endpoint
// @Description Returns the current status of the service
// @Tags root
// @Produce plain
// @Success 200 {string} string "Service running"
// @Router / [get]
func ServiceStatus() {}
