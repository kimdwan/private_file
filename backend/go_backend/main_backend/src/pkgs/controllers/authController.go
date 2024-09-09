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
