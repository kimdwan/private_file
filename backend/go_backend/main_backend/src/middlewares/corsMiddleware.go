package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var (
			origin        string   = ctx.GetHeader("Origin")
			allowed_hosts []string = strings.Split(os.Getenv("GO_ALLOWED_HOSTS"), ",")
			isAllowed     bool     = false
		)

		// URL을 파싱하고 정리한다.
		originUrl, err := url.Parse(origin)
		if err != nil {
			fmt.Println("시스템 오류: ", err.Error())
			ctx.AbortWithStatus(http.StatusBadGateway)
			return
		}
		originPathName := originUrl.Hostname()

		// 확인한다
		for _, allowed_host := range allowed_hosts {
			if originPathName == allowed_host {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			ctx.Writer.Header().Add("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
			ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type, X-Requested-With, User-Computer-Number")
			ctx.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET")
		}

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
