package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kimdwan/private_file/models"
	"github.com/kimdwan/private_file/settings"
	"github.com/kimdwan/private_file/src/dtos"
	"gorm.io/gorm"
)

// jwt 토큰의 유효성을 검사하는 미들웨어
func CheckJwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var (
			access_token    string
			computer_number uuid.UUID
			errorStatus     int
			err             error
		)

		// 컴퓨터 number와 access token을 검사
		if errorStatus, err = CheckJwtGetAccessTokenAndComputerNumberFunc(ctx, &access_token, &computer_number); err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatus(errorStatus)
			return
		}

		// jwt 토큰 검증
		var (
			payload         dtos.Payload
			jwt_secret_keys = []string{
				os.Getenv("JWT_ACCESS_SECRET"),
				os.Getenv("JWT_REFRESH_SECRET"),
			}
			counts int = 0
		)
		// access token 확인
		if errorStatus, err = CheckJwtCheckJwtTokenFunc(&payload, access_token, jwt_secret_keys[0], &counts); err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatus(errorStatus)
			return
		}

		// refresh token 검증
		if counts > 0 {
			fmt.Println("토큰 시간이 만료했습니다 REFRESH 검증으로 들어갑니다")
			var (
				db               *gorm.DB = settings.DB
				user             models.User
				refresh_token    string
				new_access_token string
			)
			c, cancel := context.WithTimeout(context.Background(), time.Second*100)
			defer cancel()

			// refresh token 가져오기
			if errorStatus, err = CheckJwtGetRefreshTokenFunc(c, db, &user, computer_number, &refresh_token); err != nil {
				fmt.Println(err.Error())
				ctx.AbortWithStatus(errorStatus)
				return
			}

			// refresh token 검증
			if errorStatus, err = CheckJwtCheckJwtTokenFunc(&payload, refresh_token, jwt_secret_keys[1], &counts); err != nil {
				fmt.Println(err.Error())
				ctx.AbortWithStatus(errorStatus)
				return
			}

			// 새로운 access token 발급
			if err = CheckJwtNewAccessTokenFunc(&payload, jwt_secret_keys[0], &new_access_token); err != nil {
				fmt.Println(err.Error())
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// 데이터 베이스에 저장 cookie 제공
			if err = CheckJwtSaveAccessTokenAndSendFunc(c, ctx, db, &user, new_access_token); err != nil {
				fmt.Println(err.Error())
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		// 새로운 payload_byte 생성 후 전달
		if err = CheckJwtMakePayloadByteFunc(ctx, &payload); err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Next()
	}
}

// access token과 computer number를 파싱함
func CheckJwtGetAccessTokenAndComputerNumberFunc(ctx *gin.Context, access_token *string, computer_number *uuid.UUID) (int, error) {

	var (
		err error
	)

	// access token 가져오기
	if *access_token, err = ctx.Cookie("Authorization"); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusUnauthorized, errors.New("cookie에 있는 데이터 베이스를 파싱하는데 오류가 발생했습니다")
	}

	// computer number 가져오기
	header_compuer_number := ctx.GetHeader("User-Computer-Number")
	if header_compuer_number == "" {
		return http.StatusBadRequest, errors.New("User-Computer-Number를 입력하지 않았습니다")
	}

	// computer number 파싱
	*computer_number = uuid.MustParse(header_compuer_number)

	return 0, nil
}

// jwt token 검증
func CheckJwtCheckJwtTokenFunc(payload *dtos.Payload, jwt_token string, jwt_secret_key string, counts *int) (int, error) {

	if *counts > 1 {
		return http.StatusUnauthorized, errors.New("refresh token도 문제가 있습니다")
	}

	// jwt token 읽어오기
	tokens, err := jwt.Parse(jwt_token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt 해석방식이 HMAC 방식이 아닙니다")
		}
		return []byte(jwt_secret_key), nil
	})

	// 시간 체크
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			*counts += 1
			return 0, nil
		} else {
			fmt.Println("시스템 오류: ", err.Error())
			return http.StatusInternalServerError, errors.New("jwt 토큰을 파싱하는데 오류가 발생했습니다")
		}
	}

	// payload 가져오기
	if claims, ok := tokens.Claims.(jwt.MapClaims); ok {
		if payload_str, ok := claims["payload"]; ok {
			payload_byte, err := json.Marshal(&payload_str)

			if err != nil {
				fmt.Println("시스템 오류: ", err.Error())
				return http.StatusInternalServerError, errors.New("payload를 byte화 시키는데 오류가 발생했습니다")
			}

			if err = json.Unmarshal(payload_byte, payload); err != nil {
				fmt.Println("시스템 오류: ", err.Error())
				return http.StatusInternalServerError, errors.New("payload를 역직렬화 하는데 오류가 발생했습니다")
			}

		} else {
			return http.StatusUnauthorized, errors.New("claims에 payload를 파싱하는데 오류가 발생했습니다")
		}
	} else {
		return http.StatusUnauthorized, errors.New("jwt token에 claims를 파싱하는데 오류가 발생했습니다")
	}

	return 0, nil
}

// refresh token 가져오기
func CheckJwtGetRefreshTokenFunc(c context.Context, db *gorm.DB, user *models.User, computer_number uuid.UUID, refresh_token *string) (int, error) {

	// 데이터 베이스에서 유저 정보 가져오기
	result := db.WithContext(c).Where("computer_number = ?", computer_number).First(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusUnauthorized, errors.New("보낸 computer number는 잘못되었습니다")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 computer number를 찾는데 오류가 발생했습니다")
		}
	}

	*refresh_token = *user.Refresh_token
	return 0, nil
}

// 새로운 access token을 발급받는 함수
func CheckJwtNewAccessTokenFunc(payload *dtos.Payload, access_secret string, new_access_token *string) error {

	var (
		access_time_str string = os.Getenv("JWT_ACCESS_TIME")
	)
	access_time, err := strconv.Atoi(access_time_str)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("시간을 문자에서 숫자로 바꾸는데 오류가 발생했습니다")
	}

	jwt_token_str := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Duration(access_time) * time.Second).Unix(),
		"payload": payload,
	})

	if *new_access_token, err = jwt_token_str.SignedString([]byte(access_secret)); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("새로운 jwt 토큰을 만드는데 오류가 발생했습니다")
	}

	return nil
}

// 새로운 데이터 베이스에 저장하고 cookie도 제출
func CheckJwtSaveAccessTokenAndSendFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *models.User, new_acess_token string) error {

	// 새롭게 업로드
	user.Access_token = &new_acess_token
	if result := db.WithContext(c).Save(user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에 user 정보를 업데이트 하는데 오류가 발생했습니다")
	}

	// 쿠키에 업로드
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", new_acess_token, 24*60*60, "", "", false, true)

	return nil
}

// payload byte 생성 후 전달
func CheckJwtMakePayloadByteFunc(ctx *gin.Context, payload *dtos.Payload) error {

	// payload를 직렬화 하기
	payload_byte, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("payload를 직렬화 하는데 오류가 발생했습니다")
	}

	// payload를 보내기
	ctx.Set("payload_byte", string(payload_byte))
	return nil
}
