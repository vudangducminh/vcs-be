package users_handler

import (
	"log"
	"net/http"
	"sms/auth"
	"sms/object"
	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Summary      Handle user Login
// @Description  (Description can be commented out too if needed)
// @Tags         users
// @Router       /users/login [post]

func HandleLogin(c *gin.Context) {
	var req object.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username and password are required",
			"error":   "Invalid credentials",
		})
		return
	}

	log.Println("Username:", req.Username)
	log.Println("Password:", req.Password)

	storedPassword := posgresql_query.GetAccountPasswordByUsername(req.Username)
	if storedPassword != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid username or password",
			"error":   "Invalid credentials",
		})
		return
	}

	// Generate JWT token before redirecting the user
	tokenString, err := auth.GenerateJWT(req.Username, req.Password)
	if err != nil {
		// handle error
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error generating token",
			"error":   "Error generating token",
		})
		return
	}
	// Set the cookie in the response

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}
