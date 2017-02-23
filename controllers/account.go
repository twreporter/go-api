package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"net/mail"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

var (
	cfg  = configs.GetConfig()
	salt = []byte(cfg.ENCRYPT.SALT)
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

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// generateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// generateEncryptedPassword returns encryptedly
// securely generated string.
func generateEncryptedPassword(password []byte) (string, error) {
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	return fmt.Sprintf("%x", key), err
}

// getEmailAndPasswordFromPOSTBody get email and password value from the POST body
func getEmailAndPasswordFromPOSTBody(c *gin.Context) (string, []byte, error) {
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

// Authenticate authenticate a reporter account
func (ac AccountController) Authenticate(c *gin.Context) {

	email, password, err := getEmailAndPasswordFromPOSTBody(c)

	if err != nil {
		c.JSON(400, gin.H{"status": err})
		return
	}

	encryptedPassword, err := generateEncryptedPassword(password)

	if err != nil {
		utils.LogError(err, "Encrypting password occurs error")
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
	} else {

		_account, err := ac.Storage.GetReporterAccountData(email)

		if err != nil {
			c.JSON(401, gin.H{"status": "Account is not existed"})
			return
		}

		if !_account.Active {
			c.JSON(401, gin.H{"status": "Account is not activated yet"})
			return
		}

		if encryptedPassword == _account.Password {
			c.JSON(200, gin.H{"status": "You are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "Invalid password"})
		}
	}
}

// Signup create/update a reporter account
func (ac AccountController) Signup(c *gin.Context) {

	email, password, err := getEmailAndPasswordFromPOSTBody(c)

	if err != nil {
		c.JSON(400, gin.H{"status": err.Error()})
		return
	}

	// Check if mail address is not malform
	_, err = mail.ParseAddress(email)

	if err != nil {
		c.JSON(400, gin.H{"status": err.Error()})
		return
	}

	encryptedPassword, err := generateEncryptedPassword(password)
	if err != nil {
		utils.LogError(err, "Encrypting password occurs error")
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	log.WithFields(log.Fields{
		"account":  email,
		"password": encryptedPassword,
	}).Info("User account and password")

	ra, err := ac.Storage.GetReporterAccountData(email)

	if err == nil {
		// account is signuped and activated
		if ra.Active {
			c.JSON(409, gin.H{"status": "Account already signed up"})
			return
		}

		// account is not activated,
		// we think the signup request as a request for changing password
		_, err = ac.Storage.UpdateReporterAccountPassword(ra, encryptedPassword)
		if err != nil {
			utils.LogError(err, "Updating account password occurs error")
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		} else {
			c.JSON(200, gin.H{"status": "Password reset"})
		}
		return
	}

	// generate active token
	activeToken, err := generateRandomString(8)

	if err != nil {
		utils.LogError(err, "Generating active token occurs error")
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	// create records both in reporter_accounts and users table
	_, err = ac.Storage.InsertUserByReporterAccount(models.ReporterAccount{
		Account:       email,
		Password:      encryptedPassword,
		Active:        false,
		ActivateToken: activeToken,
	})

	if err != nil {
		utils.LogError(err, "Inserting record into users table occurs error")
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	c.JSON(201, gin.H{"status": "Sign up successfully", "email": email})
}

// Activate make existed reporter account active
func (ac AccountController) Activate(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")

	ra, err := ac.Storage.GetReporterAccountData(email)

	if err != nil {
		c.JSON(401, gin.H{"status": "Account not found"})
	}

	if ra.ActivateToken == token {
		_, err = ac.Storage.UpdateReporterAccountActive(ra, true)

		if err != nil {
			utils.LogError(err, "Updating reporter account active occurs error")
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		}

		c.JSON(200, gin.H{"status": "Account is activated", "Account": email})
		return
	}

	c.JSON(401, gin.H{"status": "Invalid token"})
}
