package object

type UpdateServerRequest struct {
	ServerId        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	Status          string `json:"status"`            // e.g., "active", "inactive", "maintenance"
	Uptime          int    `json:"uptime"`            // e.g., "3666" for 1 hour 1 minute and 6 seconds
	CreatedTime     string `json:"created_time"`      // ISO 8601 format
	LastUpdatedTime string `json:"last_updated_time"` // ISO 8601 format
	IPv4            string `json:"ipv4"`              // IPv4 address of the server
}
