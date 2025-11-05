package init

import (
	swagger "report_service/http"
	es "report_service/infrastructure/elasticsearch/connector"
	psql "report_service/infrastructure/postgresql/connector"
	"report_service/src"
)

func Init() {
	psql.Connect()
	es.Connect()
	go src.DailyReporter()
	swagger.Init()
}
