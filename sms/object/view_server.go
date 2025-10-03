package object

type ViewServerRequest struct {
	Order  string `json:"order"`  // e.g., "asc", "desc"
	Filter string `json:"filter"` // e.g., "_id", "server_name", "ipv4", "status"
	String string `json:"string"` // Substring to search in _id, server_name, ipv4, or status
}

type ViewServerSuccessResponse struct {
	Servers []Server `json:"servers"`
}

type ViewServerBadRequestResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

type ViewServerInternalServerErrorResponse struct {
	Error string `json:"error" example:"Failed to retrieve server details"`
}
