package users_handler

import (
	"log"
	"net/http"
	"sms/object"
	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Users
// @Summary      Handle user registration
// @Description  Handle user registration by validating input and storing account information
// @Accept       json
// @Produce      json
// @Param        request body object.RegisterRequest true "Registration request"
// @Success      201 {object} object.RegisterSuccessResponse "Registration successful"
// @Failure      400 {object} object.RegisterBadRequestResponse "Invalid request body"
// @Failure      409 {object} object.RegisterConflictResponse "Account already exists"
// @Failure      500 {object} object.RegisterInternalServerErrorResponse "Internal server error
// @Router       /users/register [post]
func HandleRegister(c *gin.Context) {
	var req object.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Fullname == "" || req.Email == "" || req.Username == "" || req.Password == "" || req.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	log.Println("Registering Username:", req.Username)
	log.Println("Password:", req.Password)
	account := object.Account{
		Fullname: req.Fullname,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     "user",
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
