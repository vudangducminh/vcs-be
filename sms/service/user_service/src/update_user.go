package src

import (
	"net/http"
	"user_service/entities"
	posgresql_query "user_service/infrastructure/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         User
// @Summary      Update user information
// @Description  Update user information
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body entities.UpdateUserRequest true "Update user request"
// @Success      201 {object} entities.UpdateUserResponse "Update user successful"
// @Failure      400 {object} entities.UpdateUserBadRequestResponse "Invalid request body"
// @Failure      404 {object} entities.UpdateUserNotFound "User not found"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      500 {object} entities.UpdateUserInternalServerErrorResponse "Internal server error
// @Router       /users/update-user [put]
func UpdateUser(c *gin.Context) {
	var req entities.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
		return
	}

	account := posgresql_query.GetAccountByUsername(req.Username)
	if account == (entities.Account{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	account.Role = req.Role

	httpStatus := posgresql_query.UpdateAccountInfo(account)
	if httpStatus == http.StatusCreated {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Update successful",
		})
	} else {
		c.JSON(httpStatus, gin.H{
			"message": "Update failed",
			"error":   "Internal server error",
		})
	}
}
