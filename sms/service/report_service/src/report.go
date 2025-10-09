package src

import (
	"fmt"
	"log"
	"net/http"
	"report_service/entities"
	elastic_query "report_service/infrastructure/elasticsearch/query"
	"report_service/src/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags         Report
// @Summary      Create a request to send report email
// @Description  Create a request to send report email from YYYY-MM-DD to YYYY-MM-DD
// @Description  Example date format: 2025-07-23T12:00:00Z
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        order query string false "Order of results, either 'asc' or 'desc'. If not provided or using the wrong order format, the default order is ascending"
// @Param        filter query string false "Filter by server_name, ipv4, or status. If not provided or using the wrong filter format, then there is no filter applied"
// @Param        string query string false "Substring to search in server_name, ipv4, or status"
// @Param        request body entities.ReportRequest true "Report request"
// @Success      200 {object} entities.ReportResponse "Email sent successfully"
// @Failure      400 {object} entities.ReportInvalidRequestResponse "Invalid request"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} entities.ReportInternalServerErrorResponse "Internal server error"
// @Failure      500 {object} entities.ExportExcelFailedResponse "Failed to export into Excel file"
// @Router       /report/report [post]
func ReportRequest(c *gin.Context) {
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
	var req entities.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	beginTime := "1970-01-01T00:00:00Z"
	parsedBeginTime, err := time.Parse(time.RFC3339, beginTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse begin time"})
		return
	}
	parsedStartTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse time"})
		return
	}
	parsedEndTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse time"})
		return
	}
	if parsedStartTime.After(parsedEndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start time must be before end time"})
		return
	}

	// Only need to save email & duration in redis
	startTimeInHHMMSS := parsedStartTime.Sub(parsedBeginTime)
	endTimeInHHMMSS := parsedEndTime.Sub(parsedBeginTime)
	var currentTimeInSecond = time.Now().Unix()
	var startTimeInSecond = min(int64(startTimeInHHMMSS.Seconds()), currentTimeInSecond)
	var endTimeInSecond = min(int64(endTimeInHHMMSS.Seconds()), currentTimeInSecond)
	var roundedStartTime = startTimeInSecond - (startTimeInSecond % 1200) // Round down to nearest 20 minutes
	var roundedEndTime = endTimeInSecond - (endTimeInSecond % 1200)       // Round down to nearest 20 minutes
	var roundedCurrentTime = currentTimeInSecond - (currentTimeInSecond % 1200)
	if startTimeInSecond%1200 == 0 {
		roundedStartTime -= 1200
	}
	if endTimeInSecond%1200 == 0 {
		roundedEndTime -= 1200
	}
	if currentTimeInSecond%1200 == 0 {
		roundedCurrentTime -= 1200
	}
	var beginBlock = int((roundedCurrentTime-roundedStartTime)/1200 + 1)
	var endBlock = int((roundedCurrentTime-roundedEndTime)/1200 + 1)
	serverDataList, status, averageUptimePercentage := elastic_query.GetServerUptimeInRange(beginBlock, endBlock, order, filter, str)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "Failed to retrieve server details"})
		return
	}

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Write header
	headers := []string{"Index", "Server ID", "Server Name", "Status", "IPv4", "Uptime", "Created Time", "Last Updated Time"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data

	for rowIdx, server := range serverDataList {
		// Convert timestamps to readable format
		createdTimeStr := time.Unix(server.CreatedTime, 0).Format("2006-01-02 15:04:05")
		lastUpdatedTimeStr := time.Unix(server.LastUpdatedTime, 0).Format("2006-01-02 15:04:05")
		serverUptime := time.Unix(int64(server.Uptime[0]), 0).Format("15:04:05")
		values := []interface{}{
			rowIdx + 1,
			server.Id,
			server.ServerName,
			server.Status,
			server.IPv4,
			serverUptime,
			createdTimeStr,
			lastUpdatedTimeStr,
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

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export into Excel file"})
		return
	}
	totalActiveServer := elastic_query.GetTotalActiveServersCount(filter, str)
	totalInactiveServer := elastic_query.GetTotalInactiveServersCount(filter, str)
	totalMaintenanceServer := elastic_query.GetTotalMaintenanceServersCount(filter, str)
	totalServer := totalActiveServer + totalInactiveServer + totalMaintenanceServer
	emailBody := "Here is your requested server report." + "\n"
	emailBody += "Total servers in the system: " + fmt.Sprintf("%d", totalServer) + "\n"
	emailBody += "Number of active servers: " + fmt.Sprintf("%d", totalActiveServer) + "\n"
	emailBody += "Number of inactive servers: " + fmt.Sprintf("%d", totalInactiveServer) + "\n"
	emailBody += "Number of maintenance servers: " + fmt.Sprintf("%d", totalMaintenanceServer) + "\n"
	emailBody += "Average uptime percentage across all servers: " + fmt.Sprintf("%.2f", averageUptimePercentage) + "%" + "\n"
	// Send email with the Excel file as attachment
	status = template.SendEmail(f, req.Email, "Server Report", emailBody)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": "Failed to send email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Servers exported successfully"})
	// Query from uptime[len - beginBlock] to uptime[len - endBlock]
	// Needs to provide total uptime of every single servers during this period
	// log.Println(startTimeInSecond, " ", time.Now().Unix())

}
