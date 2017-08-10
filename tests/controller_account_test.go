package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	var email string
	var resp *httptest.ResponseRecorder

	// START - test signp endpoint //
	// form POST body
	email = "han@twreporter.org"

	resp = ServeHTTP("POST", "/v1/signup", fmt.Sprintf("email=%v&password=0000", email),
		"application/x-www-form-urlencoded", "")
	assert.Equal(t, resp.Code, 201)

	// JSON POST body
	resp = ServeHTTP("POST", "/v1/signup", `{"email":"mika@twreporter.org","password":"0000"}`,
		"application/json", "")
	assert.Equal(t, resp.Code, 201)

	// neither JSON nor form POST body
	resp = ServeHTTP("POST", "/v1/signup", "", "application/text", "")
	assert.Equal(t, resp.Code, 400)

	// signup an already account
	resp = ServeHTTP("POST", "/v1/signup", `{"email":"nickhsine@twreporter.org","password":"0000"}`,
		"application/json", "")
	assert.Equal(t, resp.Code, 409)

	// END - test signup endpoint //

	// START - test activate endpoint //
	as := storage.NewGormStorage(DB)
	user, _ := as.GetReporterAccountData(email)

	// test activate
	activateToken := user.ActivateToken
	resp = ServeHTTP("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", email, activateToken), "", "", "")
	assert.Equal(t, resp.Code, 200)

	// test activate fails
	resp = ServeHTTP("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", "mika@twreporter.org", ""), "", "", "")
	assert.Equal(t, resp.Code, 401)
	// END - test activate endpoint //
}
