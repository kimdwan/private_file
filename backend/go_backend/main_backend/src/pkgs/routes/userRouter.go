package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/pkgs/controllers"
)

func UserRouter(router *gin.Engine) {

	// 유저 관련 서비스
	userrouter := router.Group("user")
	userrouter.POST("signup", controllers.UserSignUpController)
	userrouter.POST("login", controllers.UserLoginController)
}
