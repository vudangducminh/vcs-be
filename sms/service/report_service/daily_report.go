package report_service

import (
	"net/http"
	"sms/object"
	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Report
// @Summary      Create a request to send report email
// @Description  An email will be sent to the specified recipients at 00:00:00 UTC.
// @Description  Example date format: 2025-07-23T12:00:00Z
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body object.DailyReportRequest true "Daily report request"
// @Success      201 {object} object.DailyReportResponse "Request saved successfully"
// @Failure      400 {object} object.DailyReportInvalidRequestResponse "Invalid request"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      409 {object} object.DailyReportConflictResponse "Email already exists"
// @Failure      500 {object} object.DailyReportInternalServerErrorResponse "Internal server error"
// @Router       /report/daily_report [post]
func DailyReport(c *gin.Context) {
	// Implementation for daily report
	var req object.DailyReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	email := object.Email(req)
	status := posgresql_query.AddEmailInfo(email)
	if status == http.StatusCreated {
		c.JSON(http.StatusOK, gin.H{
			"message": "Request saved successfully",
		})
	} else {
		if status == http.StatusConflict {
			c.JSON(status, gin.H{
				"error": "Email already exists",
			})
			return
		}
		c.JSON(status, gin.H{
			"message": "Failed to save request",
			"error":   "Internal server error",
		})
	}
}
