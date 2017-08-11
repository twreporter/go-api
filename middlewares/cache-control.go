package middlewares

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// SetCacheControl ...
func SetCacheControl(cc string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Avoid adding Cache-Control in the Response Header
		// while the Request Header is not containing Origin
		if c.Request.Header.Get("Origin") != "" {
			c.Writer.Header().Set("Cache-Control", cc)
		}
	}
}
