package object

type UpdateServerRequest struct {
	Id         string `json:"_id"`
	ServerName string `json:"server_name"`
	Status     string `json:"status"` // e.g., "active", "inactive", "maintenance"
}

type UpdateServerSuccessResponse struct {
	Id              string `json:"_id"`                              // Unique identifier for the server
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

type UpdateServerConflictResponse struct {
	Error string `json:"error" example:"Server already exists with the same IPv4 address"`
}
