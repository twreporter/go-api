package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// RetrieveToken - generate jwt token according to user's info
func RetrieveToken(userID uint, privilege int, firstName string, lastName string, email string) (string, error) {
	log.WithFields(log.Fields{
		"ID":        userID,
		"privilege": privilege,
		"firstName": firstName,
		"lastName":  lastName,
		"email":     email,
	}).Info("RetrieveToken")

	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// set token claims
	claims["userID"] = userID
	claims["privilege"] = privilege
	claims["firstName"] = firstName
	claims["lastName"] = lastName
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * Cfg.AppSettings.Expiration).Unix()

	/* Sign the token with our secret */
	tokenString, err := token.SignedString([]byte(Cfg.AppSettings.Token))

	if err != nil {
		return "", models.NewAppError("RetrieveToken", "utils.token.retrieve_token", err.Error(), 500)
	}

	return tokenString, nil
}
