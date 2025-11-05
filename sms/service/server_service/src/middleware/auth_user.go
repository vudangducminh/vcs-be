package middleware

import (
	"log"
	"net/http"
	"server_service/src/algorithm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Invalid JWT token:", err)
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT token is not valid")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("JWT claims are not of type jwt.MapClaims")
			c.Abort()
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Username claim not found or not a string")
			c.Abort()
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Role claim not found or not a string")
			c.Abort()
			return
		}
		if role != "admin" && role != "user" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		// Store username and role in context for further handlers
		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}
