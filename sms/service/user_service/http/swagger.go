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

	users := r.Group("/users")
	{
		users.POST("/login", src.Login)
	}
	users = r.Group("/users", middleware.AuthAdmin())
	{
		users.POST("/register", src.Register)
		users.PUT("/update-user", src.UpdateUser)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
