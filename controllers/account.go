package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"

	log "github.com/Sirupsen/logrus"
)

// LoginForm is to be binded from form values
type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// LoginJSON is to be binded from json values
type LoginJSON struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AccountController ...
type AccountController struct {
	Storage *storage.UserStorage
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateEncryptedPassword returns encryptedly
// securely generated string.
func GenerateEncryptedPassword(password []byte) (string, error) {
	salt := []byte("@#$%")
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	return fmt.Sprintf("%x", key), err
}

// LogError print log message and log error as error log level
func LogError(err error, message string) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error(message)
}

// GetEmailAndPasswordFromPOSTBody get email and password value from the POST body
func GetEmailAndPasswordFromPOSTBody(c *gin.Context) (string, []byte, error) {
	var email string
	var form LoginForm
	var json LoginJSON
	var password []byte

	// Request Header
	// Content-Type: x-www-form-urlencoded
	formErr := c.Bind(&form)

	// Content-Type: application-json
	jsonErr := c.Bind(&json)

	if formErr == nil || jsonErr == nil {
		if formErr == nil {
			email = form.Email
			password = []byte(form.Password)
		} else {
			email = json.Email
			password = []byte(json.Password)
		}

		return email, password, nil
	}

	return "", nil, errors.New("Bad syntax in POST body")
}

// Authenticate Reporter account
func (ac AccountController) Authenticate(c *gin.Context) {

	email, password, err := GetEmailAndPasswordFromPOSTBody(c)

	if err != nil {
		c.JSON(400, gin.H{"status": err})
		return
	}

	encryptedPassword, err := GenerateEncryptedPassword(password)

	if err != nil {
		LogError(err, "Encrypting password occurs error")
		c.JSON(500, gin.H{"status": "internal server error"})
	} else {

		_account, err := ac.Storage.GetReporterAccountData(email)

		if err != nil {
			c.JSON(401, gin.H{"status": "Account is not existed"})
			return
		}

		if encryptedPassword == _account.Password {
			c.JSON(200, gin.H{"status": "You are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "Invalid password"})
		}
	}
}

// Signup Reporter account
func (ac AccountController) Signup(c *gin.Context) {

	email, password, err := GetEmailAndPasswordFromPOSTBody(c)

	if err != nil {
		c.JSON(400, gin.H{"status": err})
		return
	}

	encryptedPassword, err := GenerateEncryptedPassword(password)
	if err != nil {
		LogError(err, "Encrypting password occurs error")
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	log.WithFields(log.Fields{
		"account":  email,
		"password": encryptedPassword,
	}).Info("User account and password")

	_, err = ac.Storage.GetReporterAccountData(email)

	if err == nil {
		c.JSON(409, gin.H{"status": "Email already signed up"})
		return
	}

	// generate active token
	activeToken, err := GenerateRandomString(8)

	if err != nil {
		LogError(err, "Generating active token occurs error")
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	// create records both in reporter_accounts and users table
	_, err = ac.Storage.InsertUserByReporterAccount(models.ReporterAccount{
		Email:         email,
		Password:      encryptedPassword,
		Active:        false,
		ActivateToken: activeToken,
	})

	if err != nil {
		LogError(err, "Inserting record into users table occurs error")
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	c.JSON(201, gin.H{"status": "Sign up successfully", "email": email})
}
