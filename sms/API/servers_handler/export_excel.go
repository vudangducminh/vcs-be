package servers_handler

import (
	"net/http"
	"sms/object"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ExportServersToExcel exports server data to an Excel file and sends it as a download
func ExportServersToExcel(c *gin.Context, servers []object.Server) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Write header
	headers := []string{"Server Name", "Status", "IPv4"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data
	for rowIdx, server := range servers {
		values := []interface{}{
			server.ServerName,
			server.Status,
			server.IPv4,
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	// Set response headers for file download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=servers.xlsx")
	c.Header("File-Name", "servers.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")

	// Write file to response
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export Excel"})
	}
}
