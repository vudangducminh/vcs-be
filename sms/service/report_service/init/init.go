package init

import (
	swagger "report_service/http"
	es "report_service/infrastructure/elasticsearch/connector"
	psql "report_service/infrastructure/postgresql/connector"
)

func Init() {
	psql.Connect()
	es.Connect()
	swagger.Init()
}
