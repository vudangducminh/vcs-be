package object

type SendEmailRequest struct {
	JWT       string `json:"jwt" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type SendEmailInvalidRequestResponse struct {
	Message string `json:"message" example:"Invalid request"`
	Error   string `json:"error" example:"Invalid request body"`
}

type SendEmailResponse struct {
	Message string `json:"message" example:"Email sent successfully"`
}

type SendEmailInternalServerErrorResponse struct {
	Message string `json:"message" example:"Internal server error"`
	Error   string `json:"error" example:"Failed to send email"`
}
