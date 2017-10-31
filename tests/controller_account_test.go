package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
)

func TestLogin(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Login successfully
	resp = ServeHTTP("POST", "/v1/login",
		fmt.Sprintf("email=%v&password=%v", DefaultAccount, DefaultPassword),
		"application/x-www-form-urlencoded", "")
	assert.Equal(t, resp.Code, 200)

	// Fail to login
	resp = ServeHTTP("POST", "/v1/login",
		fmt.Sprintf("email=%v&password=wrongpassword", DefaultAccount),
		"application/x-www-form-urlencoded", "") //wrong password
	assert.Equal(t, resp.Code, 401)
}
func TestSignupAndActivate(t *testing.T) {
	const testAccount = "developer@twreporter.org"
	const password = "password"
	var resp *httptest.ResponseRecorder

	// START - test signp endpoint //

	// JSON POST body
	resp = ServeHTTP("POST", "/v1/signup", fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", testAccount, password),
		"application/json", "")
	assert.Equal(t, resp.Code, 201)

	// form POST body
	resp = ServeHTTP("POST", "/v1/signup", fmt.Sprintf("email=%s&password=%s", testAccount, password),
		"application/x-www-form-urlencoded", "")
	assert.Equal(t, resp.Code, 201)

	// neither JSON nor form POST body
	resp = ServeHTTP("POST", "/v1/signup", "", "application/text", "")
	assert.Equal(t, resp.Code, 400)

	// signup an already account
	resp = ServeHTTP("POST", "/v1/signup", fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", DefaultAccount, DefaultPassword),
		"application/json", "")
	assert.Equal(t, resp.Code, 409)

	// END - test signup endpoint //

	// START - test activate endpoint //
	as := storage.NewGormStorage(DB)
	user, _ := as.GetReporterAccountData(testAccount)

	// test activate
	activateToken := user.ActivateToken
	resp = ServeHTTP("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", testAccount, activateToken), "", "", "")
	fmt.Print(resp.Body)

	assert.Equal(t, resp.Code, 200)

	// test activate fails
	resp = ServeHTTP("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", "mika@twreporter.org", ""), "", "", "")
	assert.Equal(t, resp.Code, 404)
	// END - test activate endpoint //
}

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
	resp = ServeHTTP("POST", "/v1/change-password", fmt.Sprintf("{\"email\":\"%v\",\"password\":\"%v\"}", userAccount, passwdChanged),
		"application/json", "")
	assert.Equal(t, resp.Code, 401)

	// lack of password in the POST BODY
	resp = ServeHTTP("POST", "/v1/change-password", fmt.Sprintf("{\"email\":\"%v\"}", userAccount, passwdChanged),
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(user)))
	assert.Equal(t, resp.Code, 400)

	resp = ServeHTTP("POST", "/v1/change-password", fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", userAccount, passwdChanged),
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(user)))

	assert.Equal(t, resp.Code, 200)
}

func TestForgetPassword(t *testing.T) {
	var testAccount = DefaultAccount
	var resp *httptest.ResponseRecorder

	// START - test forget-password endpoint
	// fail test case - not provide the email in the url parameters
	resp = ServeHTTP("POST", "/v1/forget-password", "",
		"application/x-www-form-urlencoded", "")
	assert.Equal(t, resp.Code, 400)

	// success test case
	resp = ServeHTTP("POST", "/v1/forget-password", fmt.Sprintf("{\"email\":\"%v\"}", testAccount),
		"application/json", "")
	assert.Equal(t, resp.Code, 200)
	// END - test forget-password endpoint
}
