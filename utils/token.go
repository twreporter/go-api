package utils

import (
	"time"

	"twreporter.org/go-api/configs"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
)

var (
	cfg          = configs.GetConfig()
	mySigningKey = []byte(cfg.APP.Token)
)

// RetrieveToken ...
func RetrieveToken(privilege int, firstName string, lastName string, email string) string {
	log.WithFields(log.Fields{
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
	claims["privilege"] = privilege
	claims["firstName"] = firstName
	claims["lastName"] = lastName
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * cfg.APP.Expiration).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}
