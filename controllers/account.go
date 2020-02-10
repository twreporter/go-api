package controllers

import (
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

const (
	// TODO: Update later after UX designed
	authErrorPage = "https://www.twreporter.org/"

	defaultRedirectPage = "https://www.twreporter.org/"

	idTokenExpiration = 60 * 60 * 24 * 30 * 6
)

var defaultPath = "/"

// SignInV2 - send email containing sign-in information to the client
func (mc *MembershipController) SignInV2(c *gin.Context) (int, gin.H, error) {
	// SignInBody is to store POST body
	type SignInBody struct {
		Email       string `json:"email" form:"email" binding:"required"`
		Destination string `json:"destination" form:"destination"`
	}

	var activateHost string
	var activeToken string
	var email string
	var err error
	var matchedUser models.User
	var ra models.ReporterAccount
	var signIn SignInBody
	var statusCode int

	switch globals.Conf.Environment {
	case "development":
		activateHost = "localhost"
	case "staging":
		activateHost = "staging-go-api.twreporter.org"
	case "production":
		activateHost = "go-api.twreporter.org"
	default:
		activateHost = "localhost"
	}

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
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "Generating active token occurs error"}, err
	}

	// get reporter account by email from reporter_account table
	ra, err = mc.Storage.GetReporterAccountData(email)
	// account is already signed in before
	if err == nil {
		// update active token and token expire time
		ra.ActivateToken = activeToken
		ra.ActExpTime = time.Now().Add(time.Duration(15) * time.Minute)
		if err = mc.Storage.UpdateReporterAccount(ra); err != nil {
			return http.StatusInternalServerError, gin.H{"status": "error", "message": "Updating DB occurs error"}, err
		}

		statusCode = http.StatusOK
	} else {
		// account is not signed in before
		if !storage.IsNotFound(err) {
			return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("internal server error. %s", err.Error())}, err
		}

		ra = models.ReporterAccount{
			Email:         email,
			ActivateToken: activeToken,
			// expire time is 15 minute
			ActExpTime: time.Now().Add(time.Duration(15) * time.Minute),
		}

		// try to find record by email in users table
		matchedUser, err = mc.Storage.GetUserByEmail(email)
		// the user record is not existed
		if err != nil {
			// create records both in reporter_accounts and users table
			_, err = mc.Storage.InsertUserByReporterAccount(ra)
			if err != nil {
				return http.StatusInternalServerError, gin.H{"status": "error", "message": "Inserting new record into DB occurs error"}, nil
			}
		} else {
			// if user existed,
			// create a record in reporter_accounts table
			// and connect these two records
			ra.UserID = matchedUser.ID
			err = mc.Storage.InsertReporterAccount(ra)
			if err != nil {
				return http.StatusInternalServerError, gin.H{"status": "error", "message": "Inserting new record into DB occurs error"}, nil
			}
		}

		statusCode = http.StatusCreated
	}

	// send activation email
	err = postMailServiceEndpoint(activationReqBody{
		Email: email,
		ActivateLink: fmt.Sprintf("%s://%s:%s/v2/auth/activate?email=%s&token=%s&destination=%s",
			globals.Conf.App.Protocol,
			activateHost,
			globals.Conf.App.Port,
			url.QueryEscape(email),
			url.QueryEscape(activeToken),
			url.QueryEscape(signIn.Destination),
		),
	}, fmt.Sprintf("http://localhost:%s/v1/%s", globals.LocalhostPort, globals.SendActivationRoutePath))

	if err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "Sending activation email occurs error"}, err
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
	var defaultDomain = globals.Conf.App.Domain
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
			// Client side error. Do not trigger error reporting
			log.Infof("%v", err)

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
		err = errors.New("ActivateToken is expired")
		return
	}

	// validate token
	if ra.ActivateToken != token {
		err = errors.New("Token is invalid")
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
		idToken = "twreporter-id-token"
	}

	// Setup Set-Cookie header in response header

	// Determine cookie property
	secure := false

	if "https" == u.Scheme {
		secure = true
	}

	parameters := u.Query()
	parameters.Add("login_time", fmt.Sprintf("%d", time.Now().Unix()))
	u.RawQuery = parameters.Encode()
	destination = u.String()

	c.SetCookie("id_token", idToken, idTokenExpiration, defaultPath, defaultDomain, secure, true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}

// TokenDispatch returns the `access_token` in payload for frontend server
func (mc *MembershipController) TokenDispatch(c *gin.Context) (int, gin.H, error) {
	const acccessTokenExpiration = 60 * 60 * 24 * 14 // 2week
	const idTokenKey = "id_token"
	var accessToken string
	var claims = new(utils.IDTokenJWTClaims)
	var err error
	var idToken string
	var user models.User

	if idToken, err = c.Cookie(idTokenKey); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{"req.Headers.Cookies.id_token": err.Error()}}, nil
	}

	if _, _, err = new(jwt.Parser).ParseUnverified(idToken, claims); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{"req.Headers.Cookies.id_token": err.Error()}}, nil
	}

	user, err = mc.Storage.GetUserByID(fmt.Sprint(claims.UserID))
	if err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "cannot get user data"}, err
	}

	accessToken, err = utils.RetrieveV2AccessToken(user.ID, user.Email.ValueOrZero(), acccessTokenExpiration)
	if err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "Error occurs during generating access_token JWT"}, err
	}

	return http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"jwt": accessToken,
	}}, nil
}

// TokenInvalidate deletes the id_token stored in the client side
func (mc *MembershipController) TokenInvalidate(c *gin.Context) {
	const signInPage = "https://accounts.twreporter.org/signin"
	var defaultDomain = globals.Conf.App.Domain

	cookieName := "id_token"
	invalidateExp := -1

	destination := c.Query("destination")

	if destination == "" {
		destination = signInPage
	}

	// If destination is unavailable or invalid, redirect back to signin page
	u, err := url.Parse(destination)
	if nil != err {
		destination = signInPage
		u, _ = url.Parse(destination)
	}

	c.SetCookie(cookieName, "", invalidateExp, defaultPath, defaultDomain, u.Scheme == "https", true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}
