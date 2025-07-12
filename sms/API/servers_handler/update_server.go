package servers_handler

import (
	"net/http"
	"sms/object"
	posgresql_query "sms/server/database/postgresql/query"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Handle updating an existing server
// @Description  Handle updating an existing server by validating input and updating server information
// @Accept       json
// @Produce      json
// @Param        request body object.UpdateServerRequest true "Update server request"
// @Success      200 {object} object.UpdateServerResponse "Server updated"
// @Router       /servers/update_server [put]
func UpdateServer(c *gin.Context) {
	var req object.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.ServerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServerId is required"})
		return
	}

	server, exists := posgresql_query.GetServerById(req.ServerId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	if req.ServerName != "" {
		server.ServerName = req.ServerName
	}
	if req.IPv4 != "" {
		server.IPv4 = req.IPv4
	}
	if req.Status != "" {
		server.Status = req.Status
	}
	server.LastUpdatedTime = time.Now().Format(time.RFC3339)

	status := posgresql_query.UpdateServerInfo(server)
	if status == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"message": "Server updated successfully"})
	} else {
		c.JSON(status, gin.H{"error": "Failed to update server"})
	}
}
