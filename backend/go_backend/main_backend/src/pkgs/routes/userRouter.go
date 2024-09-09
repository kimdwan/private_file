package routes

import "github.com/gin-gonic/gin"

func UserRouter(router *gin.Engine) {
	userrouter := router.Group("user")
	userrouter.POST("login")
}
