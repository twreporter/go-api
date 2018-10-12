package controllers

import (
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

const (
	// TODO: Update later after UX designed
	authErrorPage = "https://www.twreporter.org/"

	defaultRedirectPage = "https://www.twreporter.org/"

	idTokenExpiration = 60 * 60 * 24 * 30 * 6
)

var defaultDomain = "." + globals.Conf.App.Domain
var defaultPath = "/"

// SignIn - send email containing sign-in information to the client
func (mc *MembershipController) SignIn(c *gin.Context, mailSender *utils.EmailContext) (int, gin.H, error) {
	// SignInBody is to store POST body
	type SignInBody struct {
		Email       string `json:"email" form:"email" binding:"required"`
		Destination string `json:"destination" form:"destination"`
	}

	const errorWhere = "MembershipController.SignIn"
	const SignInMailSubject = "登入報導者"
	const activateHost = "www.twreporter.org"
	var activeToken string
	var appErr *models.AppError
	var email string
	var err error
	var matchedUser models.User
	var ra models.ReporterAccount
	var signIn SignInBody
	var statusCode int

	// extract email and password field in POST body
	if err = c.Bind(&signIn); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": SignInBody{
			Email:       "email is required",
			Destination: "destination is optional",
		}}, nil
	}

	email = signIn.Email

	// Check if mail address is not malform
	_, err = mail.ParseAddress(email)
	if err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": SignInBody{
			Email: "email is malform",
		}}, nil
	}

	// generate active token
	activeToken, err = utils.GenerateRandomString(8)
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Generating active token occurs error", err.Error(), http.StatusInternalServerError)
	}

	// get reporter account by email from reporter_account table
	ra, err = mc.Storage.GetReporterAccountData(email)
	// account is already signed in before
	if err == nil {
		// update active token and token expire time
		ra.ActivateToken = activeToken
		ra.ActExpTime = time.Now().Add(time.Duration(15) * time.Minute)
		if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Updating DB occurs error", err.Error(), http.StatusInternalServerError)
		}

		statusCode = http.StatusOK
	} else {
		// account is not signed in before
		appErr = err.(*models.AppError)

		// internal server error
		if appErr.StatusCode != http.StatusNotFound {
			return 0, gin.H{}, appErr
		}

		ra = models.ReporterAccount{
			Email:         email,
			ActivateToken: activeToken,
			// expire time is one day
			ActExpTime: time.Now().Add(time.Duration(15) * time.Minute),
		}

		// try to find record by email in users table
		matchedUser, err = mc.Storage.GetUserByEmail(email)
		// the user record is not existed
		if err != nil {
			// create records both in reporter_accounts and users table
			_, err = mc.Storage.InsertUserByReporterAccount(ra)
			if err != nil {
				return 0, gin.H{}, models.NewAppError(errorWhere, "Inserting new record into DB occurs error", err.Error(), http.StatusInternalServerError)
			}

		} else {
			// if user existed,
			// create a record in reporter_accounts table
			// and connect these two records
			ra.UserID = matchedUser.ID
			err = mc.Storage.InsertReporterAccount(ra)
			if err != nil {
				return 0, gin.H{}, models.NewAppError(errorWhere, "Inserting new record into DB occurs error", err.Error(), http.StatusInternalServerError)
			}
		}

		statusCode = http.StatusCreated
	}

	// send activation email
	err = mailSender.Send(email, SignInMailSubject, utils.GenerateActivateMailBody(email, activeToken, signIn.Destination, activateHost))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Sending activation email occurs error", err.Error(), http.StatusInternalServerError)
	}

	return statusCode, gin.H{"status": "success", "data": SignInBody{
		Email:       email,
		Destination: signIn.Destination,
	}}, nil
}

