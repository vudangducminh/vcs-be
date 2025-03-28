package database

import (
	"log"

	_ "github.com/lib/pq"
)

func GetPlayerScoreByName(name string) int {
	var personal_best []personal_best
	err := engine.Cols("id", "score").Alias("pb").Where("name = ?", name).Find(&personal_best)
	if err != nil {
		log.Fatal(err)
		return -2
	}
	if len(personal_best) == 0 {
		return -1
	}
	return personal_best[0].Score
}

func InsertPlayerScore(player string, score int) bool {
	has, err := engine.Where("name = ?", player).Count(new(personal_best))
	if err != nil {
		log.Fatal(err)
		return false
	}
	if has > 0 {
		return false
	}
	newScore := new(personal_best)
	newScore.Name = player
	newScore.Score = score
	rows, err := engine.Count(new(personal_best))
	if err != nil {
		log.Fatal(err)
		return false
	}
	newScore.Id = int(rows)
	affected, err := engine.Insert(newScore)

	if err != nil {
		log.Fatal(err)
		return false
	}
	return affected > 0
}
