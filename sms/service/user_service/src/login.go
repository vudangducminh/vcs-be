package src

import (
	"log"
	"net/http"
	"user_service/entities"
	"user_service/src/algorithm"

	posgresql_query "user_service/infrastructure/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         User
// @Summary      Handle user login
// @Description  Handle user login by validating credentials and generating a JWT token
// @Accept       json
// @Produce      json
// @Param        request body entities.LoginRequest true "Login request"
// @Success      200 {object} entities.LoginSuccessResponse "Login successful"
// @Failure      400 {object} entities.LoginBadRequestResponse "Invalid request body"
// @Failure      401 {object} entities.LoginUnauthorizedResponse "Invalid credentials"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      500 {object} entities.LoginInternalServerErrorResponse "Error generating token
// @Router       /users/login [post]
func Login(c *gin.Context) {
	var req entities.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Println("Username:", req.Username)
	log.Println("Password:", req.Password)
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Username and password are required",
			"error":   "Invalid credentials",
		})
		return
	}
	req.Password = algorithm.SHA256Hash(req.Password)
	storedPassword, status := posgresql_query.GetAccountPasswordByUsername(req.Username)
	if status == http.StatusInternalServerError {
		c.JSON(status, gin.H{"error": "Error retrieving stored password"})
		return
	}
	log.Println("Username:", req.Username)
	log.Println("Password:", req.Password)
	log.Println("Stored Password:", storedPassword)
	if storedPassword != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid username or password",
			"error":   "Invalid credentials",
		})
		return
	}

	// Generate JWT token before redirecting the user
	role := posgresql_query.GetRoleByUsername(req.Username)
	tokenString, err := algorithm.GenerateJWT(req.Username, role)
	if err != nil {
		// handle error
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error generating token",
			"error":   "Error generating token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"role":    role,
	})
}

type AccountQuery interface {
	GetAccountPasswordByUsername(username string) (string, int)
	GetRoleByUsername(username string) string
}

type JWTGenerator interface {
	GenerateJWT(username, role string) (string, error)
}

var (
	accountQuery AccountQuery
	jwtGenerator JWTGenerator
)

func ModifiedLogin(c *gin.Context) {
	var req entities.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Hash the password before comparison
	hashedPassword := algorithm.SHA256Hash(req.Password)
	storedPassword, status := accountQuery.GetAccountPasswordByUsername(req.Username)
	if status == http.StatusInternalServerError {
		c.JSON(status, gin.H{"error": "Error retrieving stored password"})
		return
	}
	if storedPassword != hashedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid username or password",
			"error":   "Invalid credentials",
		})
		return
	}

	role := accountQuery.GetRoleByUsername(req.Username)
	tokenString, err := jwtGenerator.GenerateJWT(req.Username, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Error generating token",
			"error":   "Error generating token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"role":    role,
	})
}
