package handler

import (
	"app_using_docker/elasticsearch/query"
	"app_using_docker/object"
	"log"
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Println("Username:", username)
	log.Println("Password:", password)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	account := object.Account{
		ID:       0,
		FullName: fullName,
		Email:    email,
		Username: username,
		Password: password,
	}

	log.Println("Registering Username:", username)
	log.Println("Password:", password)

	query.AddAccountInfo(account)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
