package report_service

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/xuri/excelize/v2"
)

func SendEmail(excelFile *excelize.File, recipientEmail string, subject string, body string) int {
	// SMTP configuration - use environment variables
	smtpHost := "smtp.gmail.com" // e.g., "smtp.gmail.com"
	smtpPort := "587"            // e.g., "587"
	senderEmail := "vudangducminh@gmail.com"
	senderPassword := "fzyzzdsglqrznvpw"

	if smtpHost == "" || smtpPort == "" || senderEmail == "" || senderPassword == "" {
		log.Println("SMTP configuration not set")
		return http.StatusInternalServerError
	}

	// Create the Excel file in memory
	var excelBuffer bytes.Buffer
	if err := excelFile.Write(&excelBuffer); err != nil {
		log.Printf("Failed to write Excel file: %v", err)
		return http.StatusInternalServerError
	}

	// Create email with attachment
	var emailBuffer bytes.Buffer
	writer := multipart.NewWriter(&emailBuffer)

	// Email headers
	headers := map[string]string{
		"From":    senderEmail,
		"To":      recipientEmail,
		"Subject": subject,
	}

	// Write headers
	for key, value := range headers {
		emailBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	emailBuffer.WriteString("MIME-Version: 1.0\r\n")
	emailBuffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary()))

	// Write email body
	bodyPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	})
	if err != nil {
		log.Printf("Failed to create email body part: %v", err)
		return http.StatusInternalServerError
	}
	bodyPart.Write([]byte(body))

	// Write Excel attachment
	attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"Content-Disposition": []string{fmt.Sprintf("attachment; filename=\"daily_report_%s.xlsx\"", time.Now().Format("2006-01-02"))},
	})
	if err != nil {
		log.Printf("Failed to create attachment part: %v", err)
		return http.StatusInternalServerError
	}
	attachmentPart.Write(excelBuffer.Bytes())

	writer.Close()

	// Send email
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		senderEmail,
		[]string{recipientEmail},
		emailBuffer.Bytes(),
	)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
