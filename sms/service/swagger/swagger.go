package swagger

import (
	_ "sms/docs"
	auth_service "sms/service/auth_service"
	"sms/service/healthcheck_service"
	server_service "sms/service/server_service"
	user_service "sms/service/user_service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConnectToSwagger() {
	r := gin.Default()

	r.GET("/health", healthcheck_service.ServiceHealthCheck)

	users := r.Group("api/v1/users")
	{
		users.POST("/login", user_service.HandleLogin)
		users.POST("/register", user_service.HandleRegister)
	}
	servers := r.Group("api/v1/servers", auth_service.AuthUser())
	{
		servers.GET("/view_servers/:order/:filter/:string", server_service.ViewServer)
		servers.GET("/export_excel/:order/:filter/:string", server_service.ExportDataToExcel)
		servers.GET("/daily_report/:order/:filter/:string", server_service.DailyReportRequest)
	}
	servers = r.Group("api/v1/servers", auth_service.AuthAdmin())
	{
		servers.POST("/add_server", server_service.AddServer)
		servers.PUT("/update_server", server_service.UpdateServer)
		servers.DELETE("/delete_server", server_service.DeleteServer)
		servers.POST("/import_excel", server_service.ImportExcel)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8800/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
