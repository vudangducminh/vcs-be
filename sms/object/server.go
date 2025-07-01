package object

type Server struct {
	ServerId        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	Status          string `json:"status"`            // e.g., "active", "inactive", "maintenance"
	CreatedTime     string `json:"created_time"`      // ISO 8601 format
	LastUpdatedTime string `json:"last_updated_time"` // ISO 8601 format
	IPV4            string `json:"ipv4"`
}
