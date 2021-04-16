package middlewares

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/globals"
)

const jwtUserPropertyForMailService = "mail-service-jwt"

type JWTMiddleware interface {
	ValidateAuthorization() gin.HandlerFunc
}

type mailServiceMiddleware struct {
	JWTMiddleware *jwtmiddleware.JWTMiddleware
}

func (m mailServiceMiddleware) ValidateAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := m.JWTMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func GetMailServiceMiddleware() JWTMiddleware {
	return mailServiceMiddleware{
		JWTMiddleware: jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(globals.MailServiceJWTPrefix + globals.Conf.App.JwtSecret), nil
			},
			UserProperty:  jwtUserPropertyForMailService,
			SigningMethod: jwt.SigningMethodHS256,
		}),
	}
}
