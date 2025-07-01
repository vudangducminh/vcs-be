package handler

import (
	"log"
	"net/http"
	posgresql_query "sms/database/postgresql/query"
	"sms/object"
	"text/template"
)

var templates = template.Must(template.ParseGlob("client/templates/*.html"))

func LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func DashboardPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		log.Println("Username and password are required")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println("Username:", username)
	log.Println("Password:", password)

	storedPassword := posgresql_query.GetAccountPasswordByUsername(username)
	if storedPassword != password {
		log.Println("Invalid username or password")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		ID:       0,
		FullName: fullName,
		Email:    email,
		Username: username,
		Password: password,
	}

	log.Println("Registering Username:", username)
	log.Println("Password:", password)

	posgresql_query.AddAccountInfo(account)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
