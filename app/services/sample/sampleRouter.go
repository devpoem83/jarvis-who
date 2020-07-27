package sample

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes connects the HTTP API endpoints to the handlers
func SampleRouter(r *gin.Engine) {
	sample := r.Group("samples")
	{
		sample.GET("/", Samples)
		sample.GET("/:id", ReadSample)
		sample.POST("/", CreateSample)
		sample.PUT("/:id", UpdateSample)
		sample.DELETE("/:id", DeleteSample)
	}
}
