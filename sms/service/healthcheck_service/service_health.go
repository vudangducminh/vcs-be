package healthcheck_service

import (
	"net/http"
	redis "sms/server/database/cache/redis/connector"
	elastic "sms/server/database/elasticsearch/connector"
	postgresql "sms/server/database/postgresql/connector"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Timestamp string            `json:"timestamp"`
}

// @Tags         Health
// @Summary      Health check endpoint
// @Description  Check the health status of the application and its dependencies
// @Accept       json
// @Produce      json
// @Success      200 {object} HealthResponse "Application is healthy"
// @Failure      503 {object} HealthResponse "Service unavailable"
// @Router       /health [get]
func ServiceHealthCheck(c *gin.Context) {
	services := make(map[string]string)
	overallStatus := "healthy"

	// Check PostgreSQL connection
	if postgresql.IsConnected() {
		services["postgresql"] = "healthy"
	} else {
		services["postgresql"] = "unhealthy"
		overallStatus = "unhealthy"
	}

	// Check Redis connection
	if redis.IsConnected() {
		services["redis"] = "healthy"
	} else {
		services["redis"] = "unhealthy"
		overallStatus = "unhealthy"
	}

	// Check Elasticsearch connection
	if elastic.IsConnected() {
		services["elasticsearch"] = "healthy"
	} else {
		// Try to reconnect once
		if err := elastic.Reconnect(); err == nil {
			services["elasticsearch"] = "healthy"
		} else {
			services["elasticsearch"] = "unhealthy"
			overallStatus = "unhealthy"
		}
	}

	response := HealthResponse{
		Status:    overallStatus,
		Services:  services,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if overallStatus == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
