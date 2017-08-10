package middlewares

import (
	"github.com/gin-gonic/gin"
)

// SetCacheControl ...
func SetCacheControl(cc string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", cc)
	}
}