// Activate - validate the reporter account
// if validated, then sign in successfully,
// otherwise, sign in unsuccessfully.
func (mc *MembershipController) Activate(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.Activate"
	var appErr *models.AppError
	var err error
	var jwt string
	var ra models.ReporterAccount
	var user models.User

	email := c.Query("email")
	token := c.Query("token")

	// get reporter account by email from reporter_account table
	if ra, err = mc.Storage.GetReporterAccountData(email); err != nil {
		appErr = err.(*models.AppError)
		return 0, gin.H{}, models.NewAppError(errorWhere, fmt.Sprintf("Activating email %s occurs error", email), err.Error(), appErr.StatusCode)
	}

	// check expire time
	if ra.ActExpTime.Sub(time.Now()) < time.Duration(0) {
		return 0, gin.H{}, models.NewAppError(errorWhere, "ActivateToken is expired", "", http.StatusUnauthorized)
	}

	// check token
	if ra.ActivateToken == token {
		// set active expire time to now
		// which make sure the same token not be signed in again
		ra.ActExpTime = time.Now()

		// handle internal server error - can't update the record in reporter_account table
		if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Updating reporter_account occurs error", err.Error(), http.StatusInternalServerError)
		}

		// handle internal server error - can't query the record from users table
		if user, err = mc.Storage.GetUserDataByReporterAccount(ra); err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Querying reporter_account occurs error", err.Error(), http.StatusInternalServerError)
		}

		// handle internal server error - cannot generate JWT
		if jwt, err = utils.RetrieveV1Token(user.ID, user.Email.String); err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Generating JWT occurs error", err.Error(), http.StatusInternalServerError)
		}

		return 200, gin.H{"status": "success",
			"id": user.ID, "privilege": user.Privilege, "firstname": user.FirstName.String,
			"lastname": user.LastName.String, "email": user.Email.String, "jwt": jwt}, nil
	}

	return 0, gin.H{}, models.NewAppError(errorWhere, "Token is invalid", "", http.StatusUnauthorized)
}

// RenewJWT - validate the old JWT,
// if validated, then renew a new one
func (mc *MembershipController) RenewJWT(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.RenewJWT"
	var appErr *models.AppError
	var err error
	var jwt string
	var user models.User

	userID := c.Param("userID")

	if user, err = mc.Storage.GetUserByID(userID); err != nil {
		appErr = err.(*models.AppError)
		return 0, gin.H{}, models.NewAppError(errorWhere, "Renewing JWT occurs error", appErr.Error(), appErr.StatusCode)
	}

	if jwt, err = utils.RetrieveV1Token(user.ID, user.Email.String); err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Generating JWT occurs error", err.Error(), http.StatusInternalServerError)
	}

	return http.StatusOK, gin.H{"status": "success", "data": struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type"`
	}{
		Token:     jwt,
		TokenType: "Bearer",
	}}, nil
}

// SignInV2 - send email containing sign-in information to the client
func (mc *MembershipController) SignInV2(c *gin.Context, mailSender *utils.EmailContext) (int, gin.H, error) {
	// SignInBody is to store POST body
	type SignInBody struct {
		Email       string `json:"email" form:"email" binding:"required"`
		Destination string `json:"destination" form:"destination"`
	}

	const errorWhere = "MembershipController.SignInV2"
	const SignInMailSubject = "登入報導者"
	const activateHost = "go-api.twreporter.org"
	var activeToken string
	var appErr *models.AppError
	var email string
	var err error
	var matchedUser models.User
	var ra models.ReporterAccount
	var signIn SignInBody
	var statusCode int

	// extract email and password field in POST body
	if err = c.Bind(&signIn); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": SignInBody{
			Email:       "email is required",
			Destination: "destination is optional",
		}}, nil
	}

	email = signIn.Email

	// Check if mail address is not malform
	_, err = mail.ParseAddress(email)
	if err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": SignInBody{
			Email: "email is malform",
		}}, nil
	}

	// generate active token
	activeToken, err = utils.GenerateRandomString(8)
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Generating active token occurs error", err.Error(), http.StatusInternalServerError)
	}

	// get reporter account by email from reporter_account table
	ra, err = mc.Storage.GetReporterAccountData(email)
	// account is already signed in before
	if err == nil {
		// update active token and token expire time
		ra.ActivateToken = activeToken
		ra.ActExpTime = time.Now().Add(time.Duration(15) * time.Minute)
		if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
			return 0, gin.H{}, models.NewAppError(errorWhere, "Updating DB occurs error", err.Error(), http.StatusInternalServerError)
		}

		statusCode = http.StatusOK
	} else {
		// account is not signed in before
		appErr = err.(*models.AppError)

		// internal server error
		if appErr.StatusCode != http.StatusNotFound {
			return 0, gin.H{}, appErr
		}

		ra = models.ReporterAccount{
			Email:         email,
			ActivateToken: activeToken,
			// expire time is one day
			ActExpTime: time.Now().Add(time.Duration(15) * time.Minute),
		}

		// try to find record by email in users table
		matchedUser, err = mc.Storage.GetUserByEmail(email)
		// the user record is not existed
		if err != nil {
			// create records both in reporter_accounts and users table
			_, err = mc.Storage.InsertUserByReporterAccount(ra)
			if err != nil {
				return 0, gin.H{}, models.NewAppError(errorWhere, "Inserting new record into DB occurs error", err.Error(), http.StatusInternalServerError)
			}

		} else {
			// if user existed,
			// create a record in reporter_accounts table
			// and connect these two records
			ra.UserID = matchedUser.ID
			err = mc.Storage.InsertReporterAccount(ra)
			if err != nil {
				return 0, gin.H{}, models.NewAppError(errorWhere, "Inserting new record into DB occurs error", err.Error(), http.StatusInternalServerError)
			}
		}

		statusCode = http.StatusCreated
	}

	// send activation email
	err = mailSender.Send(email, SignInMailSubject, utils.GenerateActivateMailBody(email, activeToken, signIn.Destination, activateHost))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Sending activation email occurs error", err.Error(), http.StatusInternalServerError)
	}

	return statusCode, gin.H{"status": "success", "data": SignInBody{
		Email:       email,
		Destination: signIn.Destination,
	}}, nil
}

