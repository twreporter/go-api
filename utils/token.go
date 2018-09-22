package utils

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	//log "github.com/Sirupsen/logrus"
)

type AuthTokenType int

const (
	AuthV1Token AuthTokenType = iota + 1

	AuthV2IDToken
	AuthV2AccessToken
)

const (
	IDTokenSubject     = "ID_TOKEN"
	AccessTokenSubject = "ACCESS_TOKEN"
)

// ReporterJWTClaims JWT claims we used
type ReporterJWTClaims struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// IDToken
type IDTokenJWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type AccessTokenJWTClaims struct {
	jwt.StandardClaims
}

func (idc IDTokenJWTClaims) Valid() error {
	var err error

	// Validate expiration date
	if err = idc.StandardClaims.Valid(); nil != err {
		return err
	}

	if IDTokenSubject != idc.StandardClaims.Subject {
		errMsg := "Invalid subject"
		err = *(jwt.NewValidationError(errMsg, jwt.ValidationErrorClaimsInvalid))
		return err
	}

	return nil
}

func RetrieveV1Token(userID uint, email string) (string, error) {
	return genToken(AuthV1Token, userID, email, globals.Conf.App.JwtExpiration)
}

func RetrieveV2Token(tokenType AuthTokenType, userID uint, email string, expiration int) (string, error) {
	return genToken(tokenType, userID, email, expiration)
}

// genToken - generate jwt token according to user's info
func genToken(tokenType AuthTokenType, userID uint, email string, expiration int) (string, error) {
	var errorWhere = "RetrieveToken"
	var claims jwt.Claims

	// Create the Claims
	switch tokenType {
	case AuthV1Token:
		claims = ReporterJWTClaims{
			userID,
			email,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
				Issuer:    globals.Conf.App.JwtIssuer,
				Audience:  globals.Conf.App.JwtAudience,
			},
		}

	case AuthV2IDToken:
		claims = IDTokenJWTClaims{
			userID,
			email,
			jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
				Issuer:    globals.Conf.App.JwtIssuer,
				Audience:  globals.Conf.App.JwtAudience,
				Subject:   IDTokenSubject,
			},
		}
	case AuthV2AccessToken:
		claims = AccessTokenJWTClaims{
			jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
				Issuer:    globals.Conf.App.JwtIssuer,
				Audience:  globals.Conf.App.JwtAudience,
				Subject:   AccessTokenSubject,
			},
		}
	default:
		return "", models.NewAppError(errorWhere, "Invalid token type", "Invalid token type", http.StatusInternalServerError)
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	tokenString, err := token.SignedString([]byte(globals.Conf.App.JwtSecret))

	if err != nil {
		return "", models.NewAppError("RetrieveToken", "internal server error: fail to generate token", err.Error(), http.StatusInternalServerError)
	}

	return tokenString, nil
}
