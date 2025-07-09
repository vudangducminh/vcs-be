package test

import (
	"context" // Used for client methods
	"fmt"

	// Needed for TLSConfig

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

	rdb.Set(ctx, "foo", "bar", 0)
	result, err := rdb.Get(ctx, "foo").Result()

	if err != nil {
		t.Fatalf("Failed to get value from Redis: %v", err)
		return
	}

	fmt.Println(result) // >>> bar
}
