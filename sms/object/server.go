package object

type Server struct {
	ServerId        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	Status          string `json:"status"`            // e.g., "active", "inactive", "maintenance"
	Uptime          int64  `json:"uptime"`            // e.g., "3666" for 1 hour 1 minute and 6 seconds
	CreatedTime     int64  `json:"created_time"`      // Unix timestamp
	LastUpdatedTime int64  `json:"last_updated_time"` // Unix timestamp
	IPv4            string `json:"ipv4"`
}
