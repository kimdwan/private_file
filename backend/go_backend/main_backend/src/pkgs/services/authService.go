package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kimdwan/private_file/models"
	"github.com/kimdwan/private_file/settings"
	"github.com/kimdwan/private_file/src/dtos"
	"gorm.io/gorm"
)

// middleware에서 보낸 payload를 받는 함수
func AuthParsePayloadByteService(ctx *gin.Context) (*dtos.Payload, error) {

	var (
		payload_byte string = ctx.GetString("payload_byte")
		payload      dtos.Payload
		err          error
	)

	if payload_byte == "" {
		return nil, errors.New("payload byte를 전달받지 못했습니다")
	}

	if err = json.Unmarshal([]byte(payload_byte), &payload); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("payload를 역직렬화 하는데 오류가 발생했습니다")
	}

	return &payload, nil
}

// auth에서 원하는 데이터 파싱해줌
func AuthParseAndBodyService[T dtos.FileNumberDto | dtos.FileSearchNameDto | dtos.FileIdDto](ctx *gin.Context) (*T, error) {
	var (
		body T
		err  error
	)

	// 폼 가져오기
	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(json) 클라이언트 폼을 파싱하는데 오류가 발생했습니다")
	}

	// 폼 검증하기
	validate := validator.New()
	if err = validate.Struct(body); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("(validate) 클라이언트 폼을 파싱하는데 오류가 발생했습니다")
	}

	return &body, nil
}

// 유저의 프로필 이미지를 주는 함수
func AuthGetProfileImgService(payload *dtos.Payload, imageDto *dtos.ImageDto) (int, error) {

	var (
		db   *gorm.DB = settings.DB
		user models.User
		err  error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 데이터 베이스에서 user 모델 가져오기
	if err = AuthGetProfileImgGetDataFunc(c, db, &user, payload); err != nil {
		return http.StatusInternalServerError, err
	}

	// 프로필 이미지 확인
	if user.Profile_img == nil {
		return 0, nil
	}

	// 파일 가져오기
	if err = AuthGetProfileImgParseAndSendImgDataFunc(&user, imageDto); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 데이터 베이스에서 profile 이미지가 존재하는지 확인하고 가져오기
func AuthGetProfileImgGetDataFunc(c context.Context, db *gorm.DB, user *models.User, payload *dtos.Payload) error {

	// 데이터 베이스에서 user_id에 해당하는 데이터 찾기
	if result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스에서 user_id에 해당하는 데이터를 찾는데 오류가 발생했습니다")
	}

	return nil
}

// 이미지가 있다면 보내기
func AuthGetProfileImgParseAndSendImgDataFunc(user *models.User, imageDto *dtos.ImageDto) error {

	// 주소 가져오기
	var (
		file_server  string = os.Getenv("FILE_SERVER_PATH")
		image_server string = os.Getenv("FILE_PROFILE_SERVER_PATH")
	)
	img_path := path.Join(file_server, image_server, user.User_id.String(), *user.Profile_img)

	// 파일 가져오기
	images, err := os.ReadFile(img_path)
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("이미지를 읽는데 오류가 발생했습니다")
	}

	// 파일을 base64로 인코딩
	image_file_str := base64.StdEncoding.EncodeToString(images)

	// 파일의 타입 가져오기
	file_path_list := strings.Split(*user.Profile_img, ".")
	file_type := file_path_list[len(file_path_list)-1]

	// 보내기
	imageDto.Imagebase64 = image_file_str
	imageDto.Imagetype = file_type

	return nil
}

// 프로필 이미지를 업로드 하는 로직
func AuthUploadProfileService(ctx *gin.Context, payload *dtos.Payload) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		file_name   string
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 폼데이터 파싱
	formData, err := AuthUploadProfileGetFormDataFunc(ctx)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// 데이터 옮기고 저장하기
	if errorStatus, err = AuthUploadProfileCheckSizeAndSaveDataAndGetFileNameFunc(ctx, formData, payload, &file_name); err != nil {
		return errorStatus, err
	}

	// 데이터 베이스에 저장
	if err = AuthUploadProfileResetDatabaseFunc(c, db, payload, file_name); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 폼데이터 파싱
