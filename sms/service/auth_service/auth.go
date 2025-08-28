package middleware

import (
	"net/http"

	redis_query "sms/server/database/cache/redis/query"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt := c.GetHeader("jwt")
		username := redis_query.GetUsernameByJWTToken(jwt)
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}

		// If the JWT is valid, proceed to the next handler
		c.Next()
	}
}
