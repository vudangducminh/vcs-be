package website_handler

import (
	"log"
	"net/http"
	"sms/object"

	redis_query "sms/server/database/cache/redis/query"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	var req object.AuthRequest
	username := redis_query.GetUsernameByJWTToken(req.JWT)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Println("Current username:", username)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Authentication successfully",
		"username": username,
	})

}
