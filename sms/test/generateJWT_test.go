package test

import (
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	_, err := GenerateJWT("sus", "123")
	if err != nil {
		t.Errorf("Error generating JWT: %v", err)
	}
}
