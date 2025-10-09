package object

type DailyReportRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type DailyReportResponse struct {
	Message string `json:"message" example:"Request saved successfully"`
}

type DailyReportInvalidRequestResponse struct {
	Error string `json:"error" example:"Invalid request body"`
}

type DailyReportConflictResponse struct {
	Error string `json:"error" example:"Email already exists"`
}

type DailyReportInternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}
