package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/dtos"
	"github.com/kimdwan/private_file/src/pkgs/services"
)

// 유저의 프로필 이미지를 주는 컨트롤러
func AuthGetProfileImgController(ctx *gin.Context) {

	var (
		payload     *dtos.Payload
		imageDto    dtos.ImageDto
		errorStatus int
		err         error
	)

	// payload를 파싱함
	if payload, err = services.AuthParsePayloadByteService(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// image dto 파싱하기
	if errorStatus, err = services.AuthGetProfileImgService(payload, &imageDto); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusOK, imageDto)
}

// 닉네임 가져오기
func AuthGetNickNameController(ctx *gin.Context) {
	var (
		payload *dtos.Payload
		err     error
	)

	// payload 파싱하기
	if payload, err = services.AuthParsePayloadByteService(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"nickname": payload.Sub.Nickname,
	})
}

// 프로필 이미지 업로드
func AuthUploadProfileController(ctx *gin.Context) {
	var (
		payload     *dtos.Payload
		errorStatus int
		err         error
	)

	// payload 파싱하기
	if payload, err = services.AuthParsePayloadByteService(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 파일
	if errorStatus, err = services.AuthUploadProfileService(ctx, payload); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "프로필이 업데이트 되었습니다.",
	})

}
