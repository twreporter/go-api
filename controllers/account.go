package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"net/mail"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

const (
	activateMailSubject = "啟動報導者帳號"
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
	Storage storage.UserStorage
}

// generateActivateMailBody generate the html a tag which can link to /active enpoint to activate the account
func generateActivateMailBody(mailAddress, activeToken string) string {
	href := fmt.Sprintf("%s/%s/activate?email=%s&token=%s", utils.Cfg.AppSettings.Path, utils.Cfg.AppSettings.Version, mailAddress, activeToken)

	// TBD make the activate mail more beautiful and informative
	return fmt.Sprintf("<a href=\"%s\" target=\"_blank\">Activate Your Account</a>", href)
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
	salt := []byte(utils.Cfg.EncryptSettings.Salt)
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

	return "", nil, models.NewAppError("getEmailAndPasswordFromPOSTBody", "controllers.account.parse_post_body", "POST body is neither JSON nor x-www-form-urlencoded", 500)
}

// Authenticate authenticate a reporter account
func (ac AccountController) Authenticate(c *gin.Context) {
	var account *models.ReporterAccount
	var err error
	var email string
	var password []byte
	var encryptedPassword string

	email, password, err = getEmailAndPasswordFromPOSTBody(c)

	if err != nil {
		log.Error(err.Error())
		c.JSON(400, gin.H{"status": err})
		return
	}

	encryptedPassword, err = generateEncryptedPassword(password)

	if err != nil {
		log.Error("controllers.account.authenticate.generate_encrypted_password_error \n", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
	} else {

		account, err = ac.Storage.GetReporterAccountData(email)

		if err != nil {
			c.JSON(401, gin.H{"status": "Account is not existed"})
			return
		}

		if !account.Active {
			c.JSON(401, gin.H{"status": "Account is not activated yet"})
			return
		}

		log.WithFields(log.Fields{
			"salt":              utils.Cfg.EncryptSettings.Salt,
			"encryptedPassword": encryptedPassword,
			"password":          account.Password,
		}).Info("Password Info ")

		if encryptedPassword == account.Password {
			user, err := ac.Storage.GetUserDataByReporterAccount(account)
			if err != nil {
				log.Error("controllers.account.active.get_user_data_by_reporter_account_error \n", err.Error())
				c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
			}

			jwt := utils.RetrieveToken(user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
			c.JSON(200, gin.H{"status": "You are logged in", "jwt": jwt})
		} else {
			c.JSON(401, gin.H{"status": "Invalid password"})
		}
	}
}

// Signup create/update a reporter account
func (ac AccountController) Signup(c *gin.Context, mailSender utils.EmailSender) {
	var activeToken string
	var email string
	var err error
	var encryptedPassword string
	var password []byte
	var ra *models.ReporterAccount

	email, password, err = getEmailAndPasswordFromPOSTBody(c)
	log.Info("email", email, " password ", password)

	if err != nil {
		log.Error(err.Error())
		c.JSON(400, gin.H{"status": err.Error()})
		return
	}

	// Check if mail address is not malform
	_, err = mail.ParseAddress(email)

	if err != nil {
		c.JSON(400, gin.H{"status": err.Error()})
		return
	}

	encryptedPassword, err = generateEncryptedPassword(password)
	if err != nil {
		log.Error("controllers.account.sign_up.generate_encrypted_password_error \n", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	log.WithFields(log.Fields{
		"account":  email,
		"password": encryptedPassword,
	}).Info("User account and password")

	ra, err = ac.Storage.GetReporterAccountData(email)

	if err == nil {
		// account is signuped and activated
		if ra.Active {
			c.JSON(409, gin.H{"status": "Account already signed up"})
			return
		}

		// account is not activated,
		// we think the signup request as a request for changing password
		if _, err = ac.Storage.UpdateReporterAccountPassword(ra, encryptedPassword); err != nil {
			log.Error("controllers.account.sign_up.update_db_error \n", err.Error())
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		}

		go func() {
			// re-send the activation email
			if err1 := mailSender.Send(email, activateMailSubject, generateActivateMailBody(email, ra.ActivateToken)); err1 != nil {
				log.Error("controllers.account.sign_up.send_mail \n", err1.Error())
			}
		}()

		c.JSON(200, gin.H{"status": "Password reset and activation email resent"})
		return
	}

	// generate active token
	activeToken, err = generateRandomString(8)

	if err != nil {
		log.Error("controllers.account.sign_up.generate_active_token_error \n", err.Error())
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
		log.Error("controllers.account.sign_up.insert_db_error \n", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	go func() {
		if err1 := mailSender.Send(email, activateMailSubject, generateActivateMailBody(email, activeToken)); err1 != nil {
			log.Error("controllers.account.sign_up.send_mail \n", err1.Error())
		}
	}()

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
			log.Error("controllers.account.active.update_db_error \n", err.Error())
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		}

		user, err := ac.Storage.GetUserDataByReporterAccount(ra)
		if err != nil {
			log.Error("controllers.account.active.get_user_data_by_reporter_account_error \n", err.Error())
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		}

		jwt := utils.RetrieveToken(user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
		c.JSON(200, gin.H{"status": "Account is activated", "account": email, "jwt": jwt})
		return
	}

	c.JSON(401, gin.H{"status": "Invalid token"})
}
