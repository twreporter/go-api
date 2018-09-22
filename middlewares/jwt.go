package middlewares

import (
	"fmt"
	"net/http"

	"twreporter.org/go-api/globals"

	//log "github.com/Sirupsen/logrus"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(globals.Conf.App.JwtSecret), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

// CheckJWT checks the jwt token in the Authorization header is valid or not
func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
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
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func ValidateIDToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Request.Context().Value("user")

		claims := user.(*jwt.Token).Claims

		if err := claims.Valid(); nil != err {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"data": gin.H{
					"id_token": err.Error(),
				},
			})
		}
	}
}

// SetEmailClaim get email value from jwt, and set it into gin.Context
func SetEmailClaim() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Request.Context().Value("user")
		emailClaim := user.(*jwt.Token).Claims.(jwt.MapClaims)["email"]

		c.Set("emailClaim", fmt.Sprint(emailClaim))

		// The reason why we don't validate the email value in POST body with jwt email claim HERE
		// is because c.Bind(&login) only can be executed once.
		// If we get POST body here, then in the controllers, such as controllers/account.go,
		// we cannot c.Bind(&login) to get the POST body.
		// Hence, here only set emailClaim in the gin.Context
	}
}

// ValidateAdminUsers ...
func ValidateAdminUsers() gin.HandlerFunc {
	var whiteList = []string{"nickhsine@twreporter.org", "han@twreporter.org", "yucj@twreporter.org", "developer@twreporter.org"}
	return func(c *gin.Context) {
		user := c.Request.Context().Value("user")
		userIDClaim := user.(*jwt.Token).Claims.(jwt.MapClaims)["email"]

		for _, v := range whiteList {
			if v == fmt.Sprint(userIDClaim) {
				return
			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
