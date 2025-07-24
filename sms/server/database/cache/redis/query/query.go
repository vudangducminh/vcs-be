package redis_query

import (
	"context"
	"log"
	"sms/auth"
	template "sms/notification/template"
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
	log.Println("JWT token:", tokenStr)
	username, err := connector.RedisClient.Get(context.Background(), tokenStr).Result()
	log.Println("Username from JWT token:", username)
	if err != nil {
		log.Println("Error while fetching username using JWT token")
		return ""
	}
	return username
}

func SaveDailyReportEmailRequest(email string, second int, startTime int64) bool {
	// Save the daily report email request to Redis
	err := connector.RedisClient.Set(context.Background(), "Email: "+email, startTime, time.Duration(second)*time.Second).Err()
	if err != nil {
		log.Println("Failed to save daily report email request to Redis:", err)
		return false
	}
	log.Println("Daily report email request saved to Redis successfully")
	return true
}

func SendDailyReportEmail() {
	var cursor uint64
	var allKeys []string
	ctx := context.Background()
	pattern := "Email:*"
	var to []string

	for {
		// Scan for a batch of keys matching the pattern.
		// The '10' is a hint to Redis about how many keys to check per iteration.
		keys, nextCursor, err := connector.RedisClient.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			log.Printf("Error during Redis SCAN for pattern '%s': %v", pattern, err)
			return
		}

		// Add the found keys to our list.
		allKeys = append(allKeys, keys...)

		// When the cursor returned by SCAN is 0, the iteration is complete.
		if nextCursor == 0 {
			break
		}

		// Update the cursor for the next iteration.
		cursor = nextCursor
	}

	log.Printf("Found %d email requests to process.", len(allKeys))

	for _, key := range allKeys {
		value, err := connector.RedisClient.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get value for key %s: %v", key, err)
			continue
		}
		var tm int64
		for _, char := range value {
			tm = tm*10 + int64(char-'0')
		}
		if tm > time.Now().Unix() {
			log.Printf("Skipping %s as the start time is in the future.", key)
			continue
		}
		to = append(to, key[7:])
	}

	template.SendEmail(to)
	log.Println("Daily report email sent successfully to:", to)
}
