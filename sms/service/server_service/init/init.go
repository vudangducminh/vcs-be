package init

import (
	swagger "server_service/http"
	es "server_service/infrastructure/elasticsearch/connector"
	psql "server_service/infrastructure/postgresql/connector"
)

func Init() {
	psql.Connect()
	es.Connect()
	swagger.Init()
}
