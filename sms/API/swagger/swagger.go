package swagger

import (
	"sms/API/users_handler"
	_ "sms/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConnectToSwagger() {
	r := gin.Default()
	users := r.Group("/users")
	{
		users.POST("/login", users_handler.HandleLogin)
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
	}

	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8800/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
