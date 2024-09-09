package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kimdwan/private_file/models"
	"github.com/kimdwan/private_file/settings"
	"github.com/kimdwan/private_file/src/dtos"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// user에서 클라이언트 폼을 파싱하는 함수
func UserParseAndBodyService[T dtos.UserSignUp](ctx *gin.Context) (*T, error) {
	var (
		body T
		err  error
	)

	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(json) 클라이언트 폼을 파싱하는데 오류가 발생했습니다")
	}

	validate := validator.New()
	if err = validate.Struct(body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(validate) 클라이언트 폼을 파싱하는데 오류가 발생했습니다")
	}

	return &body, nil
}

// 유저의 회원가입이 이루어지는 함수
func UserSignUpService(signUpDto *dtos.UserSignUp) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 유저 클라이언트 폼 체크
	if errorStatus, err = UserSignUpCheckDtoFunc(c, db, signUpDto); err != nil {
		return errorStatus, err
	}

	// 데이터 저장
	if err = UserSignUpCreateDatabaseFunc(c, db, signUpDto); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 유저의 클라이언트 폼을 한번 체크하는 함수
func UserSignUpCheckDtoFunc(c context.Context, db *gorm.DB, signUpDto *dtos.UserSignUp) (int, error) {

	var (
		user models.User
	)

	// 비밀번호 폼 확인
	if err := signUpDto.Check_password(); err != nil {
		return http.StatusBadRequest, err
	}

	// 이메일 존재 유무 확인
	result := db.WithContext(c).Where("email = ?", signUpDto.Email).First(&user)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저의 이메일을 찾는중 오류가 발생했습니다")
		}
	} else {
		return http.StatusNotAcceptable, errors.New("이미 존재하는 이메일 입니다")
	}

	// 닉네임 존재 유무 확인
	if result = db.WithContext(c).Where("nickname = ?", signUpDto.Nickname).First(&user); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저의 닉네임을 찾는중 오류가 발생했습니다")
		}
	} else {
		return http.StatusFailedDependency, errors.New("이미 존재하는 닉네임 입니다")
	}

	return 0, nil
}

// 유저 데이터 저장
func UserSignUpCreateDatabaseFunc(c context.Context, db *gorm.DB, signUpDto *dtos.UserSignUp) error {

	var (
		new_user models.User
	)

	// 기본 데이터 저장
	new_user.Email = signUpDto.Email
	new_user.Nickname = signUpDto.Nickname
	new_user.Term_agree_3 = signUpDto.Term_agree_3

	// 비밀번호 해쉬화 후 저장
	var (
		saltRounds int = 10
		hash       []byte
		err        error
	)
	if hash, err = bcrypt.GenerateFromPassword([]byte(signUpDto.Password), saltRounds); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("비밀번호를 해쉬화 하는데 오류가 발생했습니다")
	}
	new_user.Hash = hash

	// 데이터 저장
	if result := db.WithContext(c).Create(&new_user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에 새로운 유저를 저장하는데 오류가 발생했습니다")
	}

	return nil
}
