package middleware

import (
	"log"
	"net/http"
	posgresql_query "server_service/infrastructure/postgresql/query"
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
		password, ok := claims["password"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Password claim not found or not a string")
			c.Abort()
			return
		}
		storedPassword := posgresql_query.GetAccountPasswordByUsername(username)
		if storedPassword != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			log.Println("Password does not match")
			c.Abort()
			return
		}
		// If the JWT is valid, proceed to the next handler
		c.Next()
	}
}
