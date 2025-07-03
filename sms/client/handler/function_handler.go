package handler

import (
	"log"
	"net/http"
	"sms/auth"
	posgresql_query "sms/database/postgresql/query"
	"sms/object"
	"text/template"
	"time"
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

	// Generate JWT token before redirecting the user
	tokenString, err := auth.GenerateJWT(username, password)
	if err != nil {
		// handle error
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true, // prevents JS access (important)
		// Secure: true, // enable if using HTTPS
		Expires: time.Now().Add(time.Hour), // optional, set expiry
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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

func HandleAddServer(w http.ResponseWriter, r *http.Request) {
	log.Fatalln("Error: HandleAddServer function is not implemented yet")
	serverName := r.FormValue("server_name")
	serverAddress := r.FormValue("server_address")
	serverPort := r.FormValue("server_port")
	log.Println("Adding server:", serverName, serverAddress, serverPort)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	if serverName == "" || serverAddress == "" || serverPort == "" {
		log.Println("All fields are required")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Here you would typically save the server information to the database
	log.Printf("Adding server: %s at %s:%s\n", serverName, serverAddress, serverPort)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
