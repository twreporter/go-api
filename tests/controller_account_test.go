package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/storage"
)

func TestLogin(t *testing.T) {
	// Login successfully
	req := RequestWithBody("POST", "/v1/login",
		fmt.Sprintf("email=%v&password=%v", DefaultAccount, DefaultPassword))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	// Fail to login
	req = RequestWithBody("POST", "/v1/login",
		fmt.Sprintf("email=%v&password=wrongpassword", DefaultAccount)) //wrong password
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 401)
}

func TestSignupAndActivate(t *testing.T) {
	var email string

	// START - test signp endpoint //
	// form POST body
	email = "han@twreporter.org"
	req := RequestWithBody("POST", "/v1/signup", fmt.Sprintf("email=%v&password=0000", email))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)

	// JSON POST body
	req = RequestWithBody("POST", "/v1/signup", `{"email":"mika@twreporter.org","password":"0000"}`)
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)

	// neither JSON nor form POST body
	req = RequestWithBody("POST", "/v1/signup", "")
	req.Header.Add("Content-Type", "application/text")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400)

	// signup an already account
	req = RequestWithBody("POST", "/v1/signup", `{"email":"nickhsine@twreporter.org","password":"0000"}`)
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 409)

	// END - test signup endpoint //

	// START - test activate endpoint //
	as := storage.NewMembershipStorage(DB)
	user, _ := as.GetReporterAccountData(email)

	// test activate
	activateToken := user.ActivateToken
	req = RequestWithBody("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", email, activateToken), "")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	// test activate fails
	req = RequestWithBody("GET", fmt.Sprintf("/v1/activate?email=%v&token=%v", "mika@twreporter.org", ""), "")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 401)
	// END - test activate endpoint //
}
