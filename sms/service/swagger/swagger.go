package swagger

import (
	_ "sms/docs"
	auth_service "sms/service/auth_service"
	"sms/service/healthcheck_service"
	report_service "sms/service/report_service"
	server_service "sms/service/server_service"
	user_service "sms/service/user_service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConnectToSwagger() {
	r := gin.Default()

	health := r.Group("api/v1/health")
	{
		health.GET("/", healthcheck_service.ServiceHealthCheck)
	}
	user := r.Group("api/v1/user")
	{
		user.POST("/login", user_service.HandleLogin)
		user.POST("/register", user_service.HandleRegister)
	}
	server := r.Group("api/v1/server", auth_service.AuthUser())
	{
		server.GET("/view_servers/:order/:filter/:string", server_service.ViewServer)
		server.GET("/export_excel/:order/:filter/:string", server_service.ExportDataToExcel)
	}
	server = r.Group("api/v1/server", auth_service.AuthAdmin())
	{
		server.POST("/add_server", server_service.AddServer)
		server.PUT("/update_server", server_service.UpdateServer)
		server.DELETE("/delete_server", server_service.DeleteServer)
		server.POST("/import_excel", server_service.ImportExcel)
	}
	report := r.Group("api/v1/report", auth_service.AuthAdmin())
	{
		report.POST("/report/:order/:filter/:string", report_service.ReportRequest)
		report.POST("/daily_report", report_service.DailyReport)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8800/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
