package src

import (
	"net/http"
	"server_service/entities"
	elastic_query "server_service/infrastructure/elasticsearch/query"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags         Server
// @Summary      Add a new server
// @Description  Add a new server by validating input and storing server information
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body entities.AddServerRequest true "Add server request"
// @Success      201 {object} entities.AddServerSuccessResponse "Server added"
// @Failure      400 {object} entities.AddServerBadRequestResponse "Invalid request body"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      409 {object} entities.AddServerConflictResponse "Server already exists"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      500 {object} entities.AddServerInternalServerErrorResponse "Internal server error"
// @Router       /servers/add-server [post]
func AddServer(c *gin.Context) {
	// Implementation for adding a server
	var req entities.AddServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.ServerName == "" || req.IPv4 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServerName, and IPv4 are required"})
		return
	}
	if req.Status == "" {
		// Default status if not provided
		req.Status = "active"
	}

	server := entities.Server{
		ServerName:      req.ServerName,
		Status:          req.Status,
		Uptime:          []int{0}, // Default uptime to 0
		CreatedTime:     time.Now().Unix(),
		LastUpdatedTime: time.Now().Unix(),
		IPv4:            req.IPv4,
	}

	status := elastic_query.AddServerInfo(server)
	switch status {
	case http.StatusCreated:
		c.JSON(http.StatusCreated, gin.H{
			"message":     "Server added successfully to Elasticsearch",
			"server_name": server.ServerName,
			"ipv4":        server.IPv4,
		})
	case http.StatusConflict:
		c.JSON(http.StatusConflict, gin.H{"error": "Server already exists with the same IPv4 address in Elasticsearch"})
	default:
		c.JSON(status, gin.H{"error": "Failed to add server into Elasticsearch database"})
	}
}

type ServerAdder interface {
	AddServerInfo(server entities.Server) int
}

var serverAdder ServerAdder

func SetServerAdder(sa ServerAdder) {
	serverAdder = sa
}

func ModifiedAddServer(c *gin.Context) {
	var req entities.AddServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.ServerName == "" || req.IPv4 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServerName, and IPv4 are required"})
		return
	}
	if req.Status == "" {
		req.Status = "active"
	}

	server := entities.Server{
		ServerName:      req.ServerName,
		Status:          req.Status,
		Uptime:          []int{0},
		CreatedTime:     time.Now().Unix(),
		LastUpdatedTime: time.Now().Unix(),
		IPv4:            req.IPv4,
	}

	status := serverAdder.AddServerInfo(server)
	switch status {
	case http.StatusCreated:
		c.JSON(http.StatusCreated, gin.H{
			"message":     "Server added successfully",
			"server_name": server.ServerName,
			"ipv4":        server.IPv4,
		})
	case http.StatusConflict:
		c.JSON(http.StatusConflict, gin.H{"error": "Server already exists with the same IPv4 address"})
	default:
		c.JSON(status, gin.H{"error": "Failed to add server"})
	}
}
