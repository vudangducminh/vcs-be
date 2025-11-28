package src

import (
	"net/http"
	"net/http/httptest"
	"server_service/entities"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockServerViewer struct {
	servers    []entities.Server
	httpStatus int
}

func (m *mockServerViewer) GetServerByNameSubstr(str string) ([]entities.Server, int) {
	return m.servers, m.httpStatus
}
func (m *mockServerViewer) GetServerByIPv4Substr(str string) ([]entities.Server, int) {
	return m.servers, m.httpStatus
}
func (m *mockServerViewer) GetServerByStatus(str string) ([]entities.Server, int) {
	return m.servers, m.httpStatus
}

func setupViewServerTest(sv ServerViewer) *gin.Engine {
	SetServerViewer(sv)
	r := gin.Default()
	r.GET("/view-servers", ModifiedViewServer)
	return r
}

func TestViewServer_NoServersFound(t *testing.T) {
	r := setupViewServerTest(&mockServerViewer{servers: nil, httpStatus: http.StatusNotFound})
	req, _ := http.NewRequest("GET", "/view-servers?filter=server_name&string=notfound", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestViewServer_FailedToRetrieve(t *testing.T) {
	r := setupViewServerTest(&mockServerViewer{servers: nil, httpStatus: http.StatusInternalServerError})
	req, _ := http.NewRequest("GET", "/view-servers?filter=server_name&string=fail", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestViewServer_Success(t *testing.T) {
	servers := []entities.Server{
		{Id: "1", ServerName: "A", Status: "active", CreatedTime: 1, LastUpdatedTime: 2, IPv4: "127.0.0.1"},
		{Id: "2", ServerName: "B", Status: "inactive", CreatedTime: 3, LastUpdatedTime: 4, IPv4: "127.0.0.2"},
	}
	r := setupViewServerTest(&mockServerViewer{servers: servers, httpStatus: http.StatusOK})
	req, _ := http.NewRequest("GET", "/view-servers?filter=server_name&order=asc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestViewServer_SuccessDescOrder(t *testing.T) {
	servers := []entities.Server{
		{Id: "1", ServerName: "A", Status: "active", CreatedTime: 1, LastUpdatedTime: 2, IPv4: "127.0.0.1"},
		{Id: "2", ServerName: "B", Status: "inactive", CreatedTime: 3, LastUpdatedTime: 4, IPv4: "127.0.0.2"},
	}
	r := setupViewServerTest(&mockServerViewer{servers: servers, httpStatus: http.StatusOK})
	req, _ := http.NewRequest("GET", "/view-servers?filter=server_name&order=desc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestViewServer_SuccessAscOrderIPv4(t *testing.T) {
	servers := []entities.Server{
		{Id: "1", ServerName: "A", Status: "active", CreatedTime: 1, LastUpdatedTime: 2, IPv4: "123.123.231.132"},
		{Id: "2", ServerName: "B", Status: "inactive", CreatedTime: 3, LastUpdatedTime: 4, IPv4: "123.123.231.133"},
	}
	r := setupViewServerTest(&mockServerViewer{servers: servers, httpStatus: http.StatusOK})
	req, _ := http.NewRequest("GET", "/view-servers?filter=ipv4&order=asc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestViewServer_InternalServerError(t *testing.T) {
	r := setupViewServerTest(&mockServerViewer{servers: nil, httpStatus: http.StatusInternalServerError})
	req, _ := http.NewRequest("GET", "/view-servers?filter=server_name&string=error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
