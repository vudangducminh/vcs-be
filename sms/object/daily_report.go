package object

type DailyReportRequest struct {
	Email     string `json:"email" binding:"required,email"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type DailyReportInvalidRequestResponse struct {
	Message string `json:"message" example:"Invalid request"`
	Error   string `json:"error" example:"Invalid request body"`
}

type DailyReportResponse struct {
	Message string `json:"message" example:"Request saved successfully"`
}

type DailyReportInternalServerErrorResponse struct {
	Message string `json:"message" example:"Internal server error"`
	Error   string `json:"error" example:"Failed to save request"`
}
