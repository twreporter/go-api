package middlewares

import (
	"fmt"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/utils"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Cfg.AppSettings.Token), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// CheckJWT checks the jwt token in the Authorization header is valid or not
func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}
	}
}

// ValidateUserID checks claim userID in the jwt with :userID param in the request url.
// if the two values are not the same, return the 401 response
func ValidateUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Request.Context().Value("user")
		userIDClaim := user.(*jwt.Token).Claims.(jwt.MapClaims)["userID"]
		userID := c.Param("userID")
		if userID != fmt.Sprint(userIDClaim) {
			c.AbortWithStatus(401)
		}
	}
}

// ValidateAdminUsers ...
func ValidateAdminUsers() gin.HandlerFunc {
	var whiteList = []string{"nickhsine@twreporter.org", "hsunpei_wang@twreporter.org", "han@twreporter.org", "yucj@twreporter.org"}
	return func(c *gin.Context) {
		user := c.Request.Context().Value("user")
		userIDClaim := user.(*jwt.Token).Claims.(jwt.MapClaims)["email"]

		for _, v := range whiteList {
			if v == fmt.Sprint(userIDClaim) {
				return
			}
		}

		c.AbortWithStatus(401)
	}
}
