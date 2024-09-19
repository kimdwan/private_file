package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/middlewares"
	"github.com/kimdwan/private_file/src/pkgs/controllers"
)

func AuthRouter(router *gin.Engine) {

	authrouter := router.Group("auth")
	authrouter.Use(middlewares.CheckJwtMiddleware())

	// 기본 서비스
	authrouter.GET("getprofile", controllers.AuthGetProfileImgController)
	authrouter.GET("getnickname", controllers.AuthGetNickNameController)
	authrouter.POST("upload/profile", controllers.AuthUploadProfileController)
	authrouter.GET("logout", controllers.AuthLogoutController)

	// 파일 관련 서비스
	authfilerouter := authrouter.Group("file")
	authfilerouter.POST("getfiles", controllers.AuthGetFileListController)
	authfilerouter.POST("searchfiles", controllers.AuthSearchFileController)
	authfilerouter.POST("createfile", controllers.AuthCreateFileController)
	authfilerouter.POST("getdetailfile", controllers.AuthGetFileDetailController)
	authfilerouter.POST("downloadfile", controllers.AuthDownloadFileController)
	authfilerouter.POST("deletefile", controllers.AuthRemoveFileController)
}
