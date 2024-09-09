package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/models"
	"github.com/kimdwan/private_file/settings"
	"github.com/kimdwan/private_file/src/dtos"
	"gorm.io/gorm"
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

// 유저의 프로필 이미지를 주는 함수
func AuthGetProfileImgService(payload *dtos.Payload, imageDto *dtos.ImageDto) (int, error) {

	var (
		db   *gorm.DB = settings.DB
		user models.User
		err  error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 데이터 베이스에서 user 모델 가져오기
	if err = AuthGetProfileImgGetDataFunc(c, db, &user, payload); err != nil {
		return http.StatusInternalServerError, err
	}

	// 프로필 이미지 확인
	if user.Profile_img == nil {
		return 0, nil
	}

	// 파일 가져오기
	if err = AuthGetProfileImgParseAndSendImgDataFunc(&user, imageDto); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 데이터 베이스에서 profile 이미지가 존재하는지 확인하고 가져오기
func AuthGetProfileImgGetDataFunc(c context.Context, db *gorm.DB, user *models.User, payload *dtos.Payload) error {

	// 데이터 베이스에서 user_id에 해당하는 데이터 찾기
	if result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에서 user_id에 해당하는 데이터를 찾는데 오류가 발생했습니다")
	}

	return nil
}

// 이미지가 있다면 보내기
func AuthGetProfileImgParseAndSendImgDataFunc(user *models.User, imageDto *dtos.ImageDto) error {

	// 주소 가져오기
	var (
		file_server  string = os.Getenv("FILE_SERVER_PATH")
		image_server string = os.Getenv("FILE_PROFILE_SERVER_PATH")
	)
	img_path := path.Join(file_server, image_server, *user.Profile_img)

	// 파일 가져오기
	images, err := os.ReadFile(img_path)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("이미지를 읽는데 오류가 발생했습니다")
	}

	// 파일을 base64로 인코딩
	image_file_str := base64.StdEncoding.EncodeToString(images)

	// 파일의 타입 가져오기
	file_path_list := strings.Split(*user.Profile_img, ".")
	file_type := file_path_list[len(file_path_list)-1]

	// 보내기
	imageDto.Imagebase64 = image_file_str
	imageDto.Imagetype = file_type

	return nil
}
