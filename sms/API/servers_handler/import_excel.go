package servers_handler

import (
	"log"
	"net/http"
	"sms/algorithm"
	"sms/object"
	redis_query "sms/server/database/cache/redis/query"
	elastic_query "sms/server/database/elasticsearch/query"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags         Servers
// @Summary      Import file from excel
// @Description  Import server data from an Excel file
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Excel file to import"
// @Param        jwt header string true "JWT token for authentication"
// @Success      200 {object} object.ImportExcelSuccessResponse "Excel file imported successfully"
// @Failure      400 {object} object.ImportExcelInvalidFileFormatResponse "Invalid file format"
// @Failure      400 {object} object.ImportExcelRetrieveFailedResponse "Failed to retrieve file"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} object.ImportExcelOpenFileFailedResponse "Failed to open file"
// @Failure      500 {object} object.ImportExcelReadFileFailedResponse "Failed to read Excel rows"
// @Failure      500 {object} object.ImportExcelElasticsearchErrorResponse "Failed to add server to Elasticsearch from Excel row"
// @Router       /servers/import_excel [post]
func ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	jwtToken := c.GetHeader("jwt")
	if err != nil {
		log.Println("Error retrieving file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	username := redis_query.GetUsernameByJWTToken(jwtToken)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
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
	if err != nil {
		log.Println("Error reading rows:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Excel rows"})
		return
	}

	var isFirstRow bool = true
	for _, row := range rows {
		if isFirstRow {
			isFirstRow = false
			continue
		}
		var server object.Server
		server.ServerId = algorithm.SHA256Hash(time.Now().String())
		server.ServerName = row[1]
		server.IPv4 = row[2]
		server.Status = row[3]
		server.CreatedTime = time.Now().Format(time.RFC3339)
		server.LastUpdatedTime = server.CreatedTime
		server.Uptime = 0
		status := elastic_query.AddServerInfo(server)
		if status != http.StatusCreated {
			log.Println("Failed to add server to Elasticsearch from Excel row:", row, "Status code:", status)
			c.JSON(status, gin.H{"error": "Failed to add server to Elasticsearch from Excel row"})
			continue
		}
		log.Println("Server added successfully from Excel row:", row)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Excel file imported successfully"})
}
