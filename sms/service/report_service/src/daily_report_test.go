package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"report_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockEmailAdder struct {
	status int
}

func (m *mockEmailAdder) AddEmailInfo(email entities.Email) int {
	return m.status
}

func setupDailyReportTest(ea EmailAdder) *gin.Engine {
	SetEmailAdder(ea)
	r := gin.Default()
	r.POST("/daily-report", ModifiedDailyReport)
	return r
}

func TestDailyReport_InvalidBody(t *testing.T) {
	r := setupDailyReportTest(&mockEmailAdder{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/daily-report", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestDailyReport_MissingEmail(t *testing.T) {
	r := setupDailyReportTest(&mockEmailAdder{})
	body := []byte(`{"email": ""}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/daily-report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestDailyReport_Success(t *testing.T) {
	r := setupDailyReportTest(&mockEmailAdder{status: http.StatusCreated})
	body := []byte(`{"email": "test@example.com"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/daily-report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestDailyReport_Conflict(t *testing.T) {
	r := setupDailyReportTest(&mockEmailAdder{status: http.StatusConflict})
	body := []byte(`{"email": "test@example.com"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/daily-report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409, got %d", w.Code)
	}
}

func TestDailyReport_InternalServerError(t *testing.T) {
	r := setupDailyReportTest(&mockEmailAdder{status: http.StatusInternalServerError})
	body := []byte(`{"email": "test@example.com"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/daily-report", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
