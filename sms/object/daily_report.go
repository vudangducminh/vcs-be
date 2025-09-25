package object

type DailyReportRequest struct {
	Email string `json:"email" binding:"required,email"`
}
