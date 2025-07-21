package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends an email using credentials and server info from environment variables.
func SendEmail(to []string, subject, body string) error {
	// --- IMPORTANT ---
	// Load credentials securely from environment variables.
	// NEVER hardcode them in your code.
	from := os.Getenv("SMTP_FROM_EMAIL")
	password := os.Getenv("SMTP_APP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
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
		"%s\r\n", to[0], subject, body))

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
