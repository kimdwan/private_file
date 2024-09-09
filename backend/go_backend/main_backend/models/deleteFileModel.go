package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteFile struct {
	gorm.Model
	User_id      uuid.UUID `gorm:"type:uuid;not null;"`
	File_id      uuid.UUID `gorm:"type:uuid;not null;"`
	File_title   string    `gorm:"type:varchar(255);not null;"`
	File_comment string    `gorm:"type:varchar(5000);not null;"`
	File_path    string    `gorm:"type:varchar(500);not null;"`
}

func (DeleteFile) TableName() string {
	return "LOGAN_DELETE_FILE_TB"
}
