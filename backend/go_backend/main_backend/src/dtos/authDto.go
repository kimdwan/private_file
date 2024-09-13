package dtos

import "github.com/google/uuid"

// 파일 리스트에 들어가는 데이터
type FileDataDto struct {
	File_id      uuid.UUID `json:"file_id"`
	File_name    string    `json:"file_name"`
	File_comment string    `json:"file_comment"`
}

// 파일의 순서를 가져오는 데이터
type FileNumberDto struct {
	File_number int `json:"file_number" validate:"number,min=1,required"`
}
