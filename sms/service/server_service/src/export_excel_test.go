package src

import (
	"net/http"
	"net/http/httptest"
	"server_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockServerQuery struct {
	servers    []entities.Server
	httpStatus int
}

func (m *mockServerQuery) GetServerByNameSubstr(str string) ([]entities.Server, int) {
	return m.servers, m.httpStatus
}
func (m *mockServerQuery) GetServerByIPv4Substr(str string) ([]entities.Server, int) {
	return m.servers, m.httpStatus
}
func (m *mockServerQuery) GetServerByStatus(str string) ([]entities.Server, int) {
	return m.servers, m.httpStatus
}

func setupExportExcelTest(sq ServerQuery) *gin.Engine {
	SetServerQuery(sq)
	r := gin.Default()
	r.GET("/export-excel", ModifiedExportExcel)
	return r
}

func TestExportExcel_InvalidFilter(t *testing.T) {
	r := setupExportExcelTest(&mockServerQuery{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/export-excel?filter=wrong&order=asc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestExportExcel_NotFound(t *testing.T) {
	r := setupExportExcelTest(&mockServerQuery{servers: nil, httpStatus: http.StatusNotFound})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/export-excel?filter=server_name&order=asc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestExportExcel_FailedQuery(t *testing.T) {
	r := setupExportExcelTest(&mockServerQuery{servers: nil, httpStatus: http.StatusInternalServerError})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/export-excel?filter=server_name&order=asc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestExportExcel_Success(t *testing.T) {
	servers := []entities.Server{
		{ServerName: "A", Status: "active", IPv4: "1.1.1.1"},
		{ServerName: "B", Status: "inactive", IPv4: "2.2.2.2"},
	}
	r := setupExportExcelTest(&mockServerQuery{servers: servers, httpStatus: http.StatusOK})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/export-excel?filter=server_name&order=asc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusPartialContent {
		t.Errorf("Expected 200 or 206, got %d", w.Code)
	}
	// Optionally, check response headers for Excel file
	if w.Header().Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Errorf("Expected Excel content type, got %s", w.Header().Get("Content-Type"))
	}
}

func TestExportExcel_InternalServerError(t *testing.T) {
	r := setupExportExcelTest(&mockServerQuery{servers: []entities.Server{}, httpStatus: http.StatusInternalServerError})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/export-excel?filter=server_name&order=asc", nil)
	// Simulate an internal error by setting servers to nil after query
	serverQuery = &mockServerQuery{servers: nil, httpStatus: http.StatusInternalServerError}
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
