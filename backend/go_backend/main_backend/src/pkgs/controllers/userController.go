package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// 유저 로그인이 이루어지는 로직
func UserLoginController(ctx *gin.Context) {

	var (
		loginDto        *dtos.UserLogin
		access_token    string
		computer_number uuid.UUID
		errorStatus     int
		err             error
	)

	// 유저의 클라이언트 폼을 파싱
	if loginDto, err = services.UserParseAndBodyService[dtos.UserLogin](ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 보안 정보 저장
	if errorStatus, err = services.UserLoginService(loginDto, &access_token, &computer_number); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	// 토큰 전달
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", access_token, 24*60*60, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"computer_number": computer_number.String(),
	})
}
