package redis_query

import (
	"log"
	"sms/auth"
	"sms/server/database/cache/redis/connector"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SaveJWTToken(tokenStr string, second int) bool {
	// Save the JWT token to Redis
	token, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		log.Println("Invalid JWT token:", err)
		return false
	}
	if !token.Valid {
		log.Println("JWT token is not valid")
		return false
	}
	claims := token.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok {
		log.Println("Username claim not found or not a string")
		return false
	}
	err = connector.RedisClient.Set(connector.RedisClient.Context(), tokenStr, username, second*time.Second).Err()
	if err != nil {
		log.Println("Failed to save JWT token to Redis:", err)
		return false
	}
	log.Println("JWT token saved to Redis successfully")
	return true
}
