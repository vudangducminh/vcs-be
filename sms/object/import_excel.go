package object

type ImportExcelSuccessResponse struct {
	Message      string `json:"message" example:"Excel file imported successfully"`
	AddedServers int    `json:"added_servers" example:"5"`
	ErrorServers int    `json:"error_servers" example:"2"`
}

type ImportExcelInvalidFileFormatResponse struct {
	Error string `json:"error" example:"Invalid file format"`
}

type ImportExcelRetrieveFailedResponse struct {
	Error string `json:"error" example:"Failed to retrieve file"`
}

type ImportExcelOpenFileFailedResponse struct {
	Error string `json:"error" example:"Failed to open file"`
}

type ImportExcelReadFileFailedResponse struct {
	Error string `json:"error" example:"Failed to read Excel rows"`
}

type ImportExcelElasticsearchErrorResponse struct {
	Error string `json:"error" example:"Failed to add server to Elasticsearch from Excel row"`
}
