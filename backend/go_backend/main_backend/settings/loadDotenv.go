package settings

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadDotenv() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		panic("환경변수를 로드하는데 오류가 발생했습니다.")
	}
}
