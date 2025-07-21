package swagger

import (
	"sms/API/servers_handler"
	"sms/API/users_handler"
	"sms/API/website_handler"
	_ "sms/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ConnectToSwagger() {
	r := gin.Default()
	users := r.Group("api/v1/users")
	{
		users.POST("/login", users_handler.HandleLogin)
		users.POST("/register", users_handler.HandleRegister)
	}
	auth := r.Group("api/v1/auth")
	{
		auth.POST("/", website_handler.Authentication)
	}
	servers := r.Group("api/v1/servers")
	{
		servers.POST("/add_server", servers_handler.AddServer)
		servers.GET("/view_servers/:order/:filter/:string", servers_handler.ViewServer)
		servers.PUT("/update_server", servers_handler.UpdateServer)
		servers.DELETE("/delete_server", servers_handler.DeleteServer)
		servers.POST("/import_excel", servers_handler.ImportExcel)
		servers.GET("/export_excel/:order/:filter/:string", servers_handler.ExportDataToExcel)
	}
	// The host should match the @host annotation in main.go
	url := ginSwagger.URL("http://localhost:8800/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Run(":8800")
}
