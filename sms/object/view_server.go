package object

type ViewServerSuccessResponse struct {
	Servers []Server `json:"servers"`
}

type ViewServerBadRequestResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

type ViewServerInternalServerErrorResponse struct {
	Error string `json:"error" example:"Failed to retrieve server details"`
}
