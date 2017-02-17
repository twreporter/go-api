package oauth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var (
	cfg       = configs.GetConfig()
	oauthConf = &oauth2.Config{
		ClientID:     cfg.OAUTH.FACEBOOK.ID,
		ClientSecret: cfg.OAUTH.FACEBOOK.Secret,
		RedirectURL:  cfg.OAUTH.FACEBOOK.URL,
		Scopes:       []string{"public_profile", "email"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = cfg.OAUTH.FACEBOOK.Statestr
	loginPath        = cfg.APP.Path + "/login"
)

// Facebook ...
type Facebook struct {
	Storage *storage.UserStorage
}

// BeginAuth redirects user to the Facebook Authentication
func (o Facebook) BeginAuth(c *gin.Context) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Error("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

// Authenticate requests the user profile from Facebook
func (o Facebook) Authenticate(c *gin.Context) {
	r := c.Request
	w := c.Writer

	log.WithFields(log.Fields{"type": "Facebook"}).Info("OAuth")

	// get Facebook OAuth Token
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Warn("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return
	}
	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Warn("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return
	}

	// get user data from Facebook
	resp, err := http.Get("https://graph.facebook.com/v2.8/me?fields=id,name,email,picture,birthday,first_name,last_name,gender&access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Warn("Get: %s\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, loginPath, http.StatusTemporaryRedirect)
		return
	}

	// decode user data returned by Facebook
	fstring := string(response)
	fid := gjson.Get(fstring, "id").Str
	femail := gjson.Get(fstring, "email").Str
	fname := gjson.Get(fstring, "name").Str
	ffirst := gjson.Get(fstring, "first_name").Str
	flast := gjson.Get(fstring, "last_name").Str
	fgender := gjson.Get(fstring, "gender").Str
	fpicture := gjson.Get(fstring, "picture.data.url").Str
	fbirth := gjson.Get(fstring, "birthday").Str

	log.Info("User structure", fid, femail, fname, ffirst, flast, fgender, fpicture, fbirth)

	// find the OAuth user from the database
	matchOauth := o.Storage.GetOAuthData(fstring)
	// if the user doesn't exist
	log.Info("matchOauth: ", matchOauth, matchOauth.AId)
	if !matchOauth.AId.Valid {
		fmt.Println("is zero value")
		o.Storage.InsertUserByOAuth(models.OAuthAccount{
			Type:      "Facebook",
			AId:       utils.ToNullString(fid),
			Email:     utils.ToNullString(femail),
			Name:      utils.ToNullString(fname),
			FirstName: utils.ToNullString(ffirst),
			LastName:  utils.ToNullString(flast),
			Gender:    utils.ToNullString(fgender),
			Picture:   utils.ToNullString(fpicture),
		})
	}

	w.Write([]byte(utils.RetrieveToken(false, fname, femail)))

	log.Info("parseResponseBody: %s\n", fstring)
}
