package src

import (
	"log"
	"net/http"
	"server_service/entities"
	elastic_query "server_service/infrastructure/elasticsearch/query"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags         Server
// @Summary      Export server data to Excel
// @Description  Export server data to an Excel file with optional filtering and ordering
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        order query string false "Order of results, either 'asc' or 'desc'. If not provided or using the wrong order format, the default order is ascending"
// @Param        filter query string false "Filter by server_name, ipv4, or status. If not provided or using the wrong filter format, then there is no filter applied"
// @Param        string query string false "Substring to search in server_name, ipv4, or status"
// @Success      200 {object} entities.ExportExcelSuccessResponse "Excel file exported successfully"
// @Failure      400 {object} entities.ExportExcelBadRequestResponse "Invalid request parameters"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} entities.ExportExcelInternalServerErrorResponse "Failed to retrieve server details"
// @Failure      500 {object} entities.ExportExcelFailedResponse "Failed to export into Excel file"
// @Router       /servers/export-excel [get]
func ExportExcel(c *gin.Context) {
	var servers []entities.Server
	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "asc" // Default order if not specified
	}
	filter := c.Query("filter")
	if filter != "server_name" && filter != "ipv4" && filter != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter"})
		return
	}
	str := c.Query("string")
	log.Printf("Received request to export server with filter '%s' and substring: '%s'", filter, str)
	if str == "undefined" || str == "{string}" {
		str = ""
	}
	log.Printf("Received request to export server with filter '%s' and substring: '%s'", filter, str)
	var httpStatus int = 200
	switch filter {
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

type ServerQuery interface {
	GetServerByNameSubstr(str string) ([]entities.Server, int)
	GetServerByIPv4Substr(str string) ([]entities.Server, int)
	GetServerByStatus(str string) ([]entities.Server, int)
}

var serverQuery ServerQuery

func SetServerQuery(sq ServerQuery) {
	serverQuery = sq
}

func ModifiedExportExcel(c *gin.Context) {
	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	filter := c.Query("filter")
	if filter != "server_name" && filter != "ipv4" && filter != "status" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter"})
		return
	}
	str := c.Query("string")
	if str == "undefined" || str == "{string}" {
		str = ""
	}

	var servers []entities.Server
	var httpStatus int
	switch filter {
	case "server_name":
		servers, httpStatus = serverQuery.GetServerByNameSubstr(str)
	case "ipv4":
		servers, httpStatus = serverQuery.GetServerByIPv4Substr(str)
	case "status":
		servers, httpStatus = serverQuery.GetServerByStatus(str)
	}

	if httpStatus == http.StatusNotFound {
		c.JSON(http.StatusOK, gin.H{"message": "No servers found with the given requirements"})
		return
	} else if httpStatus != http.StatusOK {
		c.JSON(httpStatus, gin.H{"error": "Failed to retrieve server details"})
		return
	}

	sort.Slice(servers, func(i, j int) bool {
		var less bool
		switch filter {
		case "status":
			less = servers[i].Status < servers[j].Status
		case "ipv4":
			less = servers[i].IPv4 < servers[j].IPv4
		default:
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
	headers := []string{"Index", "Server Name", "Status", "IPv4"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
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

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=servers.xlsx")
	c.Header("File-Name", "servers.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export into Excel file"})
		return
	}
}
