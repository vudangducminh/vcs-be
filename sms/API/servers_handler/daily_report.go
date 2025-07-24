package servers_handler

import (
	"net/http"
	"sms/object"
	redis_query "sms/server/database/cache/redis/query"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Create a request to send daily report email
// @Description  Create a request to send daily report email from YYYY-MM-DD to YYYY-MM-DD
// @Description  An email will be sent to the specified recipients at 00:00:00 UTC.
// @Description  Example date format: 2025-07-23T12:00:00Z
// @Accept       json
// @Produce      json
// @Param        request body object.DailyReportRequest true "Send email request"
// @Success      200 {object} object.DailyReportResponse "Email sent successfully"
// @Failure      400 {object} object.DailyReportInvalidRequestResponse "Invalid request"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} object.DailyReportInternalServerErrorResponse "Internal server error"
// @Router       /servers/daily_report [post]
func DailyReportEmailRequest(c *gin.Context) {
	var req object.DailyReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username := redis_query.GetUsernameByJWTToken(req.JWT)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
		return
	}

	beginTime := "1970-01-01T00:00:00Z"
	parsedBeginTime, err := time.Parse(time.RFC3339, beginTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse begin time"})
		return
	}
	parsedStartTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse time"})
		return
	}
	parsedEndTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse time"})
		return
	}
	if parsedStartTime.After(parsedEndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
		return
	}

	// Only need to save email & duration in redis
	startTimeInHHMMSS := parsedStartTime.Sub(parsedBeginTime)
	durationInHHMMSS := parsedEndTime.Sub(parsedStartTime)
	var startTimeInSecond = int64(startTimeInHHMMSS.Seconds())
	duration := int(durationInHHMMSS.Seconds())
	// log.Println(startTimeInSecond, " ", time.Now().Unix())
	if startTimeInSecond < time.Now().Unix() {
		duration -= int(time.Now().Unix() - startTimeInSecond)
		startTimeInSecond = time.Now().Unix()
	}
	if duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be greater than 0"})
		return
	}
	ok := redis_query.SaveDailyReportEmailRequest(req.Email, duration, startTimeInSecond)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save daily report email request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Request saved successfully"})
}
