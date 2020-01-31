package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"twreporter.org/go-api/globals"
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
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// IDToken
type IDTokenJWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.StandardClaims
}

type AccessTokenJWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func (idc IDTokenJWTClaims) Valid() error {
	const verifyRequired = true
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

	if !idc.VerifyAudience(globals.Conf.App.JwtAudience, verifyRequired) {
		errMsg := "Invalid audience"
		err = *(jwt.NewValidationError(errMsg, jwt.ValidationErrorClaimsInvalid))
		return err
	}

	if !idc.VerifyIssuer(globals.Conf.App.JwtIssuer, verifyRequired) {
		errMsg := "Invalid issuer"
		err = *(jwt.NewValidationError(errMsg, jwt.ValidationErrorClaimsInvalid))
		return err
	}

	return nil
}

func RetrieveV2IDToken(userID uint, email, firstName, lastName string, expiration int) (string, error) {
	claims := IDTokenJWTClaims{
		userID,
		email,
		firstName,
		lastName,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
			Issuer:    globals.Conf.App.JwtIssuer,
			Audience:  globals.Conf.App.JwtAudience,
			Subject:   IDTokenSubject,
		},
	}
	return genToken(claims, globals.Conf.App.JwtSecret)
}

func RetrieveV2AccessToken(userID uint, email string, expiration int) (string, error) {
	claims := AccessTokenJWTClaims{
		userID,
		email,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
			Issuer:    globals.Conf.App.JwtIssuer,
			Audience:  globals.Conf.App.JwtAudience,
			Subject:   AccessTokenSubject,
		},
	}
	return genToken(claims, globals.Conf.App.JwtSecret)
}

// RetrieveMailServiceAccessToken generate JWT for mail service validation
func RetrieveMailServiceAccessToken(expiration int) (string, error) {
	var secret = globals.MailServiceJWTPrefix + globals.Conf.App.JwtSecret
	var claims = jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
		Issuer:    globals.Conf.App.JwtIssuer,
		Audience:  globals.Conf.App.JwtAudience,
		Subject:   AccessTokenSubject,
	}

	return genToken(claims, secret)
}

// genToken - generate jwt token according to user's info
func genToken(claims jwt.Claims, secret string) (string, error) {
	const errorWhere = "RetrieveToken"
	var err error
	var token *jwt.Token
	var tokenString string

	// create the token
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	tokenString, err = token.SignedString([]byte(secret))

	if err != nil {
		return "", errors.Wrap(err, "internal server error: fail to generate token")
	}
	return tokenString, nil
}
