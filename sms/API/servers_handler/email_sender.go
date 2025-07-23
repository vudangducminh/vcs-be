package servers_handler

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sms/object"
	redis_query "sms/server/database/cache/redis/query"
	"time"

	"github.com/gin-gonic/gin"
)

func Send(to []string, subject, body string) error {
	from := os.Getenv("SMTP_FROM_EMAIL")
	password := os.Getenv("SMTP_APP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST") // smtp.gmail.com
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP environment variables not set")
	}

	smtpAddr := smtpHost + ":" + smtpPort

	// Create the email message.
	// Note: Headers and body are separated by a double newline.
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	// Create an authentication object.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email.
	err := smtp.SendMail(smtpAddr, auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to %v", to)
	return nil
}

// @Tags         Servers
// @Summary      Send an email
// @Description  Send an email with a specified subject and body
// @Description  Example date format: 2025-07-23T12:00:00Z
// @Accept       json
// @Produce      json
// @Param        request body object.SendEmailRequest true "Send email request"
// @Success      200 {object} object.SendEmailResponse "Email sent successfully"
// @Failure      400 {object} object.SendEmailInvalidRequestResponse "Invalid request"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} object.SendEmailInternalServerErrorResponse "Internal server error"
// @Router       /servers/send_email [post]
func SendEmail(c *gin.Context) {
	var req object.SendEmailRequest
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
	var startTimeInSecond = int64(startTimeInHHMMSS.Seconds() + startTimeInHHMMSS.Minutes()*60 + startTimeInHHMMSS.Hours()*3600)
	duration := int64(durationInHHMMSS.Seconds() + durationInHHMMSS.Minutes()*60 + durationInHHMMSS.Hours()*3600)
	if startTimeInSecond < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be in the future"})
		return
	}
	log.Printf("Start time in seconds: %d, Duration: %d", startTimeInSecond, duration)
	log.Printf("Current time in seconds: %d", time.Now().Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
