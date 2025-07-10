package servers_handler

import (
	"sms/object"

	"github.com/gin-gonic/gin"
)

func AddServer(c *gin.Context) {
	// Implementation for adding a server
	var req object.Server
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	c.JSON(200, gin.H{
		"message": "Server added successfully",
	})
}
