package src

import (
	"net/http"
	"time"
	psql "user_service/infrastructure/postgresql/connector"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// @Tags         Health
// @Summary      Health check endpoint
// @Description  Check the health status of the user service and its dependencies
// @Accept       json
// @Produce      json
// @Success      200 {object} HealthResponse "Service is healthy"
// @Failure      503 {object} HealthResponse "Service unavailable"
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	services := make(map[string]string)
	overallStatus := "healthy"

	// Check PostgreSQL connection
	if !psql.IsConnected() {
		overallStatus = "unhealthy"
		services["postgresql"] = "unhealthy"
	}

	// Prepare response
	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Services:  services,
		Version:   "1.0.0",
	}

	// Return appropriate HTTP status
	if overallStatus == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
