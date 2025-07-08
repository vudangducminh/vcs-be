package users_handler

import (
	"log"
	"net/http"
	"sms/object"
	posgresql_query "sms/server/database/postgresql/query"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	fullName := r.FormValue("fullname")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if fullName == "" || email == "" || username == "" || password == "" || confirmPassword == "" {
		log.Println(w, "All fields are required", http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if password != confirmPassword {
		log.Println(w, "Passwords do not match", http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	account := object.Account{
		Fullname: fullName,
		Email:    email,
		Username: username,
		Password: password,
	}

	log.Println("Registering Username:", username)
	log.Println("Password:", password)

	if posgresql_query.AddAccountInfo(account) {
		// Add success text
	} else {
		// Add failed text
	}

	http.Redirect(w, r, "/register", http.StatusSeeOther)
}
