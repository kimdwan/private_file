package dtos

import (
	"errors"
	"regexp"
)

// 유저의 회원가입과 관련됨
type UserSignUp struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=4,max=16"`
	Nickname     string `json:"nickname" validate:"required,min=3,max=30"`
	Term_agree_3 bool   `json:"term_agree_3" validate:"boolean"`
}

type UserSignUpInterface interface {
	Check_password() error
}

// 비밀번호가 규정에 맞는지 확인하는 로직
func (u UserSignUp) Check_password() error {
	// 영문자 포함
	hash_letter := regexp.MustCompile(`[a-zA-Z]`).MatchString(u.Password)

	// 숫자 확인
	hash_number := regexp.MustCompile(`[0-9]`).MatchString(u.Password)

	// 특수문자 확인
	hash_special := regexp.MustCompile(`[!@#~$%^&*()_+\-=\[\]{};':"\\|,.<>/?]+`).MatchString(u.Password)

	if !hash_letter || !hash_number || !hash_special {
		return errors.New("비밀번호는 반드시 영문자, 숫자, 특수문자를 한개 이상 포함해야 합니다")
	}
	return nil
}

// 유저의 로그인과 관련됨
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=16"`
}
