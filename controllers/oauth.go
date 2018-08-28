package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

const defaultDestination = "https://www.twreporter.org/"

type basicInfo struct {
	AId    models.NullString `json:"id"`
	Email  models.NullString `json:"email"`
	Name   models.NullString `json:"name"`
	Gender models.NullString `json:"gender"`
}

type facebookOauthInfoRaw struct {
	basicInfo
	FirstName  models.NullString `json:"first_name"`
	LastName   models.NullString `json:"last_name"`
	PictureObj struct {
		Data struct {
			URL models.NullString `json:"url"`
		} `json:data`
	} `json:"picture"`
}

// Picture is used by copier to copy PictureObj.Data.URL field to Picture field
func (info *facebookOauthInfoRaw) Picture() models.NullString {
	return info.PictureObj.Data.URL
}

type googleOauthInfoRaw struct {
	basicInfo
	FirstName models.NullString `json:"given_name"`
	LastName  models.NullString `json:"family_name"`
	Picture   models.NullString `json:"picture"`
}

// beginAuth uses sessions to store users'
// 1. state
// 2. destination(go to page)
// and redirect users to oauth server.
func beginAuth(c *gin.Context, conf *oauth2.Config) {
	var state string
	var err error

	destination := c.Request.URL.Query().Get("destination")
	if destination == "" {
		destination = defaultDestination
	}

	log.Infof("destination beginAuth %s", destination)

	if state, err = utils.GenerateRandomString(32); err != nil {
		state = "twreporter-oauth-state"
	}

	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("destination", destination)
	session.Save()

	url := conf.AuthCodeURL(state)
	log.Infof("url %#v", url)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// getOauthUserInfo does the following three things
// 1. validate state
// 2. exchange code to token
// 3. get user info from oauth server by token
func getOauthUserInfo(c *gin.Context, conf *oauth2.Config, userInfoEndpoint string, oauthUser interface{}) error {
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	state := c.Request.URL.Query().Get("state")
	if state != retrievedState {
		log.Warnf("expect state is %s, but actual state is %s", retrievedState, state)
		return models.NewAppError("getOauthUserInfo", "oauth fails", "Invalid oauth state", 500)
	}

	code := c.Request.URL.Query().Get("code")
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return models.NewAppError("getOauthUserInfo", "oauth code exchange failed", err.Error(), http.StatusInternalServerError)
	}

	client := conf.Client(oauth2.NoContext, token)
	response, err := client.Get(userInfoEndpoint)

	if err != nil {
		return models.NewAppError("getOauthUserInfo", "cannot get user info by using Google API", err.Error(), http.StatusInternalServerError)
	}

	defer response.Body.Close()

	userInfo, err := ioutil.ReadAll(response.Body)
	log.Infof("userInfo: %#v", string(userInfo))

	if err != nil {
		return models.NewAppError("getOauthUserInfo", "error parsing user data", err.Error(), http.StatusInternalServerError)
	}

	if err = json.Unmarshal(userInfo, &oauthUser); err != nil {
		return models.NewAppError("getOauthUserInfo", "can not unmarshal user data", err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// In order to avoid from storing user info repeatedly,
// findOrCreateUser handles how to store oauth users in the storage.
func findOrCreateUser(oauthUser models.OAuthAccount, storage storage.MembershipStorage) (user models.User, err error) {
	var appErr *models.AppError

	// get the record from o_auth_accounts table
	_, err = storage.GetOAuthData(oauthUser.AId, oauthUser.Type)

	// oAuth account is not existed
	// sign in by oauth for the first time
	if err != nil {
		appErr = appErrorTypeAssertion(err)

		// return internal server error
		if appErr.StatusCode != http.StatusNotFound {
			return user, err
		}

		// email is provided in oAuth response
		if oauthUser.Email.Valid {
			// get the record from users table
			user, err = storage.GetUserByEmail(oauthUser.Email.String)

			// record is not existed in users table
			if err != nil {
				appErr = err.(*models.AppError)

				// return internal server error
				if appErr.StatusCode != http.StatusNotFound {
					return user, err
				}

				// no record in users table with this email
				// create a record in users table
				// and create a record in o_auth_accounts table
				if user, err = storage.InsertUserByOAuth(oauthUser); err != nil {
					return user, err
				}
			} else {
				// record existed in user table
				// create record in o_auth_accounts table
				// and connect it to the user record
				oauthUser.UserID = user.ID
				if err = storage.InsertOAuthAccount(oauthUser); err != nil {
					return user, err
				}
			}
		} else {
			// email is not provided in oAuth response
			// create a record in users table
			// and also create a record in o_auth_accounts table
			if user, err = storage.InsertUserByOAuth(oauthUser); err != nil {
				return user, err
			}
		}
	} else {
		// user signed in before
		if user, err = storage.GetUserDataByOAuth(oauthUser); err != nil {
			return user, err
		}

		// update existing OAuth data
		if _, err = storage.UpdateOAuthData(oauthUser); err != nil {
			return user, err
		}
	}

	return user, nil
}

// OAuth which stores storage connection and oauth config
type OAuth struct {
	Storage   storage.MembershipStorage
	oauthConf *oauth2.Config
}

// InitGoogleConfig initiates facebook oauth config
func (o *OAuth) InitFacebookConfig() {
	appsettings := viper.GetStringMap("appsettings")
	redirectURL := fmt.Sprintf("%s://%s:%s/v2/auth/facebook/callback", appsettings["protocol"], appsettings["host"], appsettings["port"])
	o.oauthConf = &oauth2.Config{
		ClientID:     viper.GetString("oauthsettings.facebooksettings.id"),
		ClientSecret: viper.GetString("oauthsettings.facebooksettings.secret"),
		RedirectURL:  redirectURL,
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
}

// InitGoogleConfig initiates google oauth config
func (o *OAuth) InitGoogleConfig() {
	appsettings := viper.GetStringMap("appsettings")
	redirectURL := fmt.Sprintf("%s://%s:%s/v2/auth/google/callback", appsettings["protocol"], appsettings["host"], appsettings["port"])
	o.oauthConf = &oauth2.Config{
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
}

// BeginAuth redirects user to the [facebook|google] authentication(login) page
func (o *OAuth) BeginOAuth(c *gin.Context) {
	beginAuth(c, o.oauthConf)
	return
}

// Authenticate handles [google|facebook] oauth of users and redirect them to specific URL they want
// with Set-Cookie response header which contains JWT
func (o *OAuth) Authenticate(c *gin.Context) {
	var destination = "http://testtest.twreporter.org:3000/"
	var err error
	var matchUser models.User
	var oauthType string
	var oauthUser models.OAuthAccount
	var retrievedDestination interface{}
	var session sessions.Session
	var token string
	var userInfoEndpoint string

	session = sessions.Default(c)

	if retrievedDestination = session.Get("destination"); retrievedDestination != nil {
		destination = retrievedDestination.(string)
		if destination == "" {
			destination = defaultDestination
		}
	}

	if o.oauthConf.Endpoint == google.Endpoint {
		var oauthInfo googleOauthInfoRaw
		userInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
		getOauthUserInfo(c, o.oauthConf, userInfoEndpoint, &oauthInfo)
		copier.Copy(&oauthUser, &oauthInfo)
		oauthType = constants.GoogleOAuth
	} else {
		var oauthInfo facebookOauthInfoRaw
		userInfoEndpoint = "https://graph.facebook.com/v2.8/me?fields=id,name,email,picture,birthday,first_name,last_name,gender"
		err = getOauthUserInfo(c, o.oauthConf, userInfoEndpoint, &oauthInfo)
		copier.Copy(&oauthUser, &oauthInfo)
		oauthType = constants.FacebookOAuth
	}

	if err != nil {
		log.Errorf("oauth fails while getting user info from api, error message:\n%s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, destination)
		return
	}

	log.Infof("oauthUser: %#v", oauthUser)

	oauthUser.Type = oauthType

	if matchUser, err = findOrCreateUser(oauthUser, o.Storage); err != nil {
		log.Errorf("oauth fails due to database operation error:\n%s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, destination)
		return
	}

	if token, err = utils.RetrieveToken(matchUser.ID, matchUser.Email.String); err != nil {
		log.Errorf("oauth fails due to generate JWT error:\n%s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, destination)
		return
	}

	var u *url.URL
	var secure = false
	log.Infof(" destination: %s", destination)
	u, err = url.Parse(destination)

	if u.Scheme == "https" {
		secure = true
	}

	parameters := u.Query()
	parameters.Add("login", oauthType)

	u.RawQuery = parameters.Encode()
	destination = u.String()

	log.Infof("authenticate destination: %s", destination)

	authJSON := &models.AuthenticatedResponse{ID: matchUser.ID, Privilege: matchUser.Privilege,
		FirstName: matchUser.FirstName.String, LastName: matchUser.LastName.String,
		Email: matchUser.Email.String}

	// hours to seconds
	maxAge := viper.GetInt("appsettings.expiration") * 60 * 60
	authResp, _ := json.Marshal(authJSON)

	// set domain to .twreporter.org
	// so each hostname of [www|support|tsai-tracker].twreporter.org will be applied
	c.SetCookie("auth_info", string(authResp), maxAge, "/", ".twreporter.org", secure, false)
	c.SetCookie("access_token", token, maxAge, "/", ".twreporter.org", secure, true)
	c.Redirect(http.StatusTemporaryRedirect, destination)
}
