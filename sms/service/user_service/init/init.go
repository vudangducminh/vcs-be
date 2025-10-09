package init

import (
	swagger "user_service/http"
	psql "user_service/infrastructure/postgresql/connector"
)

func Init() {
	psql.Connect()
	swagger.Init()
}
