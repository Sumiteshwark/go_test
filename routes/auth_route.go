package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/signup", controllers.SignupHandler)
		authRoutes.POST("/login", controllers.LoginHandler)
		authRoutes.GET("/userinfo", controllers.GetValidateUserInfoHandler)
	}
}
