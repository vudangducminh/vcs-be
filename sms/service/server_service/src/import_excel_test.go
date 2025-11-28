package src

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"server_service/entities"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type mockBulkServerAdder struct {
	bulkStatus int
	addStatus  int
}

func (m *mockBulkServerAdder) BulkServerInfo(servers []entities.Server) int {
	return m.bulkStatus
}
func (m *mockBulkServerAdder) AddServerInfo(server entities.Server) int {
	return m.addStatus
}

func setupImportExcelTest(bsa BulkServerAdder) *gin.Engine {
	SetBulkServerAdder(bsa)
	r := gin.Default()
	r.POST("/import-excel", ModifiedImportExcel)
	return r
}

func createTestExcelFile(t *testing.T, filename string) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	// Write header
	f.SetSheetRow(sheet, "A1", &[]interface{}{"ID", "ServerName", "Status", "IPv4"})
	// Write one valid row
	f.SetSheetRow(sheet, "A2", &[]interface{}{"1", "TestServer", "active", "127.0.0.1"})
	if err := f.SaveAs(filename); err != nil {
		t.Fatal(err)
	}
}

func TestImportExcel_InvalidFile(t *testing.T) {
	r := setupImportExcelTest(&mockBulkServerAdder{})
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	req, _ := http.NewRequest("POST", "/import-excel", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestImportExcel_Success(t *testing.T) {
	r := setupImportExcelTest(&mockBulkServerAdder{bulkStatus: http.StatusCreated, addStatus: http.StatusCreated})

	// Create a valid Excel file using excelize
	filename := "test.xlsx"
	createTestExcelFile(t, filename)
	defer os.Remove(filename)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req, _ := http.NewRequest("POST", "/import-excel", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestImportExcel_InternalServerError(t *testing.T) {
	r := setupImportExcelTest(&mockBulkServerAdder{bulkStatus: http.StatusInternalServerError, addStatus: http.StatusInternalServerError})

	filename := "test.xlsx"
	createTestExcelFile(t, filename)
	defer os.Remove(filename)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req, _ := http.NewRequest("POST", "/import-excel", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
