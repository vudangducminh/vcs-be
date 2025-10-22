package http

import (
	"time"
	_ "user_service/docs"
	src "user_service/src"
	"user_service/src/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func Init() {
	r := gin.Default()
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Second/3), 5)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	health := r.Group("api/v1/health")
	{
		health.GET("", src.HealthCheck)
	}

	user := r.Group("api/v1/user")
	{
		user.POST("/login", src.HandleLogin)
	}
	user = r.Group("api/v1/user", middleware.AuthAdmin())
	{
		user.POST("/register", src.HandleRegister)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8800/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
