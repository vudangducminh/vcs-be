package servers_handler

import (
	"net/http"
	"sms/algorithm"
	"sms/object"
	"time"

	elastic_query "sms/server/database/elasticsearch/query"
	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Handle adding a new server
// @Description  Handle adding a new server by validating input and storing server information
// @Accept       json
// @Produce      json
// @Param        request body object.AddServerRequest true "Add server request"
// @Success      201 {object} object.AddServerSuccessResponse "Server added"
// @Failure      400 {object} object.AddServerBadRequestResponse "Invalid request body"
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

	// Set current time as CreatedTime
	req.CreatedTime = time.Now().Format(time.RFC3339)
	req.LastUpdatedTime = req.CreatedTime
	server := object.Server{
		ServerId:        algorithm.SHA256Hash(time.Now().String()),
		ServerName:      req.ServerName,
		Status:          req.Status,
		Uptime:          0, // Default uptime to 0
		CreatedTime:     req.CreatedTime,
		LastUpdatedTime: req.LastUpdatedTime,
		IPv4:            req.IPv4,
	}

	status := posgresql_query.AddServerInfo(server)
	switch status {
	case http.StatusCreated:
		c.JSON(http.StatusCreated, gin.H{
			"message":     "Server added successfully to PostgreSQL",
			"server_id":   server.ServerId,
			"server_name": server.ServerName,
			"ipv4":        server.IPv4,
		})
	case http.StatusConflict:
		c.JSON(http.StatusConflict, gin.H{"error": "Server already exists with the same IPv4 address in PostgreSQL"})
	default:
		c.JSON(status, gin.H{"error": "Failed to add server into PostgreSQL database"})
	}
	status = elastic_query.AddServerInfo(server)
	switch status {
	case http.StatusCreated:
		c.JSON(http.StatusCreated, gin.H{
			"message":     "Server added successfully to Elasticsearch",
			"server_id":   server.ServerId,
			"server_name": server.ServerName,
			"ipv4":        server.IPv4,
		})
	case http.StatusConflict:
		c.JSON(http.StatusConflict, gin.H{"error": "Server already exists with the same IPv4 address in Elasticsearch"})
	default:
		c.JSON(status, gin.H{"error": "Failed to add server into Elasticsearch database"})
	}
}
