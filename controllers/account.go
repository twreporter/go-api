package controllers

import (
	"golang.org/x/crypto/bcrypt"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var (
	cfg       = configs.GetConfig()
	loginPath = cfg.APP.Path + "/login"
)

// LoginForm ...
type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AccountController ...
type AccountController struct {
	Storage *storage.UserStorage
}

// Authenticate Reporter account
func (ac AccountController) Authenticate(c *gin.Context) {
	// you can bind multipart form with explicit binding declaration:
	// c.BindWith(&form, binding.Form)
	// or you can simply use autobinding with Bind method:
	var form LoginForm
	// in this case proper binding will be automatically selected

	if c.Bind(&form) == nil {
		email := form.Email
		password := []byte(form.Password)

		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		_password := string(hashedPassword)

		if err != nil {
			c.JSON(500, gin.H{"status": "internal server error"})
		} else {
			log.WithFields(log.Fields{
				"account":  email,
				"password": _password,
			}).Info("User account and encrypted password")

			_account := ac.Storage.GetReporterAccountData(email, _password)
			if _account.Email != "" {
				c.JSON(401, gin.H{"status": "unauthorized"})
			} else {
				c.JSON(200, gin.H{"status": "you are logged in"})
			}
		}
	}
}

// Signup Reporter account
func (ac AccountController) Signup(c *gin.Context) {
	// you can bind multipart form with explicit binding declaration:
	// c.BindWith(&form, binding.Form)
	// or you can simply use autobinding with Bind method:
	var form LoginForm
	// in this case proper binding will be automatically selected

	log.Info("signup...")

	if c.Bind(&form) == nil {
		log.WithFields(log.Fields{
			"account":  form.Email,
			"password": form.Password,
		}).Info("User account and password")

		email := form.Email
		password := []byte(form.Password)

		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

		_password := string(hashedPassword)

		if err != nil {
			c.JSON(500, gin.H{"status": "internal server error"})
		} else {
			log.WithFields(log.Fields{
				"account":  email,
				"password": _password,
			}).Info("User account and encrypted password")

			_account := ac.Storage.GetReporterAccountData(email, _password)

			if _account.Email == "" {
				user := ac.Storage.InsertUserByReporterAccount(models.ReporterAccount{
					Email:         email,
					Password:      _password,
					Active:        false,
					ActivateToken: "test activate token",
				})

				log.WithFields(log.Fields{
					"Email": user.Email,
					"ID":    user.ID,
				}).Info("Inserted User")

				if user.ID != 0 && user.Email.Valid {
					c.JSON(201, gin.H{"status": "Sign up successfully", "email": user.Email})
				} else {
					c.JSON(500, gin.H{"status": "internal server error"})
				}
			} else {
				c.JSON(409, gin.H{"status": "Email already signed up"})
			}
		}
	}
}
