package google

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
	Storage *storage.UserStorage
}

func initOauthConfig() {
	if oauthConf == nil {
		oauthConf = &oauth2.Config{
			ClientID:     utils.Cfg.OauthSettings.GoogleSettings.Id,
			ClientSecret: utils.Cfg.OauthSettings.GoogleSettings.Secret,
			RedirectURL:  utils.Cfg.OauthSettings.GoogleSettings.Url,
			Scopes: []string{
				"profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
				"email",
				"openid",
			},
			Endpoint: google.Endpoint,
		}
	}
}

// BeginAuth redirects user to the Google Authentication
func (o Google) BeginAuth(c *gin.Context) {
	initOauthConfig()

	url := oauthConf.AuthCodeURL(utils.Cfg.OauthSettings.GoogleSettings.Statestr)
	log.Info("BeginAuth")
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

// Authenticate requests the user profile from Google
func (o Google) Authenticate(c *gin.Context) {

	log.WithFields(log.Fields{"type": constants.Google}).Info("OAuth")

	// get user data from Google
	fstring, err := getRemoteUserData(c.Request, c.Writer)
	if err != nil {
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

	log.WithFields(log.Fields{"type": remoteOauth.AId}).Info("Fields")

	// find the OAuth user from the database
	matchUser := o.Storage.GetUserDataByOAuth(remoteOauth)
	// if the user doesn't exist, register the user automatically
	log.Info("matchUser: ", matchUser)
	if matchUser.ID == 0 {
		fmt.Println("is zero value", constants.Google)
		o.Storage.InsertUserByOAuth(remoteOauth)
	} else {
		// update existing OAuth data
		o.Storage.UpdateOAuthData(remoteOauth)
	}

	c.Writer.Write([]byte(utils.RetrieveToken(matchUser.Privilege,
		matchUser.FirstName.String, matchUser.LastName.String, matchUser.Email.String)))

	log.Info("parseResponseBody: %s\n", fstring)
}

// getRemoteUserData fetched user data from Google
func getRemoteUserData(r *http.Request, w http.ResponseWriter) (string, error) {
	state := r.FormValue("state")
	if state != utils.Cfg.OauthSettings.GoogleSettings.Statestr {
		log.Errorf("invalid oauth state, expected '%s', got '%s'\n", utils.Cfg.OauthSettings.GoogleSettings.Statestr, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Errorf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return "", err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Error("Cannot get user info from using Google API")
		return "", err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error parsing Google user data")
		return "", err
	}
	// fmt.Fprintf(w, "Content: %s\n", contents)

	return string(contents), nil
}
