package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/dtos"
	"github.com/kimdwan/private_file/src/pkgs/services"
)

// 유저 회원가입이 이루어지는 로직
func UserSignUpController(ctx *gin.Context) {
	var (
		userSignUpDto *dtos.UserSignUp
		errorStatus   int
		err           error
	)

	// 유저의 클라이언트 폼을 파싱하는 함수
	if userSignUpDto, err = services.UserParseAndBodyService[dtos.UserSignUp](ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 유저의 정보를 데이터 베이스에 저장하는 로직
	if errorStatus, err = services.UserSignUpService(userSignUpDto); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "회원가입 되었습니다.",
	})
}
