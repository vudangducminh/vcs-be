package middleware

import (
	"log"
	"net/http"
	"sms/algorithm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtString := c.GetHeader("jwt")
		token, err := algorithm.ValidateJWT(jwtString)
		if err != nil {
			log.Println("Invalid JWT token:", err)
			return
		}
		if !token.Valid {
			log.Println("JWT token is not valid")
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("JWT claims are not of type jwt.MapClaims")
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			log.Println("Username claim not found or not a string")
			return
		}
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authenticati	on failed"})
			c.Abort()
			return
		}
		// If the JWT is valid, proceed to the next handler
		c.Next()
	}
}
