package query

import (
	"log"
	"net/http"
	"user_service/entities"
	"user_service/infrastructure/postgresql/connector"
)

func GetRoleByUsername(username string) string {
	var account entities.Account
	has, err := connector.PostgreSQL.Table("account").
		Cols("role").
		Alias("account").
		Where("username = ?", username).
		Get(&account)

	if err != nil {
		log.Println("Error retrieving account role:", err)
		return ""
	}

	if !has || account.Role == "" {
		log.Println("No account found with username:", username)
		return ""
	}
	return account.Role
}

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

func GetAccountByUsername(username string) entities.Account {
	var account entities.Account
	has, err := connector.PostgreSQL.Table("account").
		Where("username = ?", username).
		Get(&account)

	if err != nil {
		log.Println("Error retrieving account:", err)
		return entities.Account{}
	}

	if !has {
		log.Println("No account found with username:", username)
		return entities.Account{}
	}
	return account
}

func AddAccountInfo(account entities.Account) int {
	has, err := connector.PostgreSQL.Table("account").
		Where("username = ?", account.Username).
		Count(new(entities.Account))
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if has > 0 {
		log.Println("Account already exists with username:", account.Username)
		return http.StatusConflict
	}
	affected, err := connector.PostgreSQL.Insert(account)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if affected > 0 {
		log.Println("Account added successfully:", account.Username)
		return http.StatusCreated
	} else {
		log.Println("Failed to add account:", account.Username)
		return http.StatusInternalServerError
	}
}

func UpdateAccountInfo(account entities.Account) int {
	affected, err := connector.PostgreSQL.Table("account").
		Where("username = ?", account.Username).
		Update(account)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if affected > 0 {
		log.Println("Account updated successfully:", account.Username)
		return http.StatusCreated
	} else {
		log.Println("Failed to update account:", account.Username)
		return http.StatusInternalServerError
	}
}
