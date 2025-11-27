package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_service/entities"

	"github.com/gin-gonic/gin"
)

type mockUpdateQuery struct {
	account    entities.Account
	updateCode int
}

func (m *mockUpdateQuery) GetAccountByUsername(username string) entities.Account {
	return m.account
}
func (m *mockUpdateQuery) UpdateAccountInfo(account entities.Account) int {
	return m.updateCode
}

func setupUpdateTest(q UpdateQuery) *gin.Engine {
	SetUpdateQuery(q)
	r := gin.Default()
	r.PUT("/users/update-user", ModifiedUpdateUser)
	return r
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	r := setupUpdateTest(&mockUpdateQuery{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update-user", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestUpdateUser_UsernameRequired(t *testing.T) {
	r := setupUpdateTest(&mockUpdateQuery{})
	body := []byte(`{"username": ""}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update-user", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	r := setupUpdateTest(&mockUpdateQuery{account: entities.Account{}})
	body := []byte(`{"username": "notfound", "role": "user"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update-user", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	r := setupUpdateTest(&mockUpdateQuery{account: entities.Account{Username: "user"}, updateCode: http.StatusCreated})
	body := []byte(`{"username": "user", "role": "admin"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update-user", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestUpdateUser_Fail(t *testing.T) {
	r := setupUpdateTest(&mockUpdateQuery{account: entities.Account{Username: "user"}, updateCode: http.StatusInternalServerError})
	body := []byte(`{"username": "user", "role": "admin"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/update-user", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}
