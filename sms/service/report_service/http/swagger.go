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
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Second/3), 5)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	health := r.Group("api/v1/health")
	{
		health.GET("", src.HealthCheck)
	}

	report := r.Group("api/v1/report", middleware.AuthAdmin())
	{
		report.POST("/report", src.ReportRequest)
		report.POST("/daily_report", src.DailyReport)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8802/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8802")
}
