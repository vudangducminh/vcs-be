package server

import (
	"basic-server/database"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	if !database.IsConnected() {
		database.ConnectToDB()
	}
	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/q?player=Floyd", nil)
		fmt.Println(request)
		response := httptest.NewRecorder()
		fmt.Println(response)

		GetPlayerScoreOnServer(response, request)

		got := response.Body.String()
		want := "15"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
