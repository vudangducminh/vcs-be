package object

type AddServerResponse struct {
	ServerId   string `json:"server_id"`                      // Unique identifier for the server
	ServerName string `json:"server_name"`                    // Name of the server
	ServerIPv4 string `json:"ipv4"`                           // IPv4 address of the server
	Message    string `json:"message" example:"Server added"` // Confirmation message
}
