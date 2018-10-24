package middlewares

import (
	//log "github.com/Sirupsen/logrus"
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
