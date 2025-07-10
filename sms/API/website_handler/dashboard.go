package website_handler

import (
	"log"
	"net/http"
	"sms/object"

	redis_query "sms/server/database/cache/redis/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Auth
// @Summary      Handle user authentication
// @Description  Handle user authentication by validating JWT token and returning username
// @Accept       json
// @Produce      json
// @Param        request body object.AuthRequest true "Authentication request"
// @Success      200 {object} object.AuthResponse "Authentication successfully"
// @Router       /auth [post]
func Authentication(c *gin.Context) {
	var req object.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

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
