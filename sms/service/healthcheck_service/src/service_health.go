package healthcheck_service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Timestamp string            `json:"timestamp"`
}

// checkServiceHealth makes an HTTP request to check if a service is healthy
func checkServiceHealth(serviceURL string, timeout time.Duration) string {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(serviceURL)
	if err != nil {
		return "unhealthy"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "healthy"
	}
	return "unhealthy"
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

	resultsChan := make(chan bool, 3)

	go func() {
		userServiceStatus := checkServiceHealth("http://user_service:8800", 5*time.Second)
		services["user_service"] = userServiceStatus
		if userServiceStatus != "healthy" {
			overallStatus = "unhealthy"
			log.Println("User Service is unhealthy")
		}
		resultsChan <- userServiceStatus == "healthy"
	}()

	go func() {
		serverServiceStatus := checkServiceHealth("http://server_service:8801", 5*time.Second)
		services["server_service"] = serverServiceStatus
		if serverServiceStatus != "healthy" {
			overallStatus = "unhealthy"
			log.Println("Server Service is unhealthy")
		}
		resultsChan <- serverServiceStatus == "healthy"
	}()
	go func() {
		reportServiceStatus := checkServiceHealth("http://report_service:8802", 5*time.Second)
		services["report_service"] = reportServiceStatus
		if reportServiceStatus != "healthy" {
			overallStatus = "unhealthy"
			log.Println("Report Service is unhealthy")
		}
		resultsChan <- reportServiceStatus == "healthy"
	}()
	// Wait for all goroutines to finish
	for i := 0; i < 3; i++ {
		isAlive := <-resultsChan
		if !isAlive {
			overallStatus = "unhealthy"
		}
	}
	close(resultsChan)

	response := HealthResponse{
		Status:    overallStatus,
		Services:  services,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
