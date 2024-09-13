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

// 파일 내용 가져오기
func AuthGetFileListController(ctx *gin.Context) {
	var (
		payload         *dtos.Payload
		fileNumberDto   *dtos.FileNumberDto
		fileListDtos    []dtos.FileDataDto
		totalFileNumber int
		errorStatus     int
		err             error
	)

	// payload 가져오기
	if payload, err = services.AuthParsePayloadByteService(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 클라이언트에서 가져오기
	if fileNumberDto, err = services.AuthParseAndBodyService[dtos.FileNumberDto](ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 데이터 리스트 가져오기
	if errorStatus, err = services.AuthGetFileListService(payload, fileNumberDto, &fileListDtos, &totalFileNumber); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"file_datas":   fileListDtos,
		"total_number": totalFileNumber,
	})

}

// 파일을 검색해서 가져오기
func AuthSearchFileController(ctx *gin.Context) {
	var (
		payload          *dtos.Payload
		fileSearchDto    *dtos.FileSearchNameDto
		file_datas       []dtos.FileDataDto
		totalFileNumbers int
		errorStatus      int
		err              error
	)

	// payload 값 가져오기
	if payload, err = services.AuthParsePayloadByteService(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 파일 검색 후 가져오기
	if fileSearchDto, err = services.AuthParseAndBodyService[dtos.FileSearchNameDto](ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 파일 리스트 가져오가
	if errorStatus, err = services.AuthSearchFileService(payload, fileSearchDto, &file_datas, &totalFileNumbers); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"file_datas":   file_datas,
		"total_number": totalFileNumbers,
	})
}

func AuthCreateFileController(ctx *gin.Context) {
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

	// 데이터 저장
	if errorStatus, err = services.AuthCreateFileService(ctx, payload); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(errorStatus)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "데이터가 업로드 되었습니다.",
	})

}
