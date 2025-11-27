package init

import (
	"log"
	"os"
	swagger "report_service/http"
	es "report_service/infrastructure/elasticsearch/connector"
	psql "report_service/infrastructure/postgresql/connector"
	"report_service/src"
)

func Init() {
	f, err := os.OpenFile("/app/logs/report_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	psql.Connect()
	es.Connect()
	go src.DailyReporter()
	swagger.Init()
}
