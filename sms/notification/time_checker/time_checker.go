package notification

import (
	template "sms/notification/template"
	redis_query "sms/server/database/cache/redis/query"
	elastic_query "sms/server/database/elasticsearch/query"
	"time"
)

func TimeCheckerForSendingEmails() {
	var secs int64 = 86400
	for {
		var currentTime = int64(time.Now().Unix())
		// Check if it's time to send daily report emails
		var mod int = int(currentTime % secs)
		if mod < 3 {
			template.SendEmail(redis_query.GetEmailListForReporting(), CalculateAverageServerUptime())
		}
		time.Sleep(3 * time.Second)
	}
}

var totalServerUptime int64 = 0
var maxServerUptime int64 = 0

func CalculateAverageServerUptime() float32 {
	return float32(totalServerUptime) / float32(maxServerUptime) * 100
}

func CheckServerUptime() {
	for {
		maxServerUptime += int64(elastic_query.GetTotalServersCount())
		totalServerUptime += int64(elastic_query.GetTotalActiveServersCount())
		time.Sleep(10 * time.Second)
	}
}
