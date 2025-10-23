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
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Second*10), 50)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	health := r.Group("/health")
	{
		health.GET("", src.HealthCheck)
	}

	user := r.Group("/user")
	{
		user.POST("/login", src.HandleLogin)
	}
	user = r.Group("/user", middleware.AuthAdmin())
	{
		user.POST("/register", src.HandleRegister)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
