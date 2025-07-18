package object

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

type ExportExcelExportFailedResponse struct {
	Error string `json:"error" example:"Failed to export Excel"`
}
