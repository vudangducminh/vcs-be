package src

import (
	"log"
	"net/http"
	"server_service/entities"
	elastic_query "server_service/infrastructure/elasticsearch/query"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags         Server
// @Summary      Import file from excel
// @Description  Import server data from an Excel file
// @Accept       multipart/form-data
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        file formData file true "Excel file to import"
// @Success      200 {object} entities.ImportExcelSuccessResponse "Excel file imported successfully"
// @Failure      400 {object} entities.ImportExcelInvalidFileFormatResponse "Invalid file format"
// @Failure      400 {object} entities.ImportExcelRetrieveFailedResponse "Failed to retrieve file"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} entities.ImportExcelOpenFileFailedResponse "Failed to open file"
// @Failure      500 {object} entities.ImportExcelReadFileFailedResponse "Failed to read Excel rows"
// @Failure      500 {object} entities.ImportExcelElasticsearchErrorResponse "Failed to add server to Elasticsearch from Excel row"
// @Router       /server/import_excel [post]
func ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		log.Println("Error opening file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer openedFile.Close()

	// Read the Excel file using excelize
	f, err := excelize.OpenReader(openedFile)
	if err != nil {
		log.Println("Error parsing Excel file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}

	// Get all rows from the first sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	log.Println("Sheet Name:", sheetName)
	log.Println("Total Rows:", len(rows))
	if err != nil {
		log.Println("Error reading rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Excel rows"})
		return
	}

	var isFirstRow bool = true
	var servers []entities.Server
	var errorServers []entities.Server
	var successServers []entities.Server
	for _, row := range rows {
		if isFirstRow {
			isFirstRow = false
			continue
		}
		var server entities.Server
		server.ServerName = row[1]
		server.Status = row[2]
		server.IPv4 = row[3]
		server.CreatedTime = time.Now().Unix()
		server.LastUpdatedTime = server.CreatedTime
		server.Uptime = []int{0}
		// if elastic_query.CheckServerExists(server.IPv4) {
		// 	log.Println("Server already exists in Elasticsearch, skipping row:", row)
		// 	errorServers = append(errorServers, server)
		// 	continue
		// }
		servers = append(servers, server)
		if len(servers) >= 250 {
			status := elastic_query.BulkServerInfo(servers)
			if status != http.StatusCreated {
				for _, s := range servers {
					status = elastic_query.AddServerInfo(s)
					if status != http.StatusCreated {
						errorServers = append(errorServers, s)
					} else {
						successServers = append(successServers, s)
					}
				}
			} else {
				successServers = append(successServers, servers...)
			}
			servers = nil
		}
	}
	status := elastic_query.BulkServerInfo(servers)
	if status != http.StatusCreated {
		for _, s := range servers {
			status = elastic_query.AddServerInfo(s)
			if status != http.StatusCreated {
				errorServers = append(errorServers, s)
			} else {
				successServers = append(successServers, s)
			}
		}
	} else {
		successServers = append(successServers, servers...)
	}
	servers = nil
	c.JSON(http.StatusOK, gin.H{
		"message":       "Excel file imported successfully",
		"added_servers": len(successServers),
		"error_servers": len(errorServers),
	})
}
