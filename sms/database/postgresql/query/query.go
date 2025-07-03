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

func AddAccountInfo(account object.Account) bool {
	has, err := connector.Engine.Table("account").
		Where("username = ?", account.Username).Count(new(object.Account))
	if err != nil {
		log.Println(err)
		return false
	}
	if has > 0 {
		log.Println("Account already exists with username:", account.Username)
		return false
	}
	affected, err := connector.Engine.Insert(account)
	if err != nil {
		log.Println(err)
		return false
	}
	if affected > 0 {
		log.Println("Account added successfully:", account.Username)
		return true
	} else {
		log.Println("Failed to add account:", account.Username)
		return false
	}
}
