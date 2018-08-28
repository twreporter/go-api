package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google ...
type Google struct {
	Storage   storage.MembershipStorage
	oauthConf *oauth2.Config
}

// InitOauthConfig initialize google oauth config
func (g *Google) InitOauthConfig(destination string) {
	if destination == "" {
		destination = viper.GetString("consumersettings.protocol") + "://" +
			viper.GetString("consumersettings.host") + ":" + viper.GetString("consumersettings.port") + "/activate"
	}

	destination = url.QueryEscape(destination)
	redirectURL := fmt.Sprintf("%s://%s:%s/v1/auth/google/callback?destination=%s", viper.GetString("appsettings.protocol"), viper.GetString("appsettings.host"), viper.GetString("appsettings.port"), destination)

	if g.oauthConf == nil {
		g.oauthConf = &oauth2.Config{
			ClientID:     viper.GetString("oauthsettings.googlesettings.id"),
			ClientSecret: viper.GetString("oauthsettings.googlesettings.secret"),
			RedirectURL:  redirectURL,
			Scopes: []string{
				"profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
				"email",
				"openid",
			},
			Endpoint: google.Endpoint,
		}
	} else {
		g.oauthConf.RedirectURL = redirectURL
	}
}

// BeginAuth redirects user to the Google Authentication
func (g *Google) BeginAuth(c *gin.Context) {
	destination := c.Query("destination")

	g.InitOauthConfig(destination)

	url := g.oauthConf.AuthCodeURL(viper.GetString("oauthsettings.googlesettings.statestr"))

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Authenticate requests the user profile from Google
func (g *Google) Authenticate(c *gin.Context) {
	var appErr *models.AppError
	var destination string
	var err error
	var fstring string
	var matchUser models.User
	var remoteOAuth models.OAuthAccount
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

	// get user data from Google
	if fstring, err = g.GetRemoteUserData(c.Request, c.Writer); err != nil {
		return
	}

	// decode user data returned by Google oAuth
	remoteOAuth = models.OAuthAccount{
		Type:      constants.Google,
		AId:       utils.ToNullString(gjson.Get(fstring, "id").Str),
		Email:     utils.ToNullString(gjson.Get(fstring, "email").Str),
		Name:      utils.ToNullString(gjson.Get(fstring, "name").Str),
		FirstName: utils.ToNullString(gjson.Get(fstring, "given_name").Str),
		LastName:  utils.ToNullString(gjson.Get(fstring, "family_name").Str),
		Gender:    utils.GetGender(gjson.Get(fstring, "gender").Str),
		Picture:   utils.ToNullString(gjson.Get(fstring, "picture").Str),
	}

	// get the record from o_auth_accounts table
	_, err = g.Storage.GetOAuthData(remoteOAuth.AId, remoteOAuth.Type)

	// oAuth account is not existed
	// sign in by oauth for the first time
	if err != nil {
		appErr = appErrorTypeAssertion(err)

		// return internal server error
		if appErr.StatusCode != http.StatusNotFound {
			return
		}

		// email is provided in oAuth response
		if remoteOAuth.Email.Valid {
			// get the record from users table
			matchUser, err = g.Storage.GetUserByEmail(remoteOAuth.Email.String)

			// record is not existed in users table
			if err != nil {
				appErr = err.(*models.AppError)

				// return internal server error
				if appErr.StatusCode != http.StatusNotFound {
					return
				}

				// no record in users table with this email
				// create a record in users table
				// and create a record in o_auth_accounts table
				matchUser, err = g.Storage.InsertUserByOAuth(remoteOAuth)
			} else {
				// record existed in user table
				// create record in o_auth_accounts table
				// and connect it to the user record
				remoteOAuth.UserID = matchUser.ID
				err = g.Storage.InsertOAuthAccount(remoteOAuth)

				if err != nil {
					return
				}
			}
		} else {
			// email is not provided in oAuth response
			// create a record in users table
			// and also create a record in o_auth_accounts table
			matchUser, err = g.Storage.InsertUserByOAuth(remoteOAuth)
		}
	} else {
		// user signed in before
		matchUser, err = g.Storage.GetUserDataByOAuth(remoteOAuth)
		if err != nil {
			return
		}

		// update existing OAuth data
		_, err = g.Storage.UpdateOAuthData(remoteOAuth)
		if err != nil {
			return
		}
	}

	token, err = utils.RetrieveToken(matchUser.ID, matchUser.Email.String)
	if err != nil {
		return
	}

	var u *url.URL
	var secure = false
	u, err = url.Parse(destination)

	if u.Scheme == "https" {
		secure = true
	}

	parameters := u.Query()
	parameters.Add("login", "google")

	u.RawQuery = parameters.Encode()
	destination = u.String()

	authJSON := &models.AuthenticatedResponse{ID: matchUser.ID, Privilege: matchUser.Privilege, FirstName: matchUser.FirstName.String, LastName: matchUser.LastName.String, Email: matchUser.Email.String, Jwt: token}
	authResp, _ := json.Marshal(authJSON)

	c.SetCookie("auth_info", string(authResp), 100, u.Path, viper.GetString("consumersettings.domain"), secure, true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}

// GetRemoteUserData fetched user data from Google
func (g *Google) GetRemoteUserData(r *http.Request, w http.ResponseWriter) (string, error) {

	state := r.FormValue("state")
	if state != viper.GetString("oauthsettings.googlesettings.statestr") {
		log.Warnf("controllers.oauth.google.getRemoteUserData. Invalid oauth state, expected '%s', got '%s'\n", viper.GetString("oauthsettings.googlesettings.statestr"), state)
		return "", models.NewAppError("OAuth state", "controllers.oauth.google", "Invalid oauth state", 500)
	}

	code := r.FormValue("code")

	token, err := g.oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", models.NewAppError("Google.GetRemoteUserData", "code exchange failed", err.Error(), http.StatusInternalServerError)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", models.NewAppError("Google.GetRemoteUserData", "cannot get user info using Google API", err.Error(), http.StatusInternalServerError)
	}

	defer utils.Check(response.Body.Close)

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", models.NewAppError("Google.GetRemoteUserData", "error parsing Google user data", err.Error(), http.StatusInternalServerError)
	}

	return string(contents), nil
}
