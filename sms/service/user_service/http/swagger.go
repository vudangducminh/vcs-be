package http

import (
	_ "user_service/docs"
	src "user_service/src"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {
	r := gin.Default()
	user := r.Group("api/v1/user")
	{
		user.POST("/login", src.HandleLogin)
		user.POST("/register", src.HandleRegister)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8800/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
