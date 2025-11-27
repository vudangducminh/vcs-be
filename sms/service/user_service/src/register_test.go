package src

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_service/entities"

	"github.com/gin-gonic/gin"
)

type mockAccountCreator struct {
	status int
}

func (m *mockAccountCreator) AddAccountInfo(account entities.Account) int {
	return m.status
}

func setupRegisterTest(ac AccountCreator) *gin.Engine {
	SetAccountCreator(ac)
	r := gin.Default()
	r.POST("/register", ModifiedRegister)
	return r
}

func TestRegister_InvalidBody(t *testing.T) {
	r := setupRegisterTest(&mockAccountCreator{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{invalid")))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	r := setupRegisterTest(&mockAccountCreator{})
	body := []byte(`{"username": "", "password": "", "confirm_password": "", "fullname": "", "email": "", "role": ""}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestRegister_PasswordMismatch(t *testing.T) {
	r := setupRegisterTest(&mockAccountCreator{})
	body := []byte(`{"username": "user", "password": "pass", "confirm_password": "wrong", "fullname": "name", "email": "user@gmail.com", "role": "user"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestRegister_Success(t *testing.T) {
	r := setupRegisterTest(&mockAccountCreator{status: http.StatusCreated})
	body := []byte(`{"username": "user", "password": "pass", "confirm_password": "pass", "fullname": "name", "email": "user@gmail.com", "role": "user"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestRegister_Failure(t *testing.T) {
	r := setupRegisterTest(&mockAccountCreator{status: http.StatusConflict})
	body := []byte(`{"username": "user", "password": "pass", "confirm_password": "pass", "fullname": "name", "email": "user@gmail.com", "role": "user"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409, got %d", w.Code)
	}
}
