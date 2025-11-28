package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"report_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockReportService struct {
	uptimeStatus int
	emailStatus  int
}

func (m *mockReportService) GetServerUptimeInRange(start, end int64, order, filter, str string) ([]entities.Server, int) {
	return []entities.Server{{Id: "1", ServerName: "TestServer"}}, m.uptimeStatus
}
func (m *mockReportService) GetTotalActiveServersCount(filter, str string) int      { return 1 }
func (m *mockReportService) GetTotalInactiveServersCount(filter, str string) int    { return 1 }
func (m *mockReportService) GetTotalMaintenanceServersCount(filter, str string) int { return 1 }
func (m *mockReportService) SendEmailReport(data interface{}, email, subject, body string) int {
	return m.emailStatus
}

func setupReportTest(rs ReportService) *gin.Engine {
	SetReportService(rs)
	r := gin.Default()
	r.POST("/report", ModifiedReportRequest)
	return r
}

func TestReport_InvalidBody(t *testing.T) {
	r := setupReportTest(&mockReportService{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/report", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestReport_FailUptime(t *testing.T) {
	r := setupReportTest(&mockReportService{uptimeStatus: http.StatusInternalServerError, emailStatus: http.StatusOK})
	body := []byte(`{"email": "test@example.com", "start_time": "2025-01-01T00:00:00Z", "end_time": "2025-01-02T00:00:00Z"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestReport_FailEmail(t *testing.T) {
	r := setupReportTest(&mockReportService{uptimeStatus: http.StatusOK, emailStatus: http.StatusInternalServerError})
	body := []byte(`{"email": "test@example.com", "start_time": "2025-01-01T00:00:00Z", "end_time": "2025-01-02T00:00:00Z"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestReport_Success(t *testing.T) {
	r := setupReportTest(&mockReportService{uptimeStatus: http.StatusOK, emailStatus: http.StatusOK})
	body := []byte(`{"email": "test@example.com", "start_time": "2025-01-01T00:00:00Z", "end_time": "2025-01-02T00:00:00Z"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
