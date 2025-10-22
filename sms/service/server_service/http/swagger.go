package http

import (
	_ "server_service/docs" // This import is required for swagger docs to be registered
	"server_service/src"
	"server_service/src/middleware"
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

	server := r.Group("api/v1/server", middleware.AuthUser())
	{
		server.GET("/view_servers", src.ViewServer)
		server.GET("/export_excel", src.ExportDataToExcel)
	}
	server = r.Group("api/v1/server", middleware.AuthAdmin())
	{
		server.POST("/add_server", src.AddServer)
		server.PUT("/update_server", src.UpdateServer)
		server.DELETE("/delete_server", src.DeleteServer)
		server.POST("/import_excel", src.ImportExcel)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8801/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8801")
}
