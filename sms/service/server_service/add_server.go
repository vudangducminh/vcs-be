package servers_handler

import (
	"net/http"
	"sms/object"
	"time"

	elastic_query "sms/server/database/elasticsearch/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Handle adding a new server
// @Description  Handle adding a new server by validating input and storing server information
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body object.AddServerRequest true "Add server request"
// @Success      201 {object} object.AddServerSuccessResponse "Server added"
// @Failure      400 {object} object.AddServerBadRequestResponse "Invalid request body"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      409 {object} object.AddServerConflictResponse "Server already exists"
// @Failure      500 {object} object.AddServerInternalServerErrorResponse "Internal server error"
// @Router       /servers/add_server [post]
func AddServer(c *gin.Context) {
	// Implementation for adding a server
	var req object.AddServerRequest
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

	server := object.Server{
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
