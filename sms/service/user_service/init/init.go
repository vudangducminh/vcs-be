package init

import (
	"log"
	"os"
	swagger "user_service/http"
	psql "user_service/infrastructure/postgresql/connector"
)

func Init() {
	f, err := os.OpenFile("/app/logs/user_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	psql.Connect()
	swagger.Init()
}
