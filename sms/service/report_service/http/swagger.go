package http

import (
	_ "report_service/docs"
	"report_service/src"
	"report_service/src/middleware"
	"time"

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

	report := r.Group("/report", middleware.AuthUser())
	{
		report.POST("/report", src.ReportRequest)
		report.POST("/daily-report", src.DailyReport)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8802")
}
