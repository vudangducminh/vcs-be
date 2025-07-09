package test

import (
	"context"
	"sms/auth"
	"sms/server/database/cache/redis/connector"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSaveJWTToken(t *testing.T) {
	connector.ConnectToRedis()
	if !connector.IsConnected() {
		t.Error("Failed to connect to Redis")
		return
	}
	// Save the JWT token to Redis
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIwMzg0ODQsInBhc3N3b3JkIjoiYWEiLCJ1c2VybmFtZSI6ImFhIn0.8HJTSwaR2AMRsT4ZwnWyJDm_ot_DrgmlwXqPRhZyu8M"
	token, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		t.Error("Invalid JWT token:", err)
		return
	}
	if !token.Valid {
		t.Error("JWT token is not valid")
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	username, ok := claims["username"].(string)
	if !ok {
		t.Error("Username claim not found or not a string")
		return
	}
	if username != "aa" {
		t.Error("Wrong username in JWT token, expected 'aa', got:", username)
		return
	}
	err = connector.RedisClient.Set(context.Background(), tokenStr, username, 10*time.Second).Err()
	if err != nil {
		t.Error("Failed to save JWT token to Redis:", err)
		return
	}
}
