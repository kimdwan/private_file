package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	img_path := path.Join(file_server, image_server, user.User_id.String(), *user.Profile_img)

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

// 프로필 이미지를 업로드 하는 로직
func AuthUploadProfileService(ctx *gin.Context, payload *dtos.Payload) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		file_name   string
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 폼데이터 파싱
	formData, err := AuthUploadProfileGetFormDataFunc(ctx)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// 데이터 옮기고 저장하기
	if errorStatus, err = AuthUploadProfileCheckSizeAndSaveDataAndGetFileNameFunc(ctx, formData, payload, &file_name); err != nil {
		return errorStatus, err
	}

	// 데이터 베이스에 저장
	if err = AuthUploadProfileResetDatabaseFunc(c, db, payload, file_name); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 폼데이터 파싱
func AuthUploadProfileGetFormDataFunc(ctx *gin.Context) (*multipart.FileHeader, error) {
	formData, err := ctx.FormFile("img")

	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("form 데이터를 파싱하는데 오류가 발생했습니다")
	}

	return formData, nil
}

// 폼데이터 검증 후 파일 저장과 이름 채취
func AuthUploadProfileCheckSizeAndSaveDataAndGetFileNameFunc(ctx *gin.Context, formData *multipart.FileHeader, payload *dtos.Payload, file_name *string) (int, error) {

	// 파일 사이즈 확인
	if formData.Size > 10*1024*1024 {
		return http.StatusBadRequest, errors.New("파일 사이즈는 최대 10MB 입니다")
	}

	// 파일 타입 확인
	*file_name = formData.Filename
	var (
		system_img_types []string = strings.Split(os.Getenv("GO_IMAGE_TYPES"), ",")
		isTypeAllowed    bool     = false
	)
	file_name_list := strings.Split(*file_name, ".")
	for _, system_img_type := range system_img_types {
		if file_name_list[len(file_name_list)-1] == system_img_type {
			isTypeAllowed = true
			break
		}
	}
	if !isTypeAllowed {
		return http.StatusBadRequest, errors.New("파일 타입을 다시 확인해주시길 바랍니다")
	}

	// 파일 경로
	var (
		file_server  string = os.Getenv("FILE_SERVER_PATH")
		image_server string = os.Getenv("FILE_PROFILE_SERVER_PATH")
	)
	user_profile_server := path.Join(file_server, image_server, payload.User_id.String())

	// 기존 파일 검색
	var origin_files []string
	err := filepath.Walk(user_profile_server, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		origin_files = append(origin_files, path)
		return nil
	})
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("기존 파일을 검색하는데 오류가 발생했습니다")
	}

	// 기존 파일 옮기기
	var (
		currentDate   string = time.Now().Format("2006-01-02")
		currentTime   string = time.Now().Format("15:04:05")
		remove_server string = os.Getenv("FILE_REMOVE_PROFILE_SERVER_PATH")
	)
	user_remove_server := path.Join(file_server, remove_server, payload.User_id.String(), currentDate+"T"+currentTime)

	for _, origin_file := range origin_files {
		origin_file_path_list := strings.Split(origin_file, "/")
		user_remove_server_path := path.Join(user_remove_server, origin_file_path_list[len(origin_file_path_list)-1])
		if err = os.Rename(origin_file, user_remove_server_path); err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			return http.StatusInternalServerError, errors.New("파일을 옮기는데 오류가 발생했습니다")
		}
	}

	// 파일 저장하기
	file_path := path.Join(user_profile_server, *file_name)
	if err = ctx.SaveUploadedFile(formData, file_path); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("파일을 저장하는데 오류가 발생했습니다")
	}

	return 0, nil
}

// 데이터 베이스 업로드
func AuthUploadProfileResetDatabaseFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, file_name string) error {

	var (
		user models.User
	)
	// 유저 데이터 찾기
	result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(&user)
	if result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("유저 아이디에 해당하는 데이터를 찾는데 오류가 발생했습니다")
	}

	// 데이터 수정
	user.Profile_img = &file_name
	if result = db.WithContext(c).Save(&user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스를 업데이트 하는데 오류가 발생했습니다")
	}

	return nil
}
