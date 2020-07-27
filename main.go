package main // import "gitlab.eduwill.net/dev_team/jarvis-who"

import (
	"strconv"
	
	"gitlab.eduwill.net/dev_team/jarvis-who/app/services"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/base/config"
)

func main() {

	// 프로파일 초기화
	config.ProfileInit()

	// Logger 초기화
	common.LoggerInit(config.GetProfile().LoggerLevel)

	// 캐시 초기화
	config.CacheInit()

	// DB 초기화
	config.DbInit()

	// 라우트 설정
	r := services.DefaultRouter()

	// 정적리소스 설정
	config.StaticInit(r)

	r.Run(":" + strconv.Itoa(config.GetProfile().ServerPort))
}
