package object

type UpdateServerRequest struct {
	JWT             string `json:"jwt"` // JWT token for authentication
	ServerId        string `json:"server_id"`
	ServerName      string `json:"server_name"`
	Status          string `json:"status"`            // e.g., "active", "inactive", "maintenance"
	Uptime          int    `json:"uptime"`            // e.g., "3666" for 1 hour 1 minute and 6 seconds
	CreatedTime     string `json:"created_time"`      // ISO 8601 format
	LastUpdatedTime string `json:"last_updated_time"` // ISO 8601 format
	IPv4            string `json:"ipv4"`              // IPv4 address of the server
}

type UpdateServerSuccessResponse struct {
	ServerId        string `json:"server_id"`                        // Unique identifier for the server
	ServerName      string `json:"server_name"`                      // Name of the server
	ServerIPv4      string `json:"ipv4"`                             // IPv4 address of the server
	ServerStatus    string `json:"status"`                           // Status of the server, e.g., "active", "inactive", "maintenance"
	LastUpdatedTime string `json:"last_updated_time"`                // Last updated time in ISO 8601 format
	Message         string `json:"message" example:"Server updated"` // Confirmation message
}

type UpdateServerBadRequestResponse struct {
	Error string `json:"error" example:"Invalid request body"`
}

type UpdateServerStatusNotFoundResponse struct {
	Error string `json:"error" example:"Server not found"`
}

type UpdateServerInternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}
