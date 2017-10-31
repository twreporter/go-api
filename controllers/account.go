package controllers

import (
	"net/http"
	"net/mail"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

const (
	activateMailSubject = "啟動報導者帳號"
)

// LoginBody is to be binded from form values
type LoginBody struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// parseAccountPostBody get email and password value from the POST body
func (mc *MembershipController) parseAccountPostBody(c *gin.Context) (LoginBody, error) {
	var err error
	var login LoginBody

	if err = c.Bind(&login); err != nil {
		log.Error("err:", err.Error())
		return login, models.NewAppError("parseAccountPostBody", "controllers.account.parse_post_body", err.Error(), http.StatusBadRequest)
	}
	log.Info("login:", login)

	return login, nil
}

// Authenticate authenticate a reporter account
func (mc *MembershipController) Authenticate(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.Authenticate"
	var account *models.ReporterAccount
	var err error
	var login LoginBody
	var email string
	var encryptedPassword string

	login, err = mc.parseAccountPostBody(c)

	if err != nil {
		appErr := err.(models.AppError)
		return appErr.StatusCode, gin.H{"status": "fail", "data": LoginBody{
			Email:    "email is required",
			Password: "password is required",
		}}, nil
	}

	email = login.Email

	encryptedPassword, err = utils.GenerateEncryptedPassword([]byte(login.Password))

	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Encrypt password occurs error", err.Error(), http.StatusInternalServerError)
	}

	account, err = mc.Storage.GetReporterAccountData(email)

	if err != nil {
		return 0, gin.H{}, err
	}

	if !account.Active {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Account is not active", "", http.StatusUnauthorized)
	}

	if encryptedPassword == account.Password {
		user, err := mc.Storage.GetUserDataByReporterAccount(account)
		if err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Getting user data from DB occurs error", err.Error(), http.StatusUnauthorized)
		}

		jwt, err := utils.RetrieveToken(user.ID, user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
		if err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Generating JWT occurs error", err.Error(), http.StatusInternalServerError)
		}

		return 200, gin.H{"status": "You are logged in",
			"id": user.ID, "privilege": user.Privilege, "firstname": user.FirstName.String,
			"lastname": user.LastName.String, "email": user.Email.String, "jwt": jwt}, nil
	}

	return 0, gin.H{}, models.NewAppError(errorWhere, "Invalid password", "", http.StatusUnauthorized)
}

// Signup - create/update a reporter account
func (mc *MembershipController) Signup(c *gin.Context, mailSender *utils.EmailContext) (int, gin.H, error) {
	const errorWhere = "MembershipController.Signup"
	var activeToken string
	var err error
	var encryptedPassword string
	var login LoginBody
	var email string
	var ra *models.ReporterAccount

	// extract email and password field in POST body
	login, err = mc.parseAccountPostBody(c)
	if err != nil {
		appErr := err.(models.AppError)
		return appErr.StatusCode, gin.H{"status": "fail", "data": LoginBody{
			Email:    "email is required",
			Password: "password is required",
		}}, nil
	}

	email = login.Email

	// Check if mail address is not malform
	_, err = mail.ParseAddress(email)
	if err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": LoginBody{
			Email: "email is malform",
		}}, nil
	}

	// generate encrypted password
	encryptedPassword, err = utils.GenerateEncryptedPassword([]byte(login.Password))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Generating encrypted password occurs error", err.Error(), http.StatusInternalServerError)
	}

	// generate active token
	activeToken, err = utils.GenerateRandomString(8)
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Generating active token occurs error", err.Error(), http.StatusInternalServerError)
	}

	// get reporter account by email from reporter_account table
	ra, err = mc.Storage.GetReporterAccountData(email)

	// account is already signup
	if err == nil {
		// account is already activated
		if ra.Active {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Account is already signup", "", http.StatusConflict)
		}

		// account is not activated,
		// we think the signup request as a request for changing password
		// update the password, active token and active expire time
		ra.Password = encryptedPassword
		ra.ActivateToken = activeToken
		ra.ActExpTime = time.Now().Add(time.Duration(24) * time.Hour)
		if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Updating DB occurs error", err.Error(), http.StatusInternalServerError)
		}

		// re-send the activation email
		err = mailSender.Send(email, activateMailSubject, utils.GenerateActivateMailBody(email, ra.ActivateToken))
		if err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Sending activation email occurs error", err.Error(), http.StatusInternalServerError)
		}

		return http.StatusCreated, gin.H{"status": "success", "data": LoginBody{
			Email: email,
		}}, nil
	}

	// create records both in reporter_accounts and users table
	_, err = mc.Storage.InsertUserByReporterAccount(models.ReporterAccount{
		Account:       email,
		Password:      encryptedPassword,
		Active:        false,
		ActivateToken: activeToken,
		// expire time is one day
		ActExpTime: time.Now().Add(time.Duration(24) * time.Hour),
	})

	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Inserting new record into DB occurs error", err.Error(), http.StatusInternalServerError)
	}

	// send activation email
	err = mailSender.Send(email, activateMailSubject, utils.GenerateActivateMailBody(email, activeToken))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Sending activation email occurs error", err.Error(), http.StatusInternalServerError)
	}

	return http.StatusCreated, gin.H{"status": "success", "data": LoginBody{
		Email: email,
	}}, nil
}

