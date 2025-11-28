package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"server_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockServerAdder struct {
	status int
}

func (m *mockServerAdder) AddServerInfo(server entities.Server) int {
	return m.status
}

func setupAddServerTest(sa ServerAdder) *gin.Engine {
	SetServerAdder(sa)
	r := gin.Default()
	r.POST("/add-server", ModifiedAddServer)
	return r
}

func TestAddServer_InvalidBody(t *testing.T) {
	r := setupAddServerTest(&mockServerAdder{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add-server", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestAddServer_MissingFields(t *testing.T) {
	r := setupAddServerTest(&mockServerAdder{})
	body := []byte(`{"server_name": "", "ipv4": ""}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestAddServer_Success(t *testing.T) {
	r := setupAddServerTest(&mockServerAdder{status: http.StatusCreated})
	body := []byte(`{"server_name": "TestServer", "ipv4": "192.168.1.1", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestAddServer_Conflict(t *testing.T) {
	r := setupAddServerTest(&mockServerAdder{status: http.StatusConflict})
	body := []byte(`{"server_name": "TestServer", "ipv4": "192.168.1.1", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409, got %d", w.Code)
	}
}

func TestAddServer_InternalServerError(t *testing.T) {
	r := setupAddServerTest(&mockServerAdder{status: http.StatusInternalServerError})
	body := []byte(`{"server_name": "TestServer", "ipv4": "192.168.1.1", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
