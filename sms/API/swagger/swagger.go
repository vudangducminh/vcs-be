package swagger

import (
	_ "sms/API/swagger/docs"
	"sms/client/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title VCS System Management System API
// @version 1.0
// @description This is a sample server for VCS System Management System API.
// @contact.name Vu Dang Duc Minh
// @contact.email vudangducminh@gmail.com
// @contact.url https://geithub.com/vudangducminh
// @BasePath /api/v1
// @schemes http
// @host localhost:8800
// @port 8800

func ConnectToSwagger() {
	r := gin.Default()
	users := r.Group("/api/v1/users")
	{
		users.GET("/login", handler.HandleLogin)
		users.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Register endpoint",
			})
		})
		users.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Dashboard endpoint",
			})
		})
		users.POST("/add-server", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Add server endpoint",
			})
		})
		users.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
