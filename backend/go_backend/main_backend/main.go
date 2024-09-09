package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kimdwan/private_file/settings"
	"github.com/kimdwan/private_file/src/middlewares"
	"github.com/kimdwan/private_file/src/pkgs/routes"
)

func init() {
	settings.LoadDotenv()
	settings.LoadDatabase()
	settings.MigrateDatabase()
}

func main() {
	var (
		port string = os.Getenv("GO_PORT")
	)
	if port == "" {
		panic("환경변수에 port 번호를 입력하지 않았습니다.")
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(middlewares.CorsMiddleware())

	routes.TestRouter(router)
	routes.UserRouter(router)

	router.Run(port)
}
