package object

type RegisterRequest struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type RegisterSuccessResponse struct {
	Message string `json:"message" example:"Registration successful"`
}

type RegisterBadRequestResponse struct {
	Message string `json:"message" example:"Invalid request body"`
}

type RegisterConflictResponse struct {
	Message string `json:"message" example:"Account already exists"`
}

type RegisterInternalServerErrorResponse struct {
	Message string `json:"message" example:"Internal server error"`
}
