package notification

import (
	redis_query "sms/server/database/cache/redis/query"
	"time"
)

func TimeChecker() {
	var secs int64 = 86400
	for {
		var currentTime = int64(time.Now().Unix())
		// Check if it's time to send daily report emails
		var mod int = int(min(currentTime%secs, secs-currentTime%secs))
		if mod < 10 {
			redis_query.SendDailyReportEmail()
		}
		time.Sleep(10 * time.Second)
	}
}
