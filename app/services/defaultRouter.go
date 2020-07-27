package services

import (
	"gitlab.eduwill.net/dev_team/jarvis-who/app/base/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/services/health"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/services/sample"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/services/auth"
)

// SetupRoutes connects the HTTP API endpoints to the handlers
func DefaultRouter() *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(handler.DefaultHandler()) // 핸들러

	r.Use(cors.Default())
	
	health.HealthRouter(r)
	sample.SampleRouter(r)
	auth.AuthRouter(r)

	return r
}
