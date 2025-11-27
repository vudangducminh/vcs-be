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
	rateLimiter := middleware.NewIPRateLimiter(rate.Every(time.Second*10), 50)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	health := r.Group("/health")
	{
		health.GET("", src.HealthCheck)
	}

	server := r.Group("/servers", middleware.AuthViewServer())
	{
		server.GET("/view-servers", src.ViewServer)
	}
	server = r.Group("/servers", middleware.AuthExportExcel())
	{
		server.GET("/export-excel", src.ExportExcel)
	}
	server = r.Group("/servers", middleware.AuthAddServer())
	{
		server.POST("/add-server", src.AddServer)
	}
	server = r.Group("/servers", middleware.AuthUpdateServer())
	{
		server.PUT("/update-server", src.UpdateServer)
	}
	server = r.Group("/servers", middleware.AuthDeleteServer())
	{
		server.DELETE("/delete-server", src.DeleteServer)
	}
	server = r.Group("/servers", middleware.AuthImportExcel())
	{
		server.POST("/import-excel", src.ImportExcel)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8801")
}
