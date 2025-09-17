package servers_handler

import (
	"log"
	"net/http"
	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags         Servers
// @Summary      Export server data to Excel
// @Description  Export server data to an Excel file with optional filtering and ordering
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        order query string false "Order of results, either 'asc' or 'desc'. If not provided or using the wrong order format, the default order is ascending"
// @Param        filter query string false "Filter by server_id, server_name, ipv4, or status. If not provided or using the wrong filter format, then there is no filter applied"
// @Param        string path string false "Substring to search in server_id, server_name, ipv4, or status"
// @Success      200 {object} object.ExportExcelSuccessResponse "Excel file exported successfully"
// @Failure      400 {object} object.ExportExcelBadRequestResponse "Invalid request parameters"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} object.ExportExcelInternalServerErrorResponse "Failed to retrieve server details"
// @Failure      500 {object} object.ExportExcelFailedResponse "Failed to export into Excel file"
// @Router       /servers/export_excel/{order}/{filter}/{string} [get]
func ExportDataToExcel(c *gin.Context) {
	var servers []object.Server
	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "asc" // Default order if not specified
	}
	filter := c.Query("filter")
	if filter != "server_id" && filter != "server_name" && filter != "ipv4" && filter != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter"})
		return
	}
	str := c.Param("string")
	log.Printf("Received request to export server with filter '%s' and substring: '%s'", filter, str)
	if str == "undefined" || str == "{string}" {
		str = ""
	}
	log.Printf("Received request to export server with filter '%s' and substring: '%s'", filter, str)
	var httpStatus int = 200
	switch filter {
	case "server_id":
		servers, httpStatus = elastic_query.GetServerByIdSubstr(str)
	case "server_name":
		servers, httpStatus = elastic_query.GetServerByNameSubstr(str)
	case "ipv4":
		servers, httpStatus = elastic_query.GetServerByIPv4Substr(str)
	case "status":
		servers, httpStatus = elastic_query.GetServerByStatus(str)
	}

	if httpStatus == http.StatusNotFound {
		c.JSON(http.StatusOK, gin.H{"message": "No servers found with the given requirements"})
		return
	} else if httpStatus != http.StatusOK {
		c.JSON(httpStatus, gin.H{"error": "Failed to retrieve server details"})
		return
	}

	// Sort the servers based on the filter and order
	sort.Slice(servers, func(i, j int) bool {
		var less bool
		switch filter {
		case "server_id":
			less = servers[i].ServerId < servers[j].ServerId
		case "status":
			less = servers[i].Status < servers[j].Status
		case "ipv4":
			less = servers[i].IPv4 < servers[j].IPv4
		default: // Default to sorting by server_name
			less = servers[i].ServerName < servers[j].ServerName
		}
		if order == "desc" {
			return !less
		}
		return less
	})

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Write header
	headers := []string{"Index", "Server Name", "Status", "IPv4"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data
	for rowIdx, server := range servers {
		values := []interface{}{
			rowIdx + 1,
			server.ServerName,
			server.Status,
			server.IPv4,
		}
		for colIdx, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	// Set response headers for file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=servers.xlsx")
	c.Header("File-Name", "servers.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	// Write file to response
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export into Excel file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Excel file exported successfully"})
}
