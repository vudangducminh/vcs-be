package object

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginSuccessResponse struct {
	Message string `json:"message" example:"Login successful"`
}

type LoginBadRequestResponse struct {
	Message string `json:"message" example:"Invalid request body"`
	Error   string `json:"error" example:"Invalid request body"`
}

type LoginUnauthorizedResponse struct {
	Message string `json:"message" example:"Invalid credentials"`
	Error   string `json:"error" example:"Invalid credentials"`
}

type LoginInternalServerErrorResponse struct {
	Message string `json:"message" example:"Error generating token"`
	Error   string `json:"error" example:"Error generating token"`
}
