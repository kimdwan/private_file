package dtos

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// payload와 관련된 함수
type Payload struct {
	User_id uuid.UUID `json:"user_id"`
	Sub     Sub       `json:"sub"`
}

type Sub struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type PayloadInterFace interface {
	Make_jwttoken(jwt_secret string, jwt_time int) (string, error)
}

// payload로 jwt 토큰 만들기
func (p Payload) Make_jwttoken(jwt_secret string, jwt_time int) (string, error) {
	jwt_string := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Duration(jwt_time) * time.Second).Unix(),
		"payload": p,
	})

	jwt_token, err := jwt_string.SignedString([]byte(jwt_secret))
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return "", errors.New("jwt 토큰을 만드는데 오류가 발생했습니다")
	}

	return jwt_token, nil
}
