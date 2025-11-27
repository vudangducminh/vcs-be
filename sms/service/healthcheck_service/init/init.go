package init

import (
	swagger "healthcheck_service/http"
	es "healthcheck_service/infrastructure/elasticsearch/connector"
	healthcheck_service "healthcheck_service/src"
	"log"
	"os"
)

func Init() {
	// Must init everything before swagger
	f, err := os.OpenFile("/app/logs/healthcheck_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	es.Connect()
	go healthcheck_service.StartHealthCheck()
	swagger.Init()
}