// Activate - make existed reporter account active
func (mc *MembershipController) Activate(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.Activate"

	email := c.Query("email")
	token := c.Query("token")

	// get reporter account by email from reporter_account table
	ra, err := mc.Storage.GetReporterAccountData(email)

	if err != nil {
		return 0, gin.H{}, err
	}

	// check expire time
	if ra.ActExpTime.Sub(time.Now()) < time.Duration(0) {
		return 0, gin.H{}, models.NewAppError(errorWhere, "ActivateToken is expired", "", http.StatusUnauthorized)
	}

	// check token
	if ra.ActivateToken == token {
		ra.Active = true
		err = mc.Storage.UpdateReporterAccount(ra)
		// handle internal server error - can't update the record in reporter_account table
		if err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Updating reporter_account occurs error", err.Error(), http.StatusInternalServerError)
		}

		user, err := mc.Storage.GetUserDataByReporterAccount(ra)
		// handle internal server error - can't query the record from reporter_account table
		if err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Querying reporter_account occurs error", err.Error(), http.StatusInternalServerError)
		}

		jwt, err := utils.RetrieveToken(user.ID, user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
		// handle internal server error - cannot generate JWT
		if err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Generating JWT occurs error", err.Error(), http.StatusInternalServerError)
		}

		return 200, gin.H{"status": "success",
			"id": user.ID, "privilege": user.Privilege, "firstname": user.FirstName.String,
			"lastname": user.LastName.String, "email": user.Email.String, "jwt": jwt}, nil
	}

	return 0, gin.H{}, models.NewAppError(errorWhere, "Token is invalid", "", http.StatusUnauthorized)
}

// ChangePassword - change password of the user
func (mc *MembershipController) ChangePassword(c *gin.Context) (int, gin.H, error) {
	var encryptedPassword string
	var err error
	var errorWhere = "MembershipController.ChangePassword"
	var login LoginBody
	var ra *models.ReporterAccount

	// extract email and password field in POST body
	login, err = mc.parseAccountPostBody(c)
	if err != nil {
		appErr := err.(models.AppError)
		log.Error(appErr.Error())
		return appErr.StatusCode, gin.H{"status": "fail", "data": LoginBody{
			Email:    "email is required",
			Password: "password is required",
		}}, nil
	}

	// emailClaim is set by gin.middleware in middlewares/jwt.go
	emailClaim, ok := c.Get("emailClaim")

	if !ok || emailClaim != login.Email {
		return http.StatusUnauthorized, gin.H{"status": "fail", "data": LoginBody{
			Email:    "email is not validated",
			Password: "password is required",
		}}, nil
	}

	// get reporter account by email from reporter_accounts table
	ra, err = mc.Storage.GetReporterAccountData(login.Email)
	if err != nil {
		return 0, gin.H{}, err
	}

	// generate encrypted password
	encryptedPassword, err = utils.GenerateEncryptedPassword([]byte(login.Password))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Generating encrypted password occurs error", err.Error(), http.StatusInternalServerError)
	}

	// update password of reporter account
	ra.Password = encryptedPassword
	err = mc.Storage.UpdateReporterAccount(ra)
	// handle internal server error - can't update the record in reporter_account table
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Updating password occurs error", err.Error(), http.StatusInternalServerError)
	}

	return http.StatusOK, gin.H{"status": "success", "data": LoginBody{
		Email: login.Email,
	}}, nil
}

// ForgetPassword - re-send activation email to the user according to email
func (mc *MembershipController) ForgetPassword(c *gin.Context, mailSender *utils.EmailContext) (int, gin.H, error) {
	type Form struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	var email string
	var err error
	var errorWhere = "MembershipController.ForgetPassword"
	var form Form
	var ra *models.ReporterAccount

	if err = c.Bind(&form); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": Form{
			Email: "email is required",
		}}, nil
	}

	email = form.Email

	// get reporter account by email from reporter_account table
	ra, err = mc.Storage.GetReporterAccountData(email)

	if err != nil {
		return 0, gin.H{}, err
	}

	// reset expire time
	ra.ActExpTime = time.Now().Add(time.Duration(24) * time.Hour)
	if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Updating DB occurs error", err.Error(), http.StatusInternalServerError)
	}

	// re-send the activation email
	err = mailSender.Send(email, activateMailSubject, utils.GenerateActivateMailBody(email, ra.ActivateToken))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Sending activation email occurs error", err.Error(), http.StatusInternalServerError)
	}

	return http.StatusOK, gin.H{"status": "success", "data": Form{
		Email: email,
	}}, nil
}
