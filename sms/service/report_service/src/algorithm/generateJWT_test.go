package algorithm

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	_, err := GenerateJWT("sus", "admin")
	if err != nil {
		t.Errorf("Error generating JWT: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	tokenStr, err := GenerateJWT("sus", "admin")
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	token, err := ValidateJWT(tokenStr)
	if err != nil {
		t.Errorf("Error validating JWT: %v", err)
		return
	}

	if !token.Valid {
		t.Error("JWT token is not valid")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Claims are not of type jwt.MapClaims")
		return
	}

	username, ok := claims["username"].(string)
	if !ok || username != "sus" {
		t.Errorf("Expected username 'sus', got '%s'", username)
	}

	if claims["password"] != "123" {
		t.Errorf("Expected password '123', got '%s'", claims["password"])
	}
}
