package redis_query

import (
	"context"
	"log"
	"sms/auth"
	"sms/server/database/cache/redis/connector"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGetUsernameByJWTToken(t *testing.T) {
	connector.ConnectToRedis()
	str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIxNjY4MDksInBhc3N3b3JkIjoiMTIzIiwidXNlcm5hbWUiOiJLYXdhaWkifQ.50F8cJjz1tU0pfgwVyNrhnv1snc-wDFg9T86p8pq5jo"
	username, err := connector.RedisClient.Get(context.Background(), str).Result()
	log.Println("Username from JWT token:", username)
	if username != "Kawaii" || err != nil {
		t.Errorf("Expected username 'Kawaii', got '%s'", username)
		return
	}
}

func TestSaveJWTToken(t *testing.T) {
	connector.ConnectToRedis()
	if !connector.IsConnected() {
		t.Error("Failed to connect to Redis")
		return
	}
	// Save the JWT token to Redis
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIxNjU2MzcsInBhc3N3b3JkIjoiMTIzIiwidXNlcm5hbWUiOiJLYXdhaWkifQ.Ao5QaAde80cCdf-9QDmCznRW-h0iFlinIdZNnUTJ-kc"
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
	if username != "Kawaii" {
		t.Error("Wrong username in JWT token, expected 'kawaii', got:", username)
		return
	}
	err = connector.RedisClient.Set(context.Background(), tokenStr, username, 60*time.Second).Err()
	if err != nil {
		t.Error("Failed to save JWT token to Redis:", err)
		return
	}
}
