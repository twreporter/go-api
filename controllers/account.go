package controllers

import (
	"net/mail"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
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

// parseAccountPostBody get email and password value from the POST body
func (mc *MembershipController) parseAccountPostBody(c *gin.Context) (string, []byte, error) {
	var err error
	var form LoginForm
	var json LoginJSON

	contentType := c.ContentType()

	if contentType == "application/json" {
		err = c.Bind(&json)
		if err != nil {
			return "", nil, err
		}
		return json.Email, []byte(json.Password), nil
	} else if contentType == "application/x-www-form-urlencoded" {
		err = c.Bind(&form)
		if err != nil {
			return "", nil, err
		}
		return form.Email, []byte(form.Password), nil
	}

	return "", nil, models.NewAppError("parseAccountPostBody", "controllers.account.parse_post_body", "POST body is neither JSON nor x-www-form-urlencoded", 400)
}

// Authenticate authenticate a reporter account
func (mc *MembershipController) Authenticate(c *gin.Context) {
	var account *models.ReporterAccount
	var err error
	var email string
	var password []byte
	var encryptedPassword string

	email, password, err = mc.parseAccountPostBody(c)

	if err != nil {
		log.Error(err.Error())
		c.JSON(400, gin.H{"status": err})
		return
	}

	encryptedPassword, err = utils.GenerateEncryptedPassword(password)

	if err != nil {
		log.Error("controllers.account.authenticate.generate_encrypted_password_error \n", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
	} else {

		account, err = mc.Storage.GetReporterAccountData(email)

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
			user, err := mc.Storage.GetUserDataByReporterAccount(account)
			if err != nil {
				log.Error("controllers.account.active.get_user_data_by_reporter_account_error \n", err.Error())
				c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
			}

			jwt, err := utils.RetrieveToken(user.ID, user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
			if err != nil {
				c.JSON(500, gin.H{"status": "Internal Server Error", "error": err.Error()})
				return
			}

			c.JSON(200, gin.H{"status": "You are logged in",
				"id": user.ID, "privilege": user.Privilege, "firstname": user.FirstName.String,
				"lastname": user.LastName.String, "email": user.Email.String, "jwt": jwt})
		} else {
			c.JSON(401, gin.H{"status": "Invalid password"})
		}
	}
}

// Signup create/update a reporter account
func (mc *MembershipController) Signup(c *gin.Context, mailSender utils.EmailSender) {
	var activeToken string
	var email string
	var err error
	var encryptedPassword string
	var password []byte
	var ra *models.ReporterAccount

	email, password, err = mc.parseAccountPostBody(c)
	log.Info("email ", email, " password ", password)

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

	encryptedPassword, err = utils.GenerateEncryptedPassword(password)
	if err != nil {
		log.Error("controllers.account.sign_up.generate_encrypted_password_error \n", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error"})
		return
	}

	log.WithFields(log.Fields{
		"account":  email,
		"password": encryptedPassword,
	}).Info("User account and password")

	ra, err = mc.Storage.GetReporterAccountData(email)

	if err == nil {
		// account is signuped and activated
		if ra.Active {
			c.JSON(409, gin.H{"status": "Account already signed up"})
			return
		}

		// account is not activated,
		// we think the signup request as a request for changing password
		if _, err = mc.Storage.UpdateReporterAccountPassword(ra, encryptedPassword); err != nil {
			log.Error("controllers.account.sign_up.update_db_error \n", err.Error())
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
			return
		}

		go func() {
			// re-send the activation email
			if err1 := mailSender.Send(email, activateMailSubject, utils.GenerateActivateMailBody(email, ra.ActivateToken)); err1 != nil {
				log.Error("controllers.account.sign_up.send_mail \n", err1.Error())
			}
		}()

		c.JSON(200, gin.H{"status": "Password reset and activation email resent"})
		return
	}

	// generate active token
	activeToken, err = utils.GenerateRandomString(8)

	if err != nil {
		log.Error("controllers.account.sign_up.generate_active_token_error \n", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	// create records both in reporter_accounts and users table
	_, err = mc.Storage.InsertUserByReporterAccount(models.ReporterAccount{
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
		if err1 := mailSender.Send(email, activateMailSubject, utils.GenerateActivateMailBody(email, activeToken)); err1 != nil {
			log.Error("controllers.account.sign_up.send_mail \n", err1.Error())
		}
	}()

	c.JSON(201, gin.H{"status": "Sign up successfully", "email": email})
}

// Activate make existed reporter account active
func (mc *MembershipController) Activate(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")

	ra, err := mc.Storage.GetReporterAccountData(email)

	if err != nil {
		c.JSON(401, gin.H{"status": "Account not found"})
	}

	if ra.ActivateToken == token {
		_, err = mc.Storage.UpdateReporterAccountActive(ra, true)

		if err != nil {
			log.Error("controllers.account.active.update_db_error \n", err.Error())
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		}

		user, err := mc.Storage.GetUserDataByReporterAccount(ra)
		if err != nil {
			log.Error("controllers.account.active.get_user_data_by_reporter_account_error \n", err.Error())
			c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		}

		jwt, err := utils.RetrieveToken(user.ID, user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
		if err != nil {
			c.JSON(500, gin.H{"status": "Internal Server Error", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "Account is activated",
			"id": user.ID, "privilege": user.Privilege, "firstname": user.FirstName.String,
			"lastname": user.LastName.String, "email": user.Email.String, "jwt": jwt})
		return
	}

	c.JSON(401, gin.H{"status": "Invalid token"})
}
