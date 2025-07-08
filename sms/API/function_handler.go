package API

// import (
// 	"log"
// 	"net/http"
// 	"sms/auth"
// 	"sms/object"
// 	posgresql_query "sms/server/database/postgresql/query"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// )

// // var templates = template.Must(template.ParseGlob("client/templates/*.html"))

// func DashboardPage(w http.ResponseWriter, r *http.Request) {
// 	cookie, err := r.Cookie("token")
// 	if err != nil {
// 		// Cookie not found or other error
// 		log.Println("Token cookie not found:", err)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
// 	// Maybe I can use redis to save the token so that I can retrieve the user information without the need to validate the cookie
// 	token, err := auth.ValidateJWT(cookie.Value)
// 	if err != nil {
// 		log.Println("Invalid token:", err)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
// 	claims := token.Claims.(jwt.MapClaims)
// 	username, ok := claims["username"].(string)
// 	if !ok {
// 		log.Println("Username claim not found or not a string")
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
// 	password, ok := claims["password"].(string)
// 	if !ok {
// 		log.Println("Password claim not found or not a string")
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
// 	log.Println("Current username:", username)
// 	log.Println("Current password:", password)
// 	// templates.ExecuteTemplate(w, "dashboard.html", nil)
// }

// // HandleLogin handles the login request and validates user credentials
// // @Summary Handle user login
// // @Description Validates user credentials and generates a JWT token
// // @Tags Users
// // @Success 200 {object} object.LoginRequest
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/v1/users/login [post]
// // @Accept json
// // @Produce json
// // @Param loginRequest body object.LoginRequest true "Login request body"

// func HandleLogin(c *gin.Context) {
// 	var req object.LoginRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	if req.Username == "" || req.Password == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Username and password are required",
// 			"error":   "Invalid credentials",
// 		})
// 		return
// 	}

// 	log.Println("Username:", req.Username)
// 	log.Println("Password:", req.Password)

// 	storedPassword := posgresql_query.GetAccountPasswordByUsername(req.Username)
// 	if storedPassword != req.Password {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Invalid username or password",
// 			"error":   "Invalid credentials",
// 		})
// 		return
// 	}

// 	// Generate JWT token before redirecting the user
// 	tokenString, err := auth.GenerateJWT(req.Username, req.Password)
// 	if err != nil {
// 		// handle error
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Error generating token",
// 			"error":   "Error generating token",
// 		})
// 		return
// 	}
// 	// Set the cookie in the response

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Login successful",
// 		"token":   tokenString,
// 	})
// }

// func HandleRegister(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/register", http.StatusSeeOther)
// 		return
// 	}

// 	fullName := r.FormValue("fullname")
// 	email := r.FormValue("email")
// 	username := r.FormValue("username")
// 	password := r.FormValue("password")
// 	confirmPassword := r.FormValue("confirm_password")

// 	if fullName == "" || email == "" || username == "" || password == "" || confirmPassword == "" {
// 		log.Println(w, "All fields are required", http.StatusBadRequest)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}
// 	if password != confirmPassword {
// 		log.Println(w, "Passwords do not match", http.StatusBadRequest)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}

// 	account := object.Account{
// 		Fullname: fullName,
// 		Email:    email,
// 		Username: username,
// 		Password: password,
// 	}

// 	log.Println("Registering Username:", username)
// 	log.Println("Password:", password)

// 	if posgresql_query.AddAccountInfo(account) {
// 		// Add success text
// 	} else {
// 		// Add failed text
// 	}

// 	http.Redirect(w, r, "/register", http.StatusSeeOther)
// }

// func HandleAddServer(w http.ResponseWriter, r *http.Request) {
// 	serverName := r.FormValue("server_name")
// 	serverIP := r.FormValue("server_IP")
// 	serverType := r.FormValue("server_type")
// 	log.Println("Adding server:", serverName, serverIP, serverType)

// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
// 		return
// 	}

// 	log.Printf("Adding server: %s at %s:%s\n", serverName, serverIP, serverType)

// 	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
// }
