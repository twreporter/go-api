package middlewares

import (
	"fmt"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"twreporter.org/go-api/utils"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Cfg.AppSettings.Token), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// CheckJWT ...
func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
		}
	}
}

// ValidateUserID this func will validate claim userID in the JWT token of the Authentication request header
// with :userID param in the request url
func ValidateUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(401)
			return
		}
		user := context.Get(c.Request, "user")
		userIDClaim := user.(*jwt.Token).Claims.(jwt.MapClaims)["userID"]
		userID := c.Param("userID")
		if userID != fmt.Sprint(userIDClaim) {
			c.AbortWithStatus(401)
		}
	}
}
