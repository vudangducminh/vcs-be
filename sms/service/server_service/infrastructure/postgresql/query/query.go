package query

import (
	"log"
	"net/http"
	"server_service/entities"
	"server_service/infrastructure/postgresql/connector"
)

func GetAccountPasswordByUsername(username string) (string, int) {
	var account entities.Account
	has, err := connector.PostgreSQL.Table("account").
		Cols("password").
		Alias("account").
		Where("username = ?", username).
		Get(&account)

	if err != nil {
		log.Println("Error retrieving account password:", err)
		return "", http.StatusInternalServerError
	}

	if !has || account.Password == "" {
		log.Println("No account found with username:", username)
		return "", http.StatusNotFound
	}
	return account.Password, http.StatusOK
}
