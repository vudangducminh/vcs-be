package redis_query

import (
	"context"
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
	err = connector.RedisClient.Set(context.Background(), tokenStr, username, time.Duration(second)*time.Second).Err()
	if err != nil {
		log.Println("Failed to save JWT token to Redis:", err)
		return false
	}
	log.Println("JWT token saved to Redis successfully")
	return true
}

func GetUsernameByJWTToken(tokenStr string) string {
	username, err := connector.RedisClient.Get(context.Background(), tokenStr).Result()
	if err != nil {
		log.Println("Error while fetching username using JWT token")
		return ""
	}
	return username
}
