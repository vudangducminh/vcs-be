package src

import (
	"log"
	"net/http"
	"user_service/entities"
	posgresql_query "user_service/infrastructure/postgresql/query"
	algorithm "user_service/src/algorithm"

	"github.com/gin-gonic/gin"
)

// @Tags         User
// @Summary      Handle user registration
// @Description  Handle user registration by validating input and storing account information
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body entities.RegisterRequest true "Registration request"
// @Success      201 {object} entities.RegisterSuccessResponse "Registration successful"
// @Failure      400 {object} entities.RegisterBadRequestResponse "Invalid request body"
// @Failure      409 {object} entities.RegisterConflictResponse "Account already exists"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      500 {object} entities.RegisterInternalServerErrorResponse "Internal server error
// @Router       /users/register [post]
func Register(c *gin.Context) {
	var req entities.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Role == "" || req.Fullname == "" || req.Email == "" || req.Username == "" || req.Password == "" || req.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	log.Println("Registering Username:", req.Username)
	log.Println("Password:", req.Password)
	account := entities.Account{
		Fullname: req.Fullname,
		Email:    req.Email,
		Username: req.Username,
		Password: algorithm.SHA256Hash(req.Password),
		Role:     req.Role,
	}
	httpStatus := posgresql_query.AddAccountInfo(account)
	if httpStatus == http.StatusCreated {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Registration successful",
		})
	} else {
		c.JSON(httpStatus, gin.H{
			"message": "Registration failed",
			"error":   "Account already exists or internal server error",
		})
	}
}

type AccountCreator interface {
	AddAccountInfo(account entities.Account) int
}

var accountCreator AccountCreator

func SetAccountCreator(ac AccountCreator) {
	accountCreator = ac
}

func ModifiedRegister(c *gin.Context) {
	var req entities.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	account := entities.Account{
		Fullname: req.Fullname,
		Email:    req.Email,
		Username: req.Username,
		Password: algorithm.SHA256Hash(req.Password),
		Role:     req.Role,
	}
	httpStatus := accountCreator.AddAccountInfo(account)
	if httpStatus == http.StatusCreated {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Registration successful",
		})
	} else {
		c.JSON(httpStatus, gin.H{
			"message": "Registration failed",
			"error":   "Account already exists or internal server error",
		})
	}
}
