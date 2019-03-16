package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/utils"

	// log "github.com/Sirupsen/logrus"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const authUserProperty = "app-auth-jwt"

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(globals.Conf.App.JwtSecret), nil
	},
	UserProperty:  authUserProperty,
	SigningMethod: jwt.SigningMethodHS256,
	ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
		var res = map[string]interface{}{
			"status": "fail",
			"data": map[string]interface{}{
				"req.Headers.Authorization": err,
			},
		}
		var resByte, _ = json.Marshal(res)
		http.Error(w, string(resByte), http.StatusUnauthorized)
	},
})

// ValidateAuthorization checks the jwt token in the Authorization header is valid or not
func ValidateAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		const verifyRequired = true
		var err error
		var userProperty interface{}
		var claims jwt.MapClaims

		if err = jwtMiddleware.CheckJWT(c.Writer, c.Request); err != nil {
			c.Abort()
			return
		}

		userProperty = c.Request.Context().Value(authUserProperty)
		claims = userProperty.(*jwt.Token).Claims.(jwt.MapClaims)
		if !claims.VerifyAudience(globals.Conf.App.JwtAudience, verifyRequired) ||
			!claims.VerifyIssuer(globals.Conf.App.JwtIssuer, verifyRequired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"data": gin.H{
					"req.Cookies.id_token": "aud or issuer claim is invalid",
				},
			})
			return
		}

		var newRequest *http.Request

		// Set user_id with key "auth-user-id" in context to avoid hierarchy access
		newRequest = c.Request.WithContext(context.WithValue(c.Request.Context(), globals.AuthUserIDProperty, claims["user_id"]))
		*c.Request = *newRequest
	}
}

// ValidateUserID checks claim userID in the jwt with :userID param in the request url.
// if the two values are not the same, return the 401 response
func ValidateUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userID     string
			authUserID interface{}
		)

		authUserID = c.Request.Context().Value(globals.AuthUserIDProperty)
		userID = c.Param("userID")
		if userID != fmt.Sprint(authUserID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status": "fail",
				"data": gin.H{
					"req.Headers.Authorization": "the request is not permitted to reach the resource",
				},
			})
		}
	}
}

func ValidateUserIDInReqBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body = struct {
				UserID uint64 `json:"user_id" form:"user_id" binding:"required"`
			}{}
			err        error
			authUserID interface{}
		)

		// gin.Context.Bind does not support to bind `JSON` body multiple times
		// the alternative is to use gin.Context.ShouldBindBodyWith function to bind
		if err = c.ShouldBindBodyWith(&body, binding.JSON); err == nil {
			// omit intentionally
		} else if err = c.Bind(&body); err != nil {
			// bind other format rather than JSON
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{
				"req.Body.user_id": err.Error(),
			}})
			return
		}

		authUserID = c.Request.Context().Value(globals.AuthUserIDProperty)

		if fmt.Sprint(body.UserID) != fmt.Sprint(authUserID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "data": gin.H{
				"req.Headers.Authorization": "the request is not permitted to reach the resource",
			}})
			return
		}
	}
}

// ValidateAuthentication validates `req.Cookies.id_token`
// if id_token, which is a JWT, is invalid, and then return 401 status code
func ValidateAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var err error
		var token *jwt.Token

		defer func() {
			if r := recover(); r != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status": "fail",
					"data": gin.H{
						"req.Headers.Cookies.id_token": err.Error(),
					},
				})
				return
			}
		}()

		if tokenString, err = c.Cookie("id_token"); err != nil {
			panic(err)
		}

		if token, err = jwt.ParseWithClaims(tokenString, &utils.IDTokenJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(globals.Conf.App.JwtSecret), nil
		}); err != nil {
			panic(err)
		}

		if !token.Valid {
			err = errors.New("id_token is invalid")
			panic(err)
		}
	}
}
