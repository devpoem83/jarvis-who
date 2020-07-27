package health

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// SetupRoutes connects the HTTP API endpoints to the handlers
func HealthRouter(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "SUCCESS")
	})
}