func AuthUploadProfileGetFormDataFunc(ctx *gin.Context) (*multipart.FileHeader, error) {
	formData, err := ctx.FormFile("img")

	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return nil, errors.New("form 데이터를 파싱하는데 오류가 발생했습니다")
	}

	return formData, nil
}

// 폼데이터 검증 후 파일 저장과 이름 채취
func AuthUploadProfileCheckSizeAndSaveDataAndGetFileNameFunc(ctx *gin.Context, formData *multipart.FileHeader, payload *dtos.Payload, file_name *string) (int, error) {

	// 파일 사이즈 확인
	if formData.Size > 10*1024*1024 {
		return http.StatusBadRequest, errors.New("파일 사이즈는 최대 10MB 입니다")
	}

	// 파일 타입 확인
	*file_name = formData.Filename
	var (
		system_img_types []string = strings.Split(os.Getenv("GO_IMAGE_TYPES"), ",")
		isTypeAllowed    bool     = false
	)
	file_name_list := strings.Split(*file_name, ".")
	for _, system_img_type := range system_img_types {
		if file_name_list[len(file_name_list)-1] == system_img_type {
			isTypeAllowed = true
			break
		}
	}
	if !isTypeAllowed {
		return http.StatusBadRequest, errors.New("파일 타입을 다시 확인해주시길 바랍니다")
	}

	// 파일 경로
	var (
		file_server  string = os.Getenv("FILE_SERVER_PATH")
		image_server string = os.Getenv("FILE_PROFILE_SERVER_PATH")
	)
	user_profile_server := path.Join(file_server, image_server, payload.User_id.String())

	// 기존 파일 검색
	var origin_files []string
	err := filepath.Walk(user_profile_server, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		origin_files = append(origin_files, path)
		return nil
	})
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("기존 파일을 검색하는데 오류가 발생했습니다")
	}

	// 기존 파일 옮기기
	var (
		currentDate   string = time.Now().Format("2006-01-02")
		currentTime   string = time.Now().Format("15:04:05")
		remove_server string = os.Getenv("FILE_REMOVE_PROFILE_SERVER_PATH")
	)
	user_remove_server := path.Join(file_server, remove_server, payload.User_id.String(), currentDate+"T"+currentTime)

	for _, origin_file := range origin_files {
		origin_file_path_list := strings.Split(origin_file, "/")
		user_remove_server_path := path.Join(user_remove_server, origin_file_path_list[len(origin_file_path_list)-1])
		if err = os.Rename(origin_file, user_remove_server_path); err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			return http.StatusInternalServerError, errors.New("파일을 옮기는데 오류가 발생했습니다")
		}
	}

	// 파일 저장하기
	file_path := path.Join(user_profile_server, *file_name)
	if err = ctx.SaveUploadedFile(formData, file_path); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("파일을 저장하는데 오류가 발생했습니다")
	}

	return 0, nil
}

// 데이터 베이스 업로드
func AuthUploadProfileResetDatabaseFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, file_name string) error {

	var (
		user models.User
	)
	// 유저 데이터 찾기
	result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(&user)
	if result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("유저 아이디에 해당하는 데이터를 찾는데 오류가 발생했습니다")
	}

	// 데이터 수정
	user.Profile_img = &file_name
	if result = db.WithContext(c).Save(&user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("데이터 베이스를 업데이트 하는데 오류가 발생했습니다")
	}

	return nil
}

