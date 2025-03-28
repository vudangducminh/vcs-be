package server

import (
	"basic-server/database"
	"fmt"
	"net/http"
)

func GetPlayerScoreOnServer(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Query().Get("player")
	score := database.GetPlayerScoreByName(player)
	if score == -2 {
		fmt.Fprint(w, "Error")
	} else if score == -1 {
		fmt.Fprint(w, "No records")
	} else {
		fmt.Fprint(w, score)
	}
}
func AddPlayerScoreOnServer(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Query().Get("player")
	if player == "" {
		http.Error(w, "Missing 'player' parameter", http.StatusBadRequest)
		return
	}
	scoreInString := r.URL.Query().Get("score")
	if scoreInString == "" {
		http.Error(w, "Missing 'score' parameter", http.StatusBadRequest)
		return
	}
	var score int = 0
	for i := 0; i < len(scoreInString); i++ {
		score = score*10 + int(scoreInString[i]) - 48
	}
	status := database.InsertPlayerScore(player, score)
	if !status {
		fmt.Fprint(w, "Can't insert new data")
		return
	}
	fmt.Fprintf(w, "Inserted player %q with score %d", player, score)
}
