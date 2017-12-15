package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
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
			ExpiresAt: time.Now().Add(time.Hour * Cfg.AppSettings.Expiration).Unix(),
			Issuer:    Cfg.AppSettings.Protocol + "://" + Cfg.AppSettings.Host + ":" + Cfg.AppSettings.Port,
			Audience:  Cfg.ConsumerSettings.Protocol + "://" + Cfg.ConsumerSettings.Host + ":" + Cfg.ConsumerSettings.Port,
		},
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with our secret */
	tokenString, err := token.SignedString([]byte(Cfg.AppSettings.Token))

	if err != nil {
		return "", models.NewAppError("RetrieveToken", "utils.token.retrieve_token", err.Error(), 500)
	}

	return tokenString, nil
}
