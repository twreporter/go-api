package utils

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
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
			ExpiresAt: time.Now().Add(time.Hour * viper.GetDuration("appsettings.expiration")).Unix(),
			Issuer:    viper.GetString("appsettings.protocol") + "://" + viper.GetString("appsettings.host") + ":" + viper.GetString("appsettings.port"),
			Audience:  viper.GetString("consumersettings.protocol") + "://" + viper.GetString("consumersettings.host") + ":" + viper.GetString("consumersettings.port"),
		},
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	tokenString, err := token.SignedString([]byte(viper.GetString("appsettings.token")))

	if err != nil {
		return "", models.NewAppError("RetrieveToken", "internal server error: fail to generate token", err.Error(), http.StatusInternalServerError)
	}

	return tokenString, nil
}
