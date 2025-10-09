package query

import (
	"log"
	"net/http"
	"report_service/entities"
	"report_service/infrastructure/postgresql/connector"
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

func AddEmailInfo(email entities.Email) int {
	has, err := connector.PostgreSQL.Table("email_manager").
		Where("email = ?", email.Email).
		Count(new(entities.Email))
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if has > 0 {
		log.Println("Email already exists:", email.Email)
		return http.StatusConflict
	}
	affected, err := connector.PostgreSQL.Table("email_manager").
		Insert(email)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	if affected > 0 {
		log.Println("Email added successfully:", email.Email)
		return http.StatusCreated
	} else {
		log.Println("Failed to add email:", email.Email)
		return http.StatusInternalServerError
	}
}
func GetAllEmails() ([]entities.Email, int) {
	var emails []entities.Email
	err := connector.PostgreSQL.Table("email_manager").
		Find(&emails)
	if err != nil {
		log.Println("Error retrieving emails:", err)
		return nil, http.StatusInternalServerError
	}
	return emails, http.StatusOK
}
