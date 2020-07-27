package auth

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes connects the HTTP API endpoints to the handlers
func AuthRouter(r *gin.Engine) {
	auth := r.Group("auth")
	{
		auth.POST("/signin", Signin)
		auth.POST("/welcome", Welcome)
		auth.POST("/refresh", Refresh)
	}
}
