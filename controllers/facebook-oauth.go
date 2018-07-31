package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// Facebook ...
type Facebook struct {
	Storage   storage.MembershipStorage
	oauthConf *oauth2.Config
}

// SetRoute set endpoints for serving facebook oauth
func (f Facebook) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	group.GET("/auth/facebook", middlewares.SetCacheControl("cache,no-store"), f.BeginAuth)
	group.GET("/auth/facebook/callback", middlewares.SetCacheControl("no-store"), f.Authenticate)
	return group
}

// Close ...
func (f Facebook) Close() error {
	return nil
}

// InitOauthConfig initialize facebook oauth config
func (f *Facebook) InitOauthConfig(destination string) {
	consumerSettings := utils.Cfg.ConsumerSettings
	if destination == "" {
		destination = consumerSettings.Protocol + "://" + consumerSettings.Host + ":" + consumerSettings.Port
	}

	destination = url.QueryEscape(destination)
	redirectURL := utils.Cfg.OauthSettings.FacebookSettings.URL + "?destination=" + destination

	if f.oauthConf == nil {
		f.oauthConf = &oauth2.Config{
			ClientID:     utils.Cfg.OauthSettings.FacebookSettings.ID,
			ClientSecret: utils.Cfg.OauthSettings.FacebookSettings.Secret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"public_profile", "email"},
			Endpoint:     facebook.Endpoint,
		}
	} else {
		f.oauthConf.RedirectURL = redirectURL
	}
}

// BeginAuth redirects user to the Facebook Authentication
func (f *Facebook) BeginAuth(c *gin.Context) {
	destination := c.Query("destination")
	f.InitOauthConfig(destination)
	URL, err := url.Parse(f.oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Error("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", f.oauthConf.ClientID)
	parameters.Add("scope", strings.Join(f.oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", f.oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", utils.Cfg.OauthSettings.FacebookSettings.Statestr)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

// Authenticate requests the user profile from Facebook
func (f *Facebook) Authenticate(c *gin.Context) {
	var appErr *models.AppError
	var destination string
	var err error
	var matchUser models.User
	var remoteOAuth models.OAuthAccount
	var fstring string
	var token string

	defer func() {
		// for better ux, redirect users to destination due to internal server error
		if err != nil {
			appErr = appErrorTypeAssertion(err)
			log.Error(appErr.Error())
			c.Redirect(http.StatusTemporaryRedirect, destination)
		}
	}()

	destination = c.Query("destination")

	// get user data from Facebook
	if fstring, err = f.GetRemoteUserData(c.Request, c.Writer); err != nil {
		return
	}

	// decode user data returned by Facebook
	remoteOAuth = models.OAuthAccount{
		Type:      constants.Facebook,
		AId:       utils.ToNullString(gjson.Get(fstring, "id").Str),
		Email:     utils.ToNullString(gjson.Get(fstring, "email").Str),
		Name:      utils.ToNullString(gjson.Get(fstring, "name").Str),
		FirstName: utils.ToNullString(gjson.Get(fstring, "first_name").Str),
		LastName:  utils.ToNullString(gjson.Get(fstring, "last_name").Str),
		Gender:    utils.GetGender(gjson.Get(fstring, "gender").Str),
		Picture:   utils.ToNullString(gjson.Get(fstring, "picture.data.url").Str),
	}

	// get the record from o_auth_accounts table
	_, err = f.Storage.GetOAuthData(remoteOAuth.AId, remoteOAuth.Type)

	// oAuth account is not existed
	// sign in by oauth for the first time
	if err != nil {
		appErr = appErrorTypeAssertion(err)

		// internal server error occurs
		if appErr.StatusCode != http.StatusNotFound {
			return
		}

		// email is provided in oAuth response
		if remoteOAuth.Email.Valid {
			// get the record from users table
			matchUser, err = f.Storage.GetUserByEmail(remoteOAuth.Email.String)

			// record is not existed in users table
			if err != nil {
				if appErr.StatusCode != http.StatusNotFound {
					return
				}

				// no record in users table with this email
				// create a record in users table
				// and create a record in o_auth_accounts table
				matchUser = f.Storage.InsertUserByOAuth(remoteOAuth)
			} else {
				// record existed in user table
				// create record in o_auth_accounts table
				// and connect it to the user record
				remoteOAuth.UserID = matchUser.ID
				err = f.Storage.InsertOAuthAccount(remoteOAuth)

				if err != nil {
					return
				}
			}
		} else {
			// email is not provided in oAuth response
			// create a record in users table
			// and also create a record in o_auth_accounts table
			matchUser = f.Storage.InsertUserByOAuth(remoteOAuth)
		}
	} else {
		// user signed in before
		matchUser, err = f.Storage.GetUserDataByOAuth(remoteOAuth)
		if err != nil {
			return
		}

		// update existing OAuth data
		_, err = f.Storage.UpdateOAuthData(remoteOAuth)
		if err != nil {
			return
		}
	}

	token, err = utils.RetrieveToken(matchUser.ID, matchUser.Email.String)
	if err != nil {
		return
	}

	u, err := url.Parse(destination)
	var secure bool
	secure = false

	if u.Scheme == "https" {
		secure = true
	}

	parameters := u.Query()
	parameters.Add("login", "facebook")

	u.RawQuery = parameters.Encode()
	destination = u.String()

	authJSON := &models.AuthenticatedResponse{ID: matchUser.ID, Privilege: matchUser.Privilege, FirstName: matchUser.FirstName.String, LastName: matchUser.LastName.String, Email: matchUser.Email.String, Jwt: token}
	authResp, _ := json.Marshal(authJSON)

	c.SetCookie("auth_info", string(authResp), 100, u.Path, utils.Cfg.ConsumerSettings.Domain, secure, true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}

// GetRemoteUserData fetched user data from Facebook
func (f *Facebook) GetRemoteUserData(r *http.Request, w http.ResponseWriter) (string, error) {

	oauthStateString := utils.Cfg.OauthSettings.FacebookSettings.Statestr

	// get Facebook OAuth Token
	state := r.FormValue("state")
	if state != oauthStateString {
		return "", models.NewAppError("Facebook.GetRemoteUserData", "invalid oauth state", fmt.Sprintf("invalid oauth state, expected '%s', actual '%s'\n", oauthStateString, state), http.StatusInternalServerError)
	}
	code := r.FormValue("code")

	token, err := f.oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", models.NewAppError("Facebook.GetRemoteUserData", "code exchange failed", err.Error(), http.StatusInternalServerError)
	}

	// get user data from Facebook
	resp, err := http.Get("https://graph.facebook.com/v2.8/me?fields=id,name,email,picture,birthday,first_name,last_name,gender&access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		return "", models.NewAppError("Facebook.GetRemoteUserData", "cannot get user info using Facebook API", err.Error(), http.StatusInternalServerError)
	}

	defer utils.Check(resp.Body.Close)

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", models.NewAppError("Facebook.GetRemoteUserData", "error parsing Facebook user data", err.Error(), http.StatusInternalServerError)
	}

	return string(response), nil
}
