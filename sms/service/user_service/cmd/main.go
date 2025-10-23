package main

import (
	"log"
	initialize "user_service/init"
)

// @title           VCS System Management API
// @version         1.0
// @description     This is a sample server for VCS System Management API.
// @contact.name    Vu Dang Duc Minh
// @contact.email   vudangducminh@gmail.com
// @contact.url     https://github.com/vudangducminh
// @BasePath        /
// @schemes         http
// @host            user.localhost
// @Tag.name		User
// @Tag.description "Operations related to User management"
func main() {
	initialize.Init()
	log.Println("User Service is running on port 8800")
}
