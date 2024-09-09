package models

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	User_id         uuid.UUID `gorm:"type:uuid;unique;"`
	Email           string    `gorm:"type:varchar(255);unique;not null;"`
	Hash            []byte    `gorm:"type:bytea;not null;"`
	Nickname        string    `gorm:"type:varchar(255);unique;not null;"`
	Term_agree_3    bool      `gorm:"type:boolean;not null;"`
	Profile_img     *string   `gorm:"type:varchar(1000);"`
	Access_token    *string   `gorm:"type:varchar(1000);unique"`
	Refresh_token   *string   `gorm:"type:varchar(1000);unique"`
	Computer_number *string   `gorm:"type:uuid;unique"`

	File []File `gorm:"foreignKey:user_id;references:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// 생성될때 해야 할거
func (u *User) BeforeCreate(tx *gorm.DB) error {

	// 값 제공
	u.User_id = uuid.New()

	// 이메일 확인 로직
	var (
		emailMatchString string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	)
	reg := regexp.MustCompile(emailMatchString)
	isEmailAllowed := reg.MatchString(u.Email)
	if !isEmailAllowed {
		return errors.New("이메일 형식을 지켜주시길 바랍니다")
	}

	return nil
}

// 수정할때 해야할것
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Profile_img != nil {
		var (
			system_allowed_types        []string = strings.Split(os.Getenv("GO_IMAGE_TYPES"), ",")
			userProfile_image_name_list []string = strings.Split(*u.Profile_img, ".")
			userProfile_image_type      string   = strings.ToLower(userProfile_image_name_list[len(userProfile_image_name_list)-1])
			isProfileImgAllowed         bool     = false
		)

		for _, system_allowed_type := range system_allowed_types {
			if userProfile_image_type == system_allowed_type {
				isProfileImgAllowed = true
				break
			}
		}

		if !isProfileImgAllowed {
			return errors.New("이미지의 타입을 다시 확인하세요")
		}
	}

	return nil
}

func (User) TableName() string {
	return "LOGAN_USER_TB"
}
