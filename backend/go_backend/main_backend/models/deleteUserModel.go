package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteUser struct {
	gorm.Model
	User_id      uuid.UUID `gorm:"type:uuid;not null;"`
	Email        string    `gorm:"type:varchar(255);not null;"`
	Nickname     string    `gorm:"type:varchar(255);not null;"`
	Profile_img  *string   `gorm:"type:varchar(1000)"`
	Term_agree_3 bool      `gorm:"type:boolean;not null;"`
}

func (DeleteUser) TableName() string {
	return "LOGAN_DELETE_USER_TB"
}
