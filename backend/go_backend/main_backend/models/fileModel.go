package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	User_id      uuid.UUID `gorm:"type:uuid;not null;"`
	File_id      uuid.UUID `gorm:"type:uuid;unique;"`
	File_title   string    `gorm:"type:varchar(255);not null;"`
	File_commnet string    `gorm:"type:varchar(5000);not null;"`
	File_path    string    `gorm:"type:varchar(500);not null;"`
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	f.User_id = uuid.New()

	return nil
}

func (File) TableName() string {
	return "LOGAN_FILE_TB"
}
