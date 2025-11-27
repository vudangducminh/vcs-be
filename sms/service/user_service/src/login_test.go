package src

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_service/src/algorithm"

	"github.com/gin-gonic/gin"
)

func SetDependencies(aq AccountQuery, jg JWTGenerator) {
	accountQuery = aq
	jwtGenerator = jg
}

type mockAccountQuery struct {
	password string
	status   int
	role     string
}

func (m *mockAccountQuery) GetAccountPasswordByUsername(username string) (string, int) {
	return m.password, m.status
}
func (m *mockAccountQuery) GetRoleByUsername(username string) string {
	return m.role
}

type mockJWTGenerator struct {
	token string
	err   error
}

func (m *mockJWTGenerator) GenerateJWT(username, role string) (string, error) {
	return m.token, m.err
}

func setupTest(aq AccountQuery, jg JWTGenerator) *gin.Engine {
	SetDependencies(aq, jg)
	r := gin.Default()
	r.POST("/login", ModifiedLogin)
	return r
}
func TestLogin_InvalidBody(t *testing.T) {
	r := setupTest(&mockAccountQuery{}, &mockJWTGenerator{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestLogin_MissingUsernameOrPassword(t *testing.T) {
	r := setupTest(&mockAccountQuery{}, &mockJWTGenerator{})
	body := []byte(`{"username": "", "password": ""}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
func TestLogin_DBError(t *testing.T) {
	r := setupTest(&mockAccountQuery{status: http.StatusInternalServerError}, &mockJWTGenerator{})
	body := []byte(`{"username": "user", "password": "pass"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	r := setupTest(&mockAccountQuery{password: "other"}, &mockJWTGenerator{})
	body := []byte(`{"username": "user", "password": "pass"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestLogin_JWTError(t *testing.T) {
	r := setupTest(&mockAccountQuery{password: "hashed", role: "admin"}, &mockJWTGenerator{err: errors.New("fail")})
	body := []byte(`{"username": "user", "password": "pass"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	hashed := algorithm.SHA256Hash("pass")
	r := setupTest(&mockAccountQuery{password: hashed, role: "admin"}, &mockJWTGenerator{token: "token"})
	body := []byte(`{"username": "user", "password": "pass"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
