package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/storage"
)

func TestSignIn(t *testing.T) {
	// TODO: Test for SignInV2
}

func TestActivate(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// START - test activate endpoint //
	user := getReporterAccount(Globs.Defaults.Account)

	// Renew token for v2 endpoint validation
	activateToken := "Activate_Token_2"
	expTime := time.Now().Add(time.Duration(15) * time.Minute)

	as := storage.NewGormStorage(Globs.GormDB)
	if err := as.UpdateReporterAccount(models.ReporterAccount{
		ID:            user.ID,
		ActivateToken: activateToken,
		ActExpTime:    expTime,
	}); nil != err {
		fmt.Println(err.Error())
	}

	// START - test activate endpoint v2//

	// test activate
	resp = serveHTTP("GET", fmt.Sprintf("/v2/auth/activate?email=%v&token=%v", Globs.Defaults.Account, activateToken), "", "", "")
	fmt.Print(resp.Body)

	// validate status code
	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	cookies := resp.Result().Cookies()

	cookieMap := make(map[string]http.Cookie)
	for _, cookie := range cookies {
		cookieMap[cookie.Name] = *cookie
	}
	// validate Set-Cookie header
	assert.Contains(t, cookieMap, "id_token")

	// test activate fails
	resp = serveHTTP("GET", fmt.Sprintf("/v2/auth/activate?email=%v&token=%v", Globs.Defaults.Account, ""), "", "", "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	// END - test activate endpoint v2//

}
