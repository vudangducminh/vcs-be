package connector

import (
	// Used for client methods
	// Needed for TLSConfig

	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

var isConnected = false

func IsConnected() bool {
	return isConnected
}

func ConnectToRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-13990.c251.east-us-mz.azure.redns.redis-cloud.com:13990",
		Username: "default",
		Password: "c3QOPfZSpiPqTmfBCINbFuPvaKUSEMM8",
		DB:       0,
	})
	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	isConnected = true
	log.Println("Connected to Redis successfully")
}
