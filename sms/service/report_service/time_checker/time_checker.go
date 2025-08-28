package notification

import (
	redis_query "sms/server/database/cache/redis/query"
	elastic_query "sms/server/database/elasticsearch/query"
	template "sms/service/report_service/template"
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

func CalculateAverageServerUptime() float32 {
	var totalServerUptime int64 = 0
	now := time.Now().Unix()
	time, count := elastic_query.GetTotalUptime()
	totalServerUptime += time
	time, count = elastic_query.GetTotalLastUpdatedTime()
	// log.Println("Total Last Updated Time: ", time, " Count: ", count)
	// log.Println("Now: ", now)
	totalServerUptime += now*int64(count) - time
	time, count = elastic_query.GetTotalCreatedTime()
	var maxServerUptime int64 = now*int64(count) - time
	// log.Println("Total Server Uptime: ", totalServerUptime, " Max Server Uptime: ", maxServerUptime)
	if maxServerUptime == 0 {
		return 0.0
	}
	return float32(totalServerUptime) / float32(maxServerUptime) * 100
}
