package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google ...
type Google struct {
	Storage   storage.MembershipStorage
	oauthConf *oauth2.Config
}

// SetRoute set endpoints for serving google oauth
func (g Google) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	group.GET("/auth/google", middlewares.SetCacheControl("no-store"), g.BeginAuth)
	group.GET("/auth/google/callback", middlewares.SetCacheControl("no-store"), g.Authenticate)
	return group
}

// Close ...
func (g Google) Close() error {
	return nil
}

// InitOauthConfig initialize google oauth config
func (g *Google) InitOauthConfig(destination string) {
	consumerSettings := utils.Cfg.ConsumerSettings
	if destination == "" {
		destination = consumerSettings.Protocol + "://" + consumerSettings.Host + ":" + consumerSettings.Port + "/activate"
	}

	destination = url.QueryEscape(destination)
	redirectURL := utils.Cfg.OauthSettings.GoogleSettings.URL + "?destination=" + destination

	if g.oauthConf == nil {
		g.oauthConf = &oauth2.Config{
			ClientID:     utils.Cfg.OauthSettings.GoogleSettings.ID,
			ClientSecret: utils.Cfg.OauthSettings.GoogleSettings.Secret,
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

	url := g.oauthConf.AuthCodeURL(utils.Cfg.OauthSettings.GoogleSettings.Statestr)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Authenticate requests the user profile from Google
func (g *Google) Authenticate(c *gin.Context) {
	var appErr models.AppError
	var err error
	var matchUser models.User
	var remoteOAuth models.OAuthAccount

	destination := c.Query("destination")

	// get user data from Google
	fstring, err := g.GetRemoteUserData(c.Request, c.Writer)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, destination)
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
		appErr = err.(models.AppError)

		// return internal server error
		if appErr.StatusCode != http.StatusNotFound {
			c.JSON(appErr.StatusCode, gin.H{"status": "error", "error": appErr.Error()})

			return
		}

		// email is provided in oAuth response
		if remoteOAuth.Email.Valid {
			// get the record from users table
			matchUser, err = g.Storage.GetUserByEmail(remoteOAuth.Email.String)

			// record is not existed in users table
			if err != nil {
				appErr = err.(models.AppError)

				// return internal server error
				if appErr.StatusCode != http.StatusNotFound {
					c.JSON(appErr.StatusCode, gin.H{"status": "error", "error": appErr.Error()})
					return
				}

				// no record in users table with this email
				// create a record in users table
				// and create a record in o_auth_accounts table
				matchUser = g.Storage.InsertUserByOAuth(remoteOAuth)
			} else {
				// record existed in user table
				// create record in o_auth_accounts table
				// and connect it to the user record
				remoteOAuth.UserID = matchUser.ID
				err = g.Storage.InsertOAuthAccount(remoteOAuth)

				// return internal server error
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
					return
				}
			}
		} else {
			// email is not provided in oAuth response
			// create a record in users table
			// and also create a record in o_auth_accounts table
			matchUser = g.Storage.InsertUserByOAuth(remoteOAuth)
		}
	} else {
		// user signed in before
		matchUser, err = g.Storage.GetUserDataByOAuth(remoteOAuth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}

		// update existing OAuth data
		_, err = g.Storage.UpdateOAuthData(remoteOAuth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
			return
		}
	}

	token, err := utils.RetrieveToken(matchUser.ID, matchUser.Email.String)
	if err != nil {
		c.JSON(appErr.StatusCode, gin.H{"status": "error", "error": appErr.Error()})
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

	c.SetCookie("auth_info", string(authResp), 100, u.Path, utils.Cfg.ConsumerSettings.Domain, secure, true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}

// GetRemoteUserData fetched user data from Google
func (g *Google) GetRemoteUserData(r *http.Request, w http.ResponseWriter) (string, error) {

	state := r.FormValue("state")
	if state != utils.Cfg.OauthSettings.GoogleSettings.Statestr {
		log.Warnf("controllers.oauth.google.getRemoteUserData. Invalid oauth state, expected '%s', got '%s'\n", utils.Cfg.OauthSettings.GoogleSettings.Statestr, state)
		return "", models.NewAppError("OAuth state", "controllers.oauth.google", "Invalid oauth state", 500)
	}

	code := r.FormValue("code")
	token, err := g.oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Warnf("controllers.oauth.google.getRemoteUserData. Code exchange failed with '%s'\n", err)
		return "", models.NewAppError("Code exchange failed", "controllers.oauth.google", err.Error(), 500)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Warn("controllers.oauth.google.getRemoteUserData. Cannot get user info using Google API")
		return "", models.NewAppError("Cannot get user info using Google API", "controllers.oauth.google", err.Error(), 500)
	}

	defer utils.Check(response.Body.Close)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Warn("controllers.oauth.google.getRemoteUserData. Error parsing Google user data")
		return "", models.NewAppError("Error parsing Google user data", "controllers.oauth.google", err.Error(), 500)
	}
	// fmt.Fprintf(w, "Content: %s\n", contents)

	return string(contents), nil
}
