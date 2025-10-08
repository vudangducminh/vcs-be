package main

import (
	initialize "healthcheck_service/init"
	"log"
)

// @title           VCS System Management API
// @version         1.0
// @description     This is a sample server for VCS System Management API.
// @contact.name    Vu Dang Duc Minh
// @contact.email   vudangducminh@gmail.com
// @contact.url     https://github.com/vudangducminh
// @BasePath        /api/v1
// @schemes         http
// @host            localhost:8803
// @Tag.name		Healthcheck
// @Tag.description "Check the health status of services and managed servers"
func main() {
	initialize.Init()
	log.Println("Health Check Service is running on port 8803")
}
