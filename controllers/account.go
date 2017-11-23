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

// SignInBody is to store POST body
type SignInBody struct {
	Email string `json:"email" form:"email" binding:"required"`
}

// parseAccountPostBody get email and password value from the POST body
func (mc *MembershipController) parseAccountPostBody(c *gin.Context) (signIn SignInBody, err error) {
	if err = c.Bind(&signIn); err != nil {
		log.Error("err:", err.Error())
		return signIn, models.NewAppError("parseAccountPostBody", "controllers.account.parse_post_body", err.Error(), http.StatusBadRequest)
	}

	return signIn, nil
}

// SignIn - send email containing sign-in information to the client
func (mc *MembershipController) SignIn(c *gin.Context, mailSender *utils.EmailContext) (int, gin.H, error) {
	const errorWhere = "MembershipController.SignIn"
	const SignInMailSubject = "登入報導者"
	var activeToken string
	var appErr models.AppError
	var email string
	var err error
	var matchedUser models.User
	var ra models.ReporterAccount
	var signIn SignInBody
	var statusCode int

	// extract email and password field in POST body
	signIn, err = mc.parseAccountPostBody(c)
	if err != nil {
		appErr := err.(models.AppError)
		return appErr.StatusCode, gin.H{"status": "fail", "data": SignInBody{
			Email: "email is required",
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
		appErr = err.(models.AppError)

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
	err = mailSender.Send(email, SignInMailSubject, utils.GenerateActivateMailBody(email, activeToken))
	if err != nil {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Sending activation email occurs error", err.Error(), http.StatusInternalServerError)
	}

	return statusCode, gin.H{"status": "success", "data": SignInBody{
		Email: email,
	}}, nil
}

// Activate - validate the reporter account
// if validated, then sign in successfully,
// otherwise, sign in unsuccessfully.
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
		// set active expire time to now
		// which make sure the same token not be signed in again
		ra.ActExpTime = time.Now()

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
