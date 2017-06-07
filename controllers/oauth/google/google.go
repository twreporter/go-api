package google

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConf *oauth2.Config

// Google ...
type Google struct {
	Storage storage.MembershipStorage
}

func (gc Google) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	group.GET("/auth/google", gc.BeginAuth)
	group.GET("/auth/google/callback", gc.Authenticate)
	return group
}

func (gc Google) Close() error {
	return nil
}

func initOauthConfig(location string, domain string) {
	consumerSettings := utils.Cfg.ConsumerSettings
	if location == "" {
		location = consumerSettings.Protocal + "://" + consumerSettings.Host + ":" + consumerSettings.Port
	}

	if domain == "" {
		domain = consumerSettings.Domain
	}

	location = url.QueryEscape(location)
	redirectURL := utils.Cfg.OauthSettings.GoogleSettings.URL + "?location=" + location + "&domain=" + domain

	if oauthConf == nil {
		oauthConf = &oauth2.Config{
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
		oauthConf.RedirectURL = redirectURL
	}
}

// BeginAuth redirects user to the Google Authentication
func (o Google) BeginAuth(c *gin.Context) {
	location := c.Query("location")
	domain := c.Query("domain")
	initOauthConfig(location, domain)
	url := oauthConf.AuthCodeURL(utils.Cfg.OauthSettings.GoogleSettings.Statestr)
	log.Info("BeginAuth")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Authenticate requests the user profile from Google
func (o Google) Authenticate(c *gin.Context) {

	log.Info("controllers.oauth.google.authenticate. OAuth type: ", constants.Google)
	location := c.Query("location")
	domain := c.Query("domain")

	// get user data from Google
	fstring, err := getRemoteUserData(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":      "unauthorized",
			"type":        constants.Google,
			"description": "Cannot get user data from Google.",
		})
		return
	}

	// decode user data returned by Google
	remoteOauth := models.OAuthAccount{
		Type:      constants.Google,
		AId:       utils.ToNullString(gjson.Get(fstring, "id").Str),
		Email:     utils.ToNullString(gjson.Get(fstring, "email").Str),
		Name:      utils.ToNullString(gjson.Get(fstring, "name").Str),
		FirstName: utils.ToNullString(gjson.Get(fstring, "given_name").Str),
		LastName:  utils.ToNullString(gjson.Get(fstring, "family_name").Str),
		Gender:    utils.GetGender(gjson.Get(fstring, "gender").Str),
		Picture:   utils.ToNullString(gjson.Get(fstring, "picture").Str),
	}

	log.WithFields(log.Fields{
		"Type": constants.Google,
		"AId":  remoteOauth.AId,
	}).Info("controllers.oauth.google.authenticate. OAuth Login")

	// find the OAuth user from the database
	matchUser, err := o.Storage.GetUserDataByOAuth(remoteOauth)
	// if the user doesn't exist, register the user automatically
	if err != nil {
		log.WithFields(log.Fields{
			"Type": constants.Google,
			"AId":  remoteOauth.AId,
			"Name": remoteOauth.Name,
		}).Info("controllers.oauth.google.authenticate. Create OAuth User")
		matchUser = o.Storage.InsertUserByOAuth(remoteOauth)
	} else {
		// update existing OAuth data
		o.Storage.UpdateOAuthData(remoteOauth)
	}

	token, err := utils.RetrieveToken(matchUser.ID, matchUser.Privilege,
		matchUser.FirstName.String, matchUser.LastName.String, matchUser.Email.String)

	if err != nil {
		log.Error("controllers.oauth.google.authenticate_parse_location_error", err.Error())
		c.JSON(500, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	u, err := url.Parse(location)
	var secure bool
	secure = false

	if u.Scheme == "https" {
		secure = true
	}

	parameters := u.Query()
	parameters.Add("login", "google")
	u.RawQuery = parameters.Encode()
	location = u.String()

	c.SetCookie("token", token, 100, u.Path, domain, secure, true)
	c.Redirect(http.StatusTemporaryRedirect, location)
}

// getRemoteUserData fetched user data from Google
func getRemoteUserData(r *http.Request, w http.ResponseWriter) (string, error) {
	loginPath := utils.Cfg.AppSettings.Path + "/login"

	state := r.FormValue("state")
	if state != utils.Cfg.OauthSettings.GoogleSettings.Statestr {
		log.Warnf("controllers.oauth.google.getRemoteUserData. Invalid oauth state, expected '%s', got '%s'\n", utils.Cfg.OauthSettings.GoogleSettings.Statestr, state)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", models.NewAppError("OAuth state", "controllers.oauth.google", "Invalid oauth state", 500)
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Warnf("controllers.oauth.google.getRemoteUserData. Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", models.NewAppError("Code exchange failed", "controllers.oauth.google", err.Error(), 500)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Warn("controllers.oauth.google.getRemoteUserData. Cannot get user info using Google API")
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", models.NewAppError("Cannot get user info using Google API", "controllers.oauth.google", err.Error(), 500)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Warn("controllers.oauth.google.getRemoteUserData. Error parsing Google user data")
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", models.NewAppError("Error parsing Google user data", "controllers.oauth.google", err.Error(), 500)
	}
	// fmt.Fprintf(w, "Content: %s\n", contents)

	return string(contents), nil
}
