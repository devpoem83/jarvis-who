package config

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func StaticInit(r *gin.Engine) {
	// 정적리소스 설정
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	// Favicon 설정
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.StaticFile("/logo.png", "./public/logo.png")
}
