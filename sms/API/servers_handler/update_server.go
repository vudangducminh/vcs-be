package servers_handler

import (
	"net/http"
	"sms/object"
	redis_query "sms/server/database/cache/redis/query"
	elastic_query "sms/server/database/elasticsearch/query"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Handle updating an existing server
// @Description  Handle updating an existing server by validating input and updating server information
// @Accept       json
// @Produce      json
// @Param        request body object.UpdateServerRequest true "Update server request"
// @Success      200 {object} object.UpdateServerSuccessResponse "Server updated"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      400 {object} object.UpdateServerBadRequestResponse "Invalid request body"
// @Failure      404 {object} object.UpdateServerStatusNotFoundResponse "Server not found"
// @Failure      409 {object} object.UpdateServerConflictResponse "Server already exists"
// @Failure      500 {object} object.UpdateServerInternalServerErrorResponse "Internal server error"
// @Router       /servers/update_server [put]
func UpdateServer(c *gin.Context) {
	var req object.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username := redis_query.GetUsernameByJWTToken(req.JWT)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	if req.ServerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServerId is required"})
		return
	}

	server, exists := elastic_query.GetServerById(req.ServerId)
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
		if server.Status == "active" && (req.Status == "inactive" || req.Status == "maintenance") {
			server.Uptime = time.Now().Unix() - server.LastUpdatedTime
		}
		server.Status = req.Status
	}
	server.LastUpdatedTime = time.Now().Unix()

	status := elastic_query.UpdateServerInfo(server)
	if status == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"message": "Server updated successfully"})
	} else {
		c.JSON(status, gin.H{"error": "Failed to update server"})
	}
}
