package http

import (
	_ "report_service/docs"
	report_service "report_service/src"
	"report_service/src/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {
	r := gin.Default()

	report := r.Group("api/v1/report", middleware.AuthAdmin())
	{
		report.POST("/report", report_service.ReportRequest)
		report.POST("/daily_report", report_service.DailyReport)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8802/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8802")
}
