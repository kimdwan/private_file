package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/src/dtos"
)

// middleware에서 보낸 payload를 받는 함수
func AuthParsePayloadByteService(ctx *gin.Context) (*dtos.Payload, error) {

	var (
		payload_byte string = ctx.GetString("payload_byte")
		payload      dtos.Payload
		err          error
	)

	if payload_byte == "" {
		return nil, errors.New("payload byte를 전달받지 못했습니다")
	}

	if err = json.Unmarshal([]byte(payload_byte), &payload); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("payload를 역직렬화 하는데 오류가 발생했습니다")
	}

	return &payload, nil
}
