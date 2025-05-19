package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SystemRoutes(r *gin.Engine) {

	systemRoutes := r.Group("api/system")
	{
		systemRoutes.GET("/service_health_status", controllers.GetSystemHealthUpdate)
		systemRoutes.GET("/system_usage", controllers.GetSystemUsage)
		systemRoutes.GET("/ws", controllers.GetSystemUsageWS)
	}
}
