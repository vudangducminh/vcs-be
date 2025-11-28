package src

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockDBChecker struct {
	connected bool
}

func (m *mockDBChecker) IsConnected() bool {
	return m.connected
}

func setupHealthTest(checker DBChecker) *gin.Engine {
	SetDBChecker(checker)
	r := gin.Default()
	r.GET("/health", ModifiedHealthCheck)
	return r
}

func TestHealthCheck_Healthy(t *testing.T) {
	r := setupHealthTest(&mockDBChecker{connected: true})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestHealthCheck_Unhealthy(t *testing.T) {
	r := setupHealthTest(&mockDBChecker{connected: false})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503, got %d", w.Code)
	}
}
