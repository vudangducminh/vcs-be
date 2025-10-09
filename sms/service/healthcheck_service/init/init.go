package init

import (
	swagger "healthcheck_service/http"
	es "healthcheck_service/infrastructure/elasticsearch/connector"
	healthcheck_service "healthcheck_service/src"
)

func Init() {
	// Must init everything before swagger
	es.Connect()
	go healthcheck_service.HealthCheck()
	swagger.Init()
}