// 파일 리스트 가져오기
func AuthGetFileListService(payload *dtos.Payload, fileNumberDto *dtos.FileNumberDto, fileListDtos *[]dtos.FileDataDto, totalFileNumber *int) (int, error) {

	var (
		db    *gorm.DB = settings.DB
		files []models.File
		err   error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 데이터 정보 가져오기
	if err = AuthGetFileListFindDatabaseFunc(c, db, &files, payload, totalFileNumber); err != nil {
		return http.StatusInternalServerError, err
	}

	// 갯수가 0개면 넘어가기
	if *totalFileNumber == 0 {
		return 0, nil
	}

	// 파일 가져오기
	if err = AuthGetFileListGetFileAndSummaryFunc(&files, fileListDtos, fileNumberDto, totalFileNumber); err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// 데이터 베이스에서 유저 찾기
func AuthGetFileListFindDatabaseFunc(c context.Context, db *gorm.DB, files *[]models.File, payload *dtos.Payload, totalFileNumber *int) error {

	// 파일 데이터 찾기
	if result := db.WithContext(c).Where("user_id = ?", payload.User_id).Order("created_at DESC").Find(files); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return errors.New("유저 아이디에 해당하는 파일 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	// 전체 갯수 가져오기
	*totalFileNumber = len(*files)

	return nil
}

// 파일 가져오고 정리하기
func AuthGetFileListGetFileAndSummaryFunc(files *[]models.File, fileListDtos *[]dtos.FileDataDto, fileNumberDto *dtos.FileNumberDto, totalFileNumber *int) error {
	var (
		initNumber      int = (fileNumberDto.File_number - 1) * 10
		last_number_ex1 int = (fileNumberDto.File_number) * 10
		last_number     int
	)

	// 숫자 오류를 막기 위한 확인 방법
	if *totalFileNumber <= initNumber {
		return errors.New("클라이언트에서 보낸 숫자 데이터가 문제가 있습니다")
	}

	if last_number_ex1 < *totalFileNumber {
		last_number = last_number_ex1
	} else {
		last_number = *totalFileNumber
	}

	want_file_datas := (*files)[initNumber:last_number]

	// 데이터를 배정하기
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	wg.Add(1)
	go AuthGetFileListSupplySummaryFunc(&wg, &mutex, &want_file_datas, fileListDtos)
	wg.Wait()

	return nil
}

// 데이터를 빠르게 배정하기 위해 준비한 함수
func AuthGetFileListSupplySummaryFunc(wg *sync.WaitGroup, mutex *sync.Mutex, want_files *[]models.File, fileListDtos *[]dtos.FileDataDto) {
	defer wg.Done()

	for _, want_file := range *want_files {
		var (
			fileListDto dtos.FileDataDto
		)
		mutex.Lock()
		fileListDto.File_id = want_file.File_id
		fileListDto.File_name = want_file.File_title
		fileListDto.File_comment = want_file.File_commnet
		*fileListDtos = append(*fileListDtos, fileListDto)
		mutex.Unlock()
	}

}

// 파일에서 이름을 검색후 리스트 가져오기
func AuthSearchFileService(payload *dtos.Payload, fileSearchDto *dtos.FileSearchNameDto, file_datas *[]dtos.FileDataDto, totalFileNumbers *int) (int, error) {
	var (
		db    *gorm.DB = settings.DB
		files []models.File
		err   error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 파일 이름 확인하는 로직
	if err = AuthSearchFileGetFileDatasFunc(c, db, payload, fileSearchDto, &files, totalFileNumbers); err != nil {
		return http.StatusInternalServerError, err
	}

	// 갯수가 0개면 넘어간다
	if *totalFileNumbers == 0 {
		return 0, nil
	}

	// 데이터 가져오기
	if err = AuthSearchFileGetFilesFunc(&files, file_datas, fileSearchDto, totalFileNumbers); err != nil {
		return http.StatusBadRequest, err
	}

	return 0, nil
}

// 파일 이름을 확인하는 로직
func AuthSearchFileGetFileDatasFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, fileSearchDto *dtos.FileSearchNameDto, files *[]models.File, totalFileNumbers *int) error {

	if result := db.WithContext(c).Where("user_id = ? AND file_title LIKE ?", payload.User_id, fileSearchDto.File_title+"%").Order("created_at DESC").Find(files); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return errors.New("데이터 베이스에서 원하는 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	*totalFileNumbers = len(*files)
	return nil
}

// 파일 확인후 보내주기
func AuthSearchFileGetFilesFunc(files *[]models.File, file_datas *[]dtos.FileDataDto, fileSearchDto *dtos.FileSearchNameDto, totalFileNumbers *int) error {
	var (
		initNumber    int = (fileSearchDto.File_number - 1) * 10
		lastNumber_ex int = (fileSearchDto.File_number) * 10
		lastNumber    int
	)

	if *totalFileNumbers <= initNumber {
		return errors.New("클라이언트 에서 보낸 숫자 데이터를 다시 확인하세요")
	}

	// 갯수확인
	if *totalFileNumbers > lastNumber_ex {
		lastNumber = lastNumber_ex
	} else {
		lastNumber = *totalFileNumbers
	}

	// 데이터 가져오기
	want_files := (*files)[initNumber:lastNumber]
	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	wg.Add(1)
	go AuthSearchFileSupplyGetFilesFunc(&wg, &mutex, &want_files, file_datas)
	wg.Wait()

	return nil
}

// 파일 확인하는데 필요한 함수
func AuthSearchFileSupplyGetFilesFunc(wg *sync.WaitGroup, mutex *sync.Mutex, want_files *[]models.File, file_datas *[]dtos.FileDataDto) {
	defer wg.Done()

	for _, want_file := range *want_files {
		var (
			file_data dtos.FileDataDto
		)
		mutex.Lock()
		file_data.File_id = want_file.File_id
		file_data.File_name = want_file.File_title
		file_data.File_comment = want_file.File_commnet
		*file_datas = append(*file_datas, file_data)
		mutex.Unlock()
	}

}

// formData파싱하기
func AuthCreateFileService(ctx *gin.Context, payload *dtos.Payload) (int, error) {
	var (
		fileDataDto dtos.FileDataDto
		file_names  string
		errorStatus int
		err         error
	)

	// 파일에서 데이터를 추출하고 파일을 저장한다
	if errorStatus, err = AuthCreateFileGetFormDataFunc(ctx, payload, &fileDataDto, &file_names); err != nil {
		return errorStatus, err
	}

	// 데이터를 저장한다
	var (
		db *gorm.DB = settings.DB
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	if err = AuthCreateFileSaveDataFunc(c, db, payload, &fileDataDto, &file_names); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 파일에서 데이터를 추출하고 파일을 저장한다
func AuthCreateFileGetFormDataFunc(ctx *gin.Context, payload *dtos.Payload, fileDataDto *dtos.FileDataDto, file_names *string) (int, error) {

	// 폼데이터 파싱하기
	formDatas, err := ctx.MultipartForm()
	if err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusBadRequest, errors.New("폼 데이터를 파싱하는데 오류가 발생했습니다")
	}

	// 파일 정보 파싱하고 대입하기
	var (
		clientData map[string]string = map[string]string{}
		json_names                   = []string{
			"file_name",
			"file_comment",
		}
	)
	for idx, formName := range []string{"title", "comment"} {
		if fromValue, ok := formDatas.Value[formName]; ok && len(fromValue) > 0 {
			clientData[json_names[idx]] = fromValue[0]
		}
	}
	file_datas_byte, err := json.Marshal(&clientData)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusInternalServerError, errors.New("파일 데이터를 바이트화 하는데 오류가 발생했습니다")
	}
	if err = json.Unmarshal(file_datas_byte, fileDataDto); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return http.StatusInternalServerError, errors.New("파일 데이터를 역직렬화 하는데 오류가 발생했습니다")
	}
	fileDataDto.MakeFileId()

	// 파일 데이터 파싱하기
	if formFile, ok := formDatas.File["file"]; ok && len(formFile) > 0 {
		formFileData := formFile[0]

		// 파일 크기 확인
		if formFileData.Size > 100*1024*1024 {
			return http.StatusBadRequest, errors.New("파일 데이터는 100mb를 넘을수 없습니다")
		}

		// 파일 이름 배정
		*file_names = formFileData.Filename

		// 파일 저장
		var (
			file_server      string = os.Getenv("FILE_SERVER_PATH")
			file_data_server string = os.Getenv("FILE_DATA_SERVER_PATH")
		)
		file_path := filepath.Join(file_server, file_data_server, payload.User_id.String(), fileDataDto.File_id.String(), *file_names)

		if err = ctx.SaveUploadedFile(formFileData, file_path); err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			return http.StatusInternalServerError, errors.New("파일을 저장하는데 오류가 발생했습니다")
		}
	} else {
		return http.StatusBadRequest, errors.New("폼 데이터에 file 데이터가 존재하지 않습니다")
	}

	return 0, nil
}

// 파일 데이터 배이스를 저장하는 로직
func AuthCreateFileSaveDataFunc(c context.Context, db *gorm.DB, paylaod *dtos.Payload, fileDataDto *dtos.FileDataDto, file_names *string) error {
	var (
		file models.File
	)

	file.File_id = fileDataDto.File_id
	file.File_title = fileDataDto.File_name
	file.File_commnet = fileDataDto.File_comment
	file.File_path = *file_names
	file.User_id = paylaod.User_id
	if result := db.WithContext(c).Create(&file); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("새로운 파일 데이터를 생성하는데 오류가 발생했습니다")
	}

	return nil
}

// 파일의 디테일한 부분을 가져오는 함수
func AuthGetFileDetailService(fileIdDto *dtos.FileIdDto, detailFileDatas *dtos.FileDetailDataDto) (int, error) {
	var (
		db *gorm.DB = settings.DB
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 파일의 세부정보를 가져오는 로직
	if errorStatus, err := AuthGetFileDetailGetDataFunc(c, db, fileIdDto, detailFileDatas); err != nil {
		return errorStatus, err
	}

	return 0, nil
}

// 파일의 디테일한 부분을 가져오는 로직
func AuthGetFileDetailGetDataFunc(c context.Context, db *gorm.DB, fileIdDto *dtos.FileIdDto, detailFileDatas *dtos.FileDetailDataDto) (int, error) {
	var (
		file models.File
	)

	// 데이터 찾기
	if result := db.WithContext(c).Where("file_id = ?", fileIdDto.File_id).First(&file); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, errors.New("클라이언트에서 보낸 file_id는 존재하지 않습니다")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 파일 아이디에 해당하는 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	// 데이터 옮기기
	detailFileDatas.File_title = file.File_title
	detailFileDatas.File_comment = file.File_commnet
	detailFileDatas.File_path = file.File_path

	return 0, nil
}

// 파일을 다운로드 하는 로직
func AuthDownloadFileService(ctx *gin.Context, fileIdDto *dtos.FileIdDto) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		file        models.File
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 파일의 주소를 가져오기 위함
	if errorStatus, err = AuthDownloadFileGetFileDataFunc(c, db, &file, fileIdDto); err != nil {
		return errorStatus, err
	}

	// 파일을 다운로드 함
	AuthDownloadFileStartDataFunc(ctx, &file)

	return 0, nil
}

// 파일 데이터의 주소를 일단 가져온다
func AuthDownloadFileGetFileDataFunc(c context.Context, db *gorm.DB, file *models.File, fileIdDto *dtos.FileIdDto) (int, error) {

	// 데이터 에서 찾아보기
	if result := db.WithContext(c).Where("file_id = ?", fileIdDto.File_id).First(file); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, errors.New("클라이언트에서 보낸 파일 아이디를 다시 확인하세요")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 파일 아이디에 해당하는 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	return 0, nil
}

// 파일을 다운로드 함
func AuthDownloadFileStartDataFunc(ctx *gin.Context, file *models.File) {

	var (
		file_server      string = os.Getenv("FILE_SERVER_PATH")
		file_data_server string = os.Getenv("FILE_DATA_SERVER_PATH")
	)
	data_path := filepath.Join(file_server, file_data_server, file.User_id.String(), file.File_id.String(), file.File_path)

	ctx.FileAttachment(data_path, file.File_path)

}

// 파일을 삭제하는 로직
func AuthRemoveFileService(fileIdDto *dtos.FileIdDto) (int, error) {

	var (
		db          *gorm.DB = settings.DB
		file        models.File
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 파일 찾는 로직
	if errorStatus, err = AuthRemoveFileFindDatabaseFunc(c, db, &file, fileIdDto); err != nil {
		return errorStatus, err
	}

	// 파일을 옮기고 삭제 테이블에 추가함
	if err = AuthRemoveFileMovePathFunc(c, db, &file); err != nil {
		return http.StatusInternalServerError, err
	}

	// 파일을 삭제 함
	if err = AuthRemoveFileDeleteFunc(c, db, &file); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 파일을 찾는 로직
func AuthRemoveFileFindDatabaseFunc(c context.Context, db *gorm.DB, files *models.File, fileIdDto *dtos.FileIdDto) (int, error) {

	// 데이터 찾기
	if result := db.WithContext(c).Where("file_id = ?", fileIdDto.File_id).First(files); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusBadRequest, errors.New("클라이언트에서 보낸 파일 아이디가 문제가 있습니다")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 파일아이디에 해당하는 데이터를 찾는데 오류가 발생했습니다")
		}
	}

	return 0, nil
}

// 파일 옮기고 삭제하기
func AuthRemoveFileMovePathFunc(c context.Context, db *gorm.DB, files *models.File) error {

	var (
		file_server             string = os.Getenv("FILE_SERVER_PATH")
		file_data_server        string = os.Getenv("FILE_DATA_SERVER_PATH")
		file_data_remove_server string = os.Getenv("FILE_REMOVE_DATA_SERVER_PATH")
	)
	origin_file_path := filepath.Join(file_server, file_data_server, files.User_id.String(), files.File_id.String(), files.File_path)
	new_file_path := filepath.Join(file_server, file_data_remove_server, files.User_id.String(), files.File_id.String(), files.File_path)

	if err := os.Rename(origin_file_path, new_file_path); err != nil {
		fmt.Println("시스템 오류: ", err.Error())
		return errors.New("파일을 옮기는데 오류가 발생했습니다")
	}

	// 삭제 데이터 베이스에 추가
	var (
		remove_file models.DeleteFile
	)
	remove_file.User_id = files.User_id
	remove_file.File_id = files.File_id
	remove_file.File_title = files.File_title
	remove_file.File_comment = files.File_commnet
	remove_file.File_path = files.File_path
	if result := db.WithContext(c).Create(&remove_file); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("삭제 데이터 베이스에 데이터를 추가하는데 오류가 발생했습니다")
	}

	return nil
}

