package posgresql_query

import (
	"log"
	"sms/database/postgresql/connector"
	"sms/object"
)

func GetAccountPasswordByUsername(username string) string {
	var account object.Account
	has, err := connector.Engine.Table("account").Cols("password").Alias("account").
		Where("username = ?", username).
		Get(&account)

	if err != nil {
		log.Println("Error retrieving account password:", err)
		return ""
	}

	if !has || account.Password == "" {
		log.Println("No account found with username:", username)
		return ""
	}
	return account.Password
}

func AddAccountInfo(account object.Account) {
	// This function should connect to the PostgreSQL database and add the account information.
	// For now, we will just print the account information for demonstration purposes.
	println("Adding account:", account.Username, "with password:", account.Password)
}
