package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/dtos"
	"github.com/kimdwan/private_file/src/pkgs/services"
)

func AuthGetProfileImgController(ctx *gin.Context) {

	var (
		payload *dtos.Payload
		err     error
	)

	if payload, err = services.AuthParsePayloadByteService(ctx); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fmt.Println(payload)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "사진 보냈습니다.",
	})

}
