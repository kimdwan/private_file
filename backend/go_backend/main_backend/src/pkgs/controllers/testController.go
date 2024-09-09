package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/pkgs/services"
)

func TestController1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "연결이 완료되었습니다.",
	})
}

func TestController2(ctx *gin.Context) {
	var (
		errorStatus int
		err         error
	)
	if errorStatus, err = services.TestService2(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "2번째 테스트 완료",
	})
}
