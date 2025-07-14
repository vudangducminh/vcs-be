package object

type RegisterRequest struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"Registration successful"`
}
