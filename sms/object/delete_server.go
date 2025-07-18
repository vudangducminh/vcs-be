package object

type DeleteServerResponse struct {
	Message    string `json:"message" example:"Server deleted successfully"`
	ServerId   string `json:"server_id"`   // Unique identifier for the deleted server
	ServerName string `json:"server_name"` // Name of the deleted server
	ServerIPv4 string `json:"ipv4"`        // IPv4 address of the deleted server
}

type DeleteServerStatusNotFoundResponse struct {
	Error string `json:"error" example:"Server not found"`
}

type DeleteServerInternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}