// ActivateV2 - validate the reporter account
// if validated, then sign in successfully,
// otherwise, sign in unsuccessfully.
func (mc *MembershipController) ActivateV2(c *gin.Context) {
	const errorWhere = "MembershipController.ActivateV2"
	var err error
	var ra models.ReporterAccount
	var user models.User

	email := c.Query("email")
	token := c.Query("token")
	destination := c.Query("destination")

	// If destination is unavailable or invalid, redirect back to main site.
	u, err := url.Parse(destination)
	if nil != err {
		destination = defaultRedirectPage
		u, _ = url.Parse(destination)
	}

	// Error clean up
	defer func() {
		if nil != err {
			appErr := err.(*models.AppError)
			log.Error(appErr.Error())

			//Always redirect to a designated page
			c.Redirect(http.StatusTemporaryRedirect, authErrorPage)
		}
	}()

	// get reporter account by email from reporter_account table
	if ra, err = mc.Storage.GetReporterAccountData(email); err != nil {
		return
	}

	// check expire time
	if ra.ActExpTime.Sub(time.Now()) < time.Duration(0) {
		err = models.NewAppError(errorWhere, "", "ActivateToken is expired", http.StatusUnauthorized)
		return
	}

	// validate token
	if ra.ActivateToken != token {
		err = models.NewAppError(errorWhere, "", "Token is invalid", http.StatusUnauthorized)
		return
	}

	// set active expire time to now to ensure the same token only being used once
	ra.ActExpTime = time.Now()

	// Error occurs during updating the record in reporter_account table
	if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
		return
	}

	// Error occurs during querying the record from users table
	if user, err = mc.Storage.GetUserDataByReporterAccount(ra); err != nil {
		return
	}

	// Create id token for jwt endpoint retrival
	idToken, err := utils.RetrieveV2IDToken(user.ID, user.Email.ValueOrZero(), user.FirstName.ValueOrZero(), user.LastName.ValueOrZero(), idTokenExpiration)
	if nil != err {
		log.Error(errorWhere + "(): " + err.Error())
		idToken = "twreporter-id-token"
	}

	// Setup Set-Cookie header in response header

	// Determine cookie property
	secure := false

	if "https" == u.Scheme {
		secure = true
	}

	c.SetCookie("id_token", idToken, idTokenExpiration, defaultPath, defaultDomain, secure, true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}

// TokenDispatch returns the `access_token` in payload for frontend server
func (mc *MembershipController) TokenDispatch(c *gin.Context) {
	const acccessTokenExpiration = 60 * 60 * 24 * 14 // 2week

	errorWhere := "MembershipController.TokenDispatch"

	type reqBody struct {
		UserID uint `json:"user_id"`
	}

	// Validate the request body
	body := reqBody{}
	if err := c.Bind(&body); nil != err {
		log.Error(errorWhere + "():" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{
			"user_id": "invalid",
		}})
	}

	user, err := mc.Storage.GetUserByID(fmt.Sprint(body.UserID))
	if nil != err {
		appErr := err.(*models.AppError)
		log.Error(appErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "cannot get user data"})
		return
	}

	jwt, err := utils.RetrieveV2AccessToken(user.ID, user.Email.ValueOrZero(), acccessTokenExpiration)
	if err != nil {
		appErr := err.(*models.AppError)
		log.Error(appErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error occurs during generating access_token JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"jwt": jwt,
	}})
}

// TokenInvalidate deletes the id_token stored in the client side
func (mc *MembershipController) TokenInvalidate(c *gin.Context) {
	cookieName := "id_token"
	invalidateExp := -1

	defer func() {
		c.JSON(http.StatusNoContent, gin.H{})
	}()

	idToken, err := c.Request.Cookie(cookieName)

	// If no valid `id_token` cookie provide, ignore it
	if nil != err {
		return
	}

	c.SetCookie(cookieName, idToken.Value, invalidateExp, defaultPath, defaultDomain, false, true)
}
