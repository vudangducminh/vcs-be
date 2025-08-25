package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	elastic_query "sms/server/database/elasticsearch/query"
)

// SendEmail sends an email using credentials and server info from environment variables.
func SendEmail(to []string, averageUptime float32) error {
	from := os.Getenv("SMTP_FROM_EMAIL")
	password := os.Getenv("SMTP_APP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST") // smtp.gmail.com
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP environment variables not set")
	}

	smtpAddr := smtpHost + ":" + smtpPort

	// Template for the email subject and body
	// Template for the email subject and body
	var subject string = "Daily report from VCS System Management API"
	var totalServers int = elastic_query.GetTotalServersCount()
	var activeServers int = elastic_query.GetTotalActiveServersCount()
	var inactiveServers int = elastic_query.GetTotalInactiveServersCount()
	var maintenanceServers int = elastic_query.GetTotalMaintenanceServersCount()
	var otherServers int = totalServers - (activeServers + inactiveServers + maintenanceServers)
	var body string = "Number of servers in the system: " + fmt.Sprintf("%d", totalServers) +
		"\nNumber of active servers: " + fmt.Sprintf("%d", activeServers) +
		"\nNumber of inactive servers: " + fmt.Sprintf("%d", inactiveServers) +
		"\nNumber of servers in maintenance: " + fmt.Sprintf("%d", maintenanceServers) +
		"\nNumber of other servers: " + fmt.Sprintf("%d", otherServers) +
		"\nAverage server uptime: " + fmt.Sprintf("%.2f", averageUptime) + "%"

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
