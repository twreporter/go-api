package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/utils"
)

func TestSignIn(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// START - test signp endpoint //

	// JSON POST body
	resp = serveHTTP("POST", "/v1/signin", fmt.Sprintf("{\"email\":\"%s\"}", Globs.Defaults.Account),
		"application/json", "")
	assert.Equal(t, resp.Code, 200)

	// form POST body
	resp = serveHTTP("POST", "/v1/signin", fmt.Sprintf("email=%s", Globs.Defaults.Account),
		"application/x-www-form-urlencoded", "")
	assert.Equal(t, resp.Code, 200)

	// neither JSON nor form POST body
	resp = serveHTTP("POST", "/v1/signin", "", "application/text", "")
	assert.Equal(t, resp.Code, 400)

	// sign in with different email
	resp = serveHTTP("POST", "/v1/signin", fmt.Sprintf("{\"email\":\"%s\"}", "contact@twreporter.org"),
		"application/json", "")
	assert.Equal(t, resp.Code, 201)

	// END - test signup endpoint //
}

func TestActivate(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// START - test activate endpoint //
	user := getReporterAccount(Globs.Defaults.Account)

	// test activate
	activateToken := user.ActivateToken
	resp = serveHTTP("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", Globs.Defaults.Account, activateToken), "", "", "")
	fmt.Print(resp.Body)
	assert.Equal(t, resp.Code, 200)

	// test activate fails
	resp = serveHTTP("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", Globs.Defaults.Account, ""), "", "", "")
	assert.Equal(t, resp.Code, 401)
	// END - test activate endpoint //
}

func TestRenewJWT(t *testing.T) {
	user := getReporterAccount(Globs.Defaults.Account)
	jwt, _ := utils.RetrieveToken(user.ID, user.Email)

	// START - test renew jwt endpoint //
	// renew jwt successfully
	resp := serveHTTP("GET", fmt.Sprintf("/v1/token/%v", user.ID), "", "application/json", fmt.Sprintf("Bearer %v", jwt))
	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := struct {
		Status string `json:"status"`
		Data   struct {
			Token     string `json:"token"`
			TokenType string `json:"token_type"`
		} `json:"data"`
	}{}
	json.Unmarshal(body, &res)

	assert.Equal(t, resp.Code, 200)
	assert.Equal(t, res.Status, "success")
	assert.Equal(t, res.Data.TokenType, "Bearer")
	assert.NotEmpty(t, res.Data.Token)

	// fail to renew jwt
	jwt = "testjwt"
	resp = serveHTTP("GET", fmt.Sprintf("/v1/token/%v", user.ID), "", "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 401)
	// End - test renew jwt endpoint //
}

/*
func TestChangePassword(t *testing.T) {
	const userAccount = "test@twreporter.org"
	const userPasswd = "passwd"
	const passwdChanged = "passwdChanged"
	var resp *httptest.ResponseRecorder

	// create an existing active user
	ms := storage.NewGormStorage(DB)
	ra := models.ReporterAccount{
		Account:       userAccount,
		Password:      userPasswd,
		Active:        true,
		ActivateToken: "",
		ActExpTime:    time.Now(),
	}
	user, _ := ms.InsertUserByReporterAccount(ra)

	// lack of JWT in request header
	resp = serveHTTP("POST", "/v1/change-password", fmt.Sprintf("{\"email\":\"%v\",\"password\":\"%v\"}", userAccount, passwdChanged),
		"application/json", "")
	assert.Equal(t, resp.Code, 401)

	// lack of password in the POST BODY
	resp = serveHTTP("POST", "/v1/change-password", fmt.Sprintf("{\"email\":\"%v\"}", userAccount, passwdChanged),
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(user)))
	assert.Equal(t, resp.Code, 400)

	resp = serveHTTP("POST", "/v1/change-password", fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", userAccount, passwdChanged),
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(user)))

	assert.Equal(t, resp.Code, 200)
}

func TestForgetPassword(t *testing.T) {
	var testAccount = Globs.Defaults.Account
	var resp *httptest.ResponseRecorder

	// START - test forget-password endpoint
	// fail test case - not provide the email in the url parameters
	resp = serveHTTP("POST", "/v1/forget-password", "",
		"application/x-www-form-urlencoded", "")
	assert.Equal(t, resp.Code, 400)

	// success test case
	resp = serveHTTP("POST", "/v1/forget-password", fmt.Sprintf("{\"email\":\"%v\"}", testAccount),
		"application/json", "")
	assert.Equal(t, resp.Code, 200)
	// END - test forget-password endpoint
}
*/
