package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"server_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockServerUpdater struct {
	getServerByIdResp entities.Server
	getServerByIdOk   bool
	updateStatus      int
}

func (m *mockServerUpdater) GetServerById(id string) (entities.Server, bool) {
	return m.getServerByIdResp, m.getServerByIdOk
}
func (m *mockServerUpdater) UpdateServerInfo(server entities.Server) int {
	return m.updateStatus
}

func setupUpdateServerTest(su ServerUpdater) *gin.Engine {
	SetServerUpdater(su)
	r := gin.Default()
	r.PUT("/update-server", ModifiedUpdateServer)
	return r
}

func TestUpdateServer_InvalidBody(t *testing.T) {
	r := setupUpdateServerTest(&mockServerUpdater{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update-server", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestUpdateServer_MissingId(t *testing.T) {
	r := setupUpdateServerTest(&mockServerUpdater{})
	body := []byte(`{"server_name": "TestServer", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestUpdateServer_NotFound(t *testing.T) {
	r := setupUpdateServerTest(&mockServerUpdater{
		getServerByIdResp: entities.Server{},
		getServerByIdOk:   false,
	})
	body := []byte(`{"_id": "notfound", "server_name": "TestServer", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestUpdateServer_Success(t *testing.T) {
	r := setupUpdateServerTest(&mockServerUpdater{
		getServerByIdResp: entities.Server{Id: "123", ServerName: "OldName", Status: "inactive"},
		getServerByIdOk:   true,
		updateStatus:      http.StatusOK,
	})
	body := []byte(`{"_id": "123", "server_name": "NewName", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestUpdateServer_Conflict(t *testing.T) {
	r := setupUpdateServerTest(&mockServerUpdater{
		getServerByIdResp: entities.Server{Id: "123", ServerName: "OldName", Status: "inactive"},
		getServerByIdOk:   true,
		updateStatus:      http.StatusConflict,
	})
	body := []byte(`{"_id": "123", "server_name": "NewName", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409, got %d", w.Code)
	}
}

func TestUpdateServer_InternalServerError(t *testing.T) {
	r := setupUpdateServerTest(&mockServerUpdater{
		getServerByIdResp: entities.Server{Id: "123", ServerName: "OldName", Status: "inactive"},
		getServerByIdOk:   true,
		updateStatus:      http.StatusInternalServerError,
	})
	body := []byte(`{"_id": "123", "server_name": "NewName", "status": "active"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/update-server", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
