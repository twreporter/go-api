package utils

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	//log "github.com/Sirupsen/logrus"
)

// ReporterJWTClaims JWT claims we used
type ReporterJWTClaims struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// RetrieveToken - generate jwt token according to user's info
func RetrieveToken(userID uint, email string) (string, error) {
	// Create the Claims
	claims := ReporterJWTClaims{
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(globals.Conf.App.JwtExpiration)).Unix(),
			Issuer:    globals.Conf.App.JwtIssuer,
			Audience:  globals.Conf.App.JwtAudience,
		},
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
