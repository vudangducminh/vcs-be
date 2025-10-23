package http

import (
	_ "healthcheck_service/docs"
	"healthcheck_service/src"
	"healthcheck_service/src/middleware"
	"time"

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

	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8803")
}
