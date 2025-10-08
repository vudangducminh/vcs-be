package main

import (
	"log"
	initialize "server_service/init"
)

// @title           VCS System Management API
// @version         1.0
// @description     This is a sample server for VCS System Management API.
// @contact.name    Vu Dang Duc Minh
// @contact.email   vudangducminh@gmail.com
// @contact.url     https://github.com/vudangducminh
// @BasePath        /api/v1
// @schemes         http
// @host            localhost:8801
// @Tag.name		Server
// @Tag.description "Operations related to server management"
func main() {
	initialize.Init()
	log.Println("Server Service is running on port 8801")
}
