package test

import (
	"context"    // Used for client methods
	"crypto/tls" // Needed for TLSConfig
	"fmt"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func TestConnectToRedis(t *testing.T) {
	// Best practice: Get sensitive information from environment variables
	// Replace with your actual host, port, and password found in Azure portal
	redisHost := "vcsbe.redis.cache.windows.net"
	redisPort := "6380"
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisPassword := "9rLVW1004AEMuAzqQCBIvb30gHQDqRAzmAzCaLR0JpY="
	redisUsername := "default"

	ctx := context.Background()

	// Configure Redis client options
	options := &redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // The Access Key from Azure Portal
		Username: redisUsername, // Typically "default" for Access Key auth
		DB:       0,             // Default Redis database (usually 0)

		// Since Non-SSL port 6379 is disabled, you MUST connect using SSL/TLS
		TLSConfig: &tls.Config{
			// For connecting to Azure Cache for Redis publicly, an empty TLSConfig is usually sufficient.
			// It enables the client to negotiate TLS.
		},
	}

	// Create the Redis client
	RedisClient = redis.NewClient(options)

	// Test the connection using the PING command
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		t.Errorf("Could not connect to Redis: %v", err)
		return
	}

	log.Printf("Connected to Redis successfully! PING response: %s", pong)

	// --- Example Usage ---

	// Set a key
	err = RedisClient.Set(ctx, "mykey", "hello from go", 10).Err() // 0 means no expiration
	if err != nil {
		t.Errorf("Could not set key: %v", err)
		return
	}
	log.Println("Set 'mykey' to 'hello from go'")

	// Get a key
	val, err := RedisClient.Get(ctx, "mykey").Result()
	if err == redis.Nil {
		log.Println("Key 'mykey' does not exist")
	} else if err != nil {
		t.Errorf("Could not get key: %v", err)
		return
	} else {
		log.Printf("Got 'mykey': %s", val)
	}

	// Delete a key
	// client.Del(ctx, "mykey")
}
