package middlewares

import (
	"github.com/gin-gonic/gin"
)

// SetCacheControl ...
func SetCacheControl(cc string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		// Append Etag or Last-Modified on response header for validation
		c.Writer.Header().Set("Cache-Control", cc)
	}
}
