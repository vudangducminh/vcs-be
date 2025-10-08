package http

import (
	_ "healthcheck_service/docs"
	healthcheck_service "healthcheck_service/src"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {
	r := gin.Default()

	health := r.Group("api/v1/health")
	{
		health.GET("/", healthcheck_service.ServiceHealthCheck)
	}

	url := ginSwagger.URL("http://localhost:8803/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8803")
}
