package init

import (
	swagger "server_service/http"
	es "server_service/infrastructure/elasticsearch/connector"
)

func Init() {
	es.Connect()
	swagger.Init()
}
