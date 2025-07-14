package servers_handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// @Tags         Servers
// @Summary      Import file from excel
// @Description  Import server data from an Excel file
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Excel file to import"
// @Success      200 {object} object.ImportCSVResponse "Excel file imported successfully"
// @Router       /servers/import_excel [post]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
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

	// Example: Print each row
	for _, row := range rows {
		log.Println(row)
		// You can process each row and insert into your database here
	}

	c.JSON(http.StatusOK, gin.H{"message": "Excel file imported successfully"})
}
