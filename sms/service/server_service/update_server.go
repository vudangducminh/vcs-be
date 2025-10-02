package servers_handler

import (
	"net/http"
	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags         Server
// @Summary      Update an existing server
// @Description  Update an existing server by validating input and updating server information
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body object.UpdateServerRequest true "Update server request"
// @Success      200 {object} object.UpdateServerSuccessResponse "Server updated"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      400 {object} object.UpdateServerBadRequestResponse "Invalid request body"
// @Failure      404 {object} object.UpdateServerStatusNotFoundResponse "Server not found"
// @Failure      409 {object} object.UpdateServerConflictResponse "Server IP already exists"
// @Failure      500 {object} object.UpdateServerInternalServerErrorResponse "Internal server error"
// @Router       /server/update_server [put]
func UpdateServer(c *gin.Context) {
	var req object.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServerId is required"})
		return
	}

	server, exists := elastic_query.GetServerById(req.Id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	if req.ServerName != "" {
		server.ServerName = req.ServerName
	}
	if req.IPv4 != "" {
		if elastic_query.CheckServerExists(req.IPv4) && server.IPv4 != req.IPv4 {
			c.JSON(http.StatusConflict, gin.H{"error": "IPv4 already exists"})
			return
		}
		server.IPv4 = req.IPv4
	}
	server.Status = req.Status
	server.LastUpdatedTime = time.Now().Unix()

	status := elastic_query.UpdateServerInfo(server)
	if status == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server updated successfully",
			"server":  server,
		})
	} else {
		c.JSON(status, gin.H{"error": "Failed to update server"})
	}
}
