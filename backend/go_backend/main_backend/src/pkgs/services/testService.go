package services

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TestService2SecretPassword struct {
	Test_password string `json:"test_password" validate:"required"`
}

type TestService2SecretPasswordInterface interface {
	Check_password() error
}

func (t TestService2SecretPassword) Check_password() error {
	var (
		test_password string = os.Getenv("TEST_PASSWORD")
	)

	if t.Test_password != test_password {
		return errors.New("테스트의 패스워드가 잘못되었습니다")
	}
	return nil
}

// 2번째 테스트 확인
func TestService2(ctx *gin.Context) (int, error) {
	var (
		errorStatus int
		err         error
	)

	if errorStatus, err = TestService2CheckClientAndPasswordFunc(ctx); err != nil {
		return errorStatus, err
	}

	return 0, nil
}

func TestService2CheckClientAndPasswordFunc(ctx *gin.Context) (int, error) {
	var (
		testPassword *TestService2SecretPassword
		err          error
	)

	// 클라이언트 폼을 파싱함
	if testPassword, err = TestService2ParseAndBodyFunc[TestService2SecretPassword](ctx); err != nil {
		return http.StatusBadRequest, err
	}

	// 비밀번호 확인
	if err = testPassword.Check_password(); err != nil {
		return http.StatusNotAcceptable, err
	}

	return 0, nil

}

func TestService2ParseAndBodyFunc[T TestService2SecretPassword](ctx *gin.Context) (*T, error) {
	var (
		body T
	)
	err := ctx.ShouldBindBodyWithJSON(&body)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(json) 클라이언트에서 보낸 폼을 확인하시오")
	}

	validate := validator.New()
	if err = validate.Struct(body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(validate) 클라이언트에서 보낸 폼을 확인하시오")
	}

	return &body, nil
}
