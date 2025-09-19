package report_service

import (
	"log"
	"time"
)

var emailList []string

func Reporter() {
	time.Sleep(10 * time.Second) // Initial delay to allow other services to start

	for {
		if len(emailList) == 0 {
			log.Println("Attempting to fetch email list")
		}
		time.Sleep(1 * time.Minute) // Check every minute
	}
}
