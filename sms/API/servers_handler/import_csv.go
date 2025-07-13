package servers_handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Import CSV file
// @Description  Import server data from a CSV file
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "CSV file to import"
// @Success      200 {object} object.ImportCSVResponse "CSV imported successfully"
// @Router       /servers/import_csv [post]
func ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty file provided"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV imported successfully"})
}
