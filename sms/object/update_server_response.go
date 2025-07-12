package object

type UpdateServerResponse struct {
	ServerId        string `json:"server_id"`                        // Unique identifier for the server
	ServerName      string `json:"server_name"`                      // Name of the server
	ServerIPv4      string `json:"ipv4"`                             // IPv4 address of the server
	ServerStatus    string `json:"status"`                           // Status of the server, e.g., "active", "inactive", "maintenance"
	LastUpdatedTime string `json:"last_updated_time"`                // Last updated time in ISO 8601 format
	Message         string `json:"message" example:"Server updated"` // Confirmation message
}
