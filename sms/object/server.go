package object

type Server struct {
	Id              string `json:"_id"` // Unique identifier for the server
	ServerName      string `json:"server_name"`
	Status          string `json:"status"`            // e.g., "active", "inactive", "maintenance"
	Uptime          []int  `json:"uptime"`            // e.g., "3666" for 1 hour 1 minute and 6 seconds
	CreatedTime     int64  `json:"created_time"`      // Unix timestamp
	LastUpdatedTime int64  `json:"last_updated_time"` // Unix timestamp
	IPv4            string `json:"ipv4"`
}
