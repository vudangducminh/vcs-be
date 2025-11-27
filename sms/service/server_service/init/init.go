package init

import (
	"log"
	"os"
	swagger "server_service/http"
	es "server_service/infrastructure/elasticsearch/connector"
)

func Init() {
	f, err := os.OpenFile("/app/logs/server_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	es.Connect()
	swagger.Init()
}
