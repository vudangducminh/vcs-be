package report_service

import (
	"bytes"
	"encoding/base64"
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

	// Verify Excel file size
	log.Printf("Excel file size: %d bytes", excelBuffer.Len())
	if excelBuffer.Len() == 0 {
		log.Println("Excel file is empty")
		return http.StatusInternalServerError
	}

	// Create email message
	var emailBuffer bytes.Buffer
	writer := multipart.NewWriter(&emailBuffer)

	// Email headers
	emailBuffer.WriteString(fmt.Sprintf("From: %s\r\n", senderEmail))
	emailBuffer.WriteString(fmt.Sprintf("To: %s\r\n", recipientEmail))
	emailBuffer.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	emailBuffer.WriteString("MIME-Version: 1.0\r\n")
	emailBuffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary()))

	// Write email body part
	bodyPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	})
	if err != nil {
		log.Printf("Failed to create email body part: %v", err)
		return http.StatusInternalServerError
	}
	if _, err := bodyPart.Write([]byte(body)); err != nil {
		log.Printf("Failed to write email body: %v", err)
		return http.StatusInternalServerError
	}

	// Write Excel attachment with proper encoding
	filename := fmt.Sprintf("daily_report_%s.xlsx", time.Now().Format("2006-01-02"))
	attachmentPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"Content-Disposition":       []string{fmt.Sprintf("attachment; filename=\"%s\"", filename)},
		"Content-Transfer-Encoding": []string{"base64"}, // ← Add this
	})
	if err != nil {
		log.Printf("Failed to create attachment part: %v", err)
		return http.StatusInternalServerError
	}

	// Encode Excel file in base64
	encoder := base64.NewEncoder(base64.StdEncoding, attachmentPart)
	if _, err := encoder.Write(excelBuffer.Bytes()); err != nil {
		log.Printf("Failed to write attachment: %v", err)
		return http.StatusInternalServerError
	}
	encoder.Close()

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
