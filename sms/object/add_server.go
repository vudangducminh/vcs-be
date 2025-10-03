package object

type AddServerRequest struct {
	ServerName string `json:"server_name"`
	Status     string `json:"status"` // e.g., "active", "inactive", "maintenance"
	IPv4       string `json:"ipv4"`   // IPv4 address of the server
}

type AddServerSuccessResponse struct {
	Id         string `json:"_id"`                            // Unique identifier for the server
	ServerName string `json:"server_name"`                    // Name of the server
	ServerIPv4 string `json:"ipv4"`                           // IPv4 address of the server
	Message    string `json:"message" example:"Server added"` // Confirmation message
}

type AddServerBadRequestResponse struct {
	Error string `json:"error" example:"Invalid request body"`
}

type AddServerConflictResponse struct {
	Error string `json:"error" example:"Server already exists"`
}

type AddServerInternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}
