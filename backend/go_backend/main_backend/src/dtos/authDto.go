package dtos

import "github.com/google/uuid"

// 파일 리스트에 들어가는 데이터
type FileDataDto struct {
	File_id      uuid.UUID `json:"file_id,omitempty"`
	File_name    string    `json:"file_name"`
	File_comment string    `json:"file_comment"`
}

// 파일 아이디를 만들고 싶으면 이렇게 만들면 됨
type FileDataDtoInterFace interface {
	MakeFileId()
}

func (f *FileDataDto) MakeFileId() {
	f.File_id = uuid.New()
}

// 파일의 순서를 가져오는 데이터
type FileNumberDto struct {
	File_number int `json:"file_number" validate:"number,min=1,required"`
}

// 파일 검색한 후 순서를 가져오는 데이터
type FileSearchNameDto struct {
	File_title  string `json:"file_title" validate:"required"`
	File_number int    `json:"file_number" validate:"required,min=1,required"`
}

// 파일 아이디를 가져오는 로직
type FileIdDto struct {
	File_id uuid.UUID `json:"file_id" validate:"required,uuid"`
}

// 파일 데이터 가져오는 로직
type FileDetailDataDto struct {
	File_title   string `json:"file_title"`
	File_comment string `json:"file_comment"`
	File_path    string `json:"file_path"`
}
