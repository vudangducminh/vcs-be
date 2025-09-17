package object

type ExportExcelRequest struct {
	Order  string `json:"order"`  // e.g., "asc", "desc"
	Filter string `json:"filter"` // e.g., "server_id", "server_name", "ipv4", "status"
	String string `json:"string"` // Substring to search in server_id, server_name, ipv4, or status
}
type ExportExcelSuccessResponse struct {
	Message string `json:"message" example:"Excel file exported successfully"`
}

type ExportExcelBadRequestResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

type ExportExcelStatusNotFoundResponse struct {
	Message string `json:"message" example:"No servers found with the given requirements"`
}

type ExportExcelInternalServerErrorResponse struct {
	Error string `json:"error" example:"Failed to retrieve server details"`
}

type ExportExcelFailedResponse struct {
	Error string `json:"error" example:"Failed to export into Excel file"`
}
