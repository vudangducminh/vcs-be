package test

import (
	"context"
	"log"
	"sms/server/database/cache/redis/connector"
	"testing"
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
