package main

import (
	"log"
	initialize "report_service/init"
)

// @title           VCS System Management API
// @version         1.0
// @description     This is a sample server for VCS System Management API.
// @contact.name    Vu Dang Duc Minh
// @contact.email   vudangducminh@gmail.com
// @contact.url     https://github.com/vudangducminh
// @BasePath        /
// @schemes         http
// @host            report.localhost
// @Tag.name		Report
// @Tag.description "Operations related to generating and sending reports"
func main() {
	initialize.Init()
	log.Println("Report Service is running on port 8802")
}
