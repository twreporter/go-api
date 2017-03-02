package google

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

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
	r := c.Request
	w := c.Writer
	state := r.FormValue("state")
	if state != utils.Cfg.OauthSettings.GoogleSettings.Statestr {
		log.Errorf("invalid oauth state, expected '%s', got '%s'\n", utils.Cfg.OauthSettings.GoogleSettings.Statestr, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Errorf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Error("Cannot get userinfo from using Google API")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("Error parsing Google user data")
	}
	fmt.Fprintf(w, "Content: %s\n", contents)

	// log.WithFields(log.Fields{"type": constants.Google}).Info("OAuth")
	//
	// // get user data from Google
	// fstring, err := getRemoteUserData(c.Request, c.Writer)
	// if err != nil {
	// 	return
	// }
	//
	// // decode user data returned by Google
	// remoteOauth := models.OAuthAccount{
	// 	Type:      constants.Google,
	// 	AId:       utils.ToNullString(gjson.Get(fstring, "id").Str),
	// 	Email:     utils.ToNullString(gjson.Get(fstring, "email").Str),
	// 	Name:      utils.ToNullString(gjson.Get(fstring, "name").Str),
	// 	FirstName: utils.ToNullString(gjson.Get(fstring, "first_name").Str),
	// 	LastName:  utils.ToNullString(gjson.Get(fstring, "last_name").Str),
	// 	Gender:    getGender(gjson.Get(fstring, "gender").Str),
	// 	Picture:   utils.ToNullString(gjson.Get(fstring, "picture.data.url").Str),
	// }
	//
	// // find the OAuth user from the database
	// matchUser := o.Storage.GetUserDataByOAuth(remoteOauth)
	// // if the user doesn't exist, register the user automatically
	// log.Info("matchUser: ", matchUser)
	// if matchUser.ID == 0 {
	// 	fmt.Println("is zero value", constants.Google)
	// 	o.Storage.InsertUserByOAuth(remoteOauth)
	// } else {
	// 	// update existing OAuth data
	// 	o.Storage.UpdateOAuthData(remoteOauth)
	// }
	//
	// c.Writer.Write([]byte(utils.RetrieveToken(matchUser.Privilege,
	// 	matchUser.FirstName.String, matchUser.LastName.String, matchUser.Email.String)))
	//
	// log.Info("parseResponseBody: %s\n", fstring)
}

// getRemoteUserData fetched user data from Google
func getRemoteUserData(r *http.Request, w http.ResponseWriter) (string, error) {

	oauthStateString := utils.Cfg.OauthSettings.GoogleSettings.Statestr
	loginPath := utils.Cfg.AppSettings.Path + "/login"

	// get Google OAuth Token
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Warn("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", errors.New("invalid oauth state")
	}
	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Warn("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", err
	}

	// get user data from Google
	resp, err := http.Get("https://graph.facebook.com/v2.8/me?fields=id,name,email,picture,birthday,first_name,last_name,gender&access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Warn("Get: %s\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", err
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return "", err
	}

	return string(response), nil
}

func getGender(s string) sql.NullString {
	var ngender sql.NullString
	switch s {
	case "":
		ngender = utils.GetNullString()
	case "male":
		ngender = utils.ToNullString(constants.GenderMale)
	case "female":
		ngender = utils.ToNullString(constants.GenderFemale)
	default:
		// Other gender
		ngender = utils.ToNullString(constants.GenderOthers)
	}
	return ngender
}
