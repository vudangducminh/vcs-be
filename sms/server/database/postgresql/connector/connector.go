package connector

import (
	"log"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var Engine *xorm.Engine
var isConnected = false

func IsConnected() bool {
	return isConnected
}

func ConnectToDB() {
	conns := "postgres://vudangducminh:Amogus69420@localhost:5432/postgres?sslmode=disable;"
	var err error
	Engine, err = xorm.NewEngine("postgres", conns)
	if err != nil {
		log.Fatal(err)
		return
	}
	isConnected = true
}