// 본 데이터 베이스에서 파일을 삭제함
func AuthRemoveFileDeleteFunc(c context.Context, db *gorm.DB, file *models.File) error {

	if result := db.WithContext(c).Unscoped().Delete(file); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("오리지널 파일을 삭제하는데 오류가 발생했습니다")
	}

	return nil
}

// 유저를 로그아웃 해줌
func AuthLogoutService(ctx *gin.Context, payload *dtos.Payload) (int, error) {
	var (
		db          *gorm.DB = settings.DB
		user        models.User
		errorStatus int
		err         error
	)
	c, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// 유저 정보를 찾음
	if errorStatus, err = AuthLogoutFindUserFunc(c, db, payload, &user); err != nil {
		return errorStatus, err
	}

	// 정보 업데이트
	if err = AuthLogoutRemoveJwtTokenFunc(c, ctx, db, &user); err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

// 유저 정보를 찾음
func AuthLogoutFindUserFunc(c context.Context, db *gorm.DB, payload *dtos.Payload, user *models.User) (int, error) {

	// 데이터 찾음
	result := db.WithContext(c).Where("user_id = ?", payload.User_id).First(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return http.StatusUnauthorized, errors.New("유저의 정보를 찾을수 없습니다 다시 로그인해주세요")
		} else {
			fmt.Println("시스템 오류: ", result.Error.Error())
			return http.StatusInternalServerError, errors.New("데이터 베이스에서 유저 정보를 찾는데 오류가 발생했습니다")
		}
	}

	return 0, nil
}

// 유저의 정보를 업로드하고 정리함
func AuthLogoutRemoveJwtTokenFunc(c context.Context, ctx *gin.Context, db *gorm.DB, user *models.User) error {

	// 유저 정보부터 삭제
	user.Access_token = nil
	user.Refresh_token = nil
	user.Computer_number = nil

	// 유저 정보 업데이트
	if result := db.WithContext(c).Save(user); result.Error != nil {
		fmt.Println("시스템 오류: ", result.Error.Error())
		return errors.New("유저의 정보를 업로드 하는데 오류가 발생했습니다")
	}

	// 쿠키 보내기
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", "", 24*60*60, "", "", false, true)
	return nil
}
