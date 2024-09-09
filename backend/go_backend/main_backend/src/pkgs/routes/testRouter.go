package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/pkgs/controllers"
)

func TestRouter(router *gin.Engine) {
	testrouter := router.Group("test")
	testrouter.GET("1", controllers.TestController1)
	testrouter.POST("2", controllers.TestController2)
}
