package connector

import (
	"context"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestConnectToRedis(t *testing.T) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-13990.c251.east-us-mz.azure.redns.redis-cloud.com:13990",
		Username: "default",
		Password: "c3QOPfZSpiPqTmfBCINbFuPvaKUSEMM8",
		DB:       0,
	})

	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		t.Fatalf("Failed to set value in Redis: %v", err)
	}

	result, err := rdb.Get(ctx, "foo").Result()

	if err != nil {
		t.Fatalf("Failed to get value from Redis: %v", err)
	}

	if result != "bar" {
		t.Errorf("expected 'bar', got '%s'", result)
	}

	log.Println("Successfully connected to Redis, set and retrieved value.")
}
