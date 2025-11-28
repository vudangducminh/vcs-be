package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"server_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockServerDeleter struct {
	getServerByIdResp entities.Server
	getServerByIdOk   bool
	deleteStatus      int
}

func (m *mockServerDeleter) GetServerById(id string) (entities.Server, bool) {
	return m.getServerByIdResp, m.getServerByIdOk
}
func (m *mockServerDeleter) DeleteServer(id string) int {
	return m.deleteStatus
}

func setupDeleteServerTest(sd ServerDeleter) *gin.Engine {
	SetServerDeleter(sd)
	r := gin.Default()
	r.DELETE("/delete-server", ModifiedDeleteServer)
	return r
}

func TestDeleteServer_InvalidBody(t *testing.T) {
	r := setupDeleteServerTest(&mockServerDeleter{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete-server", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestDeleteServer_NotFound(t *testing.T) {
	r := setupDeleteServerTest(&mockServerDeleter{
		getServerByIdResp: entities.Server{},
		getServerByIdOk:   false,
	})
	body := []byte(`{"_id": "notfound"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestDeleteServer_Success(t *testing.T) {
	r := setupDeleteServerTest(&mockServerDeleter{
		getServerByIdResp: entities.Server{Id: "123", ServerName: "TestServer", IPv4: "127.0.0.1"},
		getServerByIdOk:   true,
		deleteStatus:      http.StatusOK,
	})
	body := []byte(`{"_id": "123"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestDeleteServer_DeleteNotFound(t *testing.T) {
	r := setupDeleteServerTest(&mockServerDeleter{
		getServerByIdResp: entities.Server{Id: "123", ServerName: "TestServer", IPv4: "127.0.0.1"},
		getServerByIdOk:   true,
		deleteStatus:      http.StatusNotFound,
	})
	body := []byte(`{"_id": "123"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestDeleteServer_InternalServerError(t *testing.T) {
	r := setupDeleteServerTest(&mockServerDeleter{
		getServerByIdResp: entities.Server{Id: "123", ServerName: "TestServer", IPv4: "127.0.0.1"},
		getServerByIdOk:   true,
		deleteStatus:      http.StatusInternalServerError,
	})
	body := []byte(`{"_id": "123"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/delete-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
