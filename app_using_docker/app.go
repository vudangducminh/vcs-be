package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func loginPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
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

func handleRegister(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/register_submit", handleRegister)

	// If you have files like static/style.css or static/script.js
	// then you can serve them using the following line.
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
