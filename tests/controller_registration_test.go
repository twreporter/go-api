package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

type userJSON struct {
}

type registrationJSON struct {
	ID            int    `json:"ID"`
	CreatedAt     string `json:"CreatedAt"`
	UpdatedAt     string `json:"UpdatedAt"`
	DeletedAt     string `json:"DeletedAt"`
	ServiceID     string `json:"ServiceID"`
	Service       serviceJSON
	UserID        int      `json:"UserID"`
	User          userJSON `json:"User"`
	Active        bool     `json:"Active"`
	ActivateToken string   `json:"ActivateToken"`
}

type registrationResponse struct {
	Status        string             `json:"status"`
	Registrations []registrationJSON `json:"records"`
	Count         int                `json:"count"`
}

func TestRegisterAndDeregister(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path = "/v1/registrations/default_service"
	var contentType = "application/json"

	// ===== TestRegister START ===== //
	// Wrong JSON POST body
	resp = serveHTTP("POST", path, "", contentType, "")
	assert.Equal(t, resp.Code, 400)

	// Success to register
	resp = serveHTTP("POST", "/v1/registrations/default_service", `{"email":"han@twreporter.org"}`, contentType, "")
	assert.Equal(t, resp.Code, 201)

	// Fail to create the existed registration
	resp = serveHTTP("POST", "/v1/registrations/default_service", `{"email":"han@twreporter.org"}`, contentType, "")
	assert.Equal(t, resp.Code, 409)

	// Fail to create the registration when service is not existed
	resp = serveHTTP("POST", "/v1/registrations/non_existed_service", `{"email":"han@twreporter.org"}`, contentType, "")
	assert.Equal(t, resp.Code, 404)
	// ===== TestRegister END ===== //

	// ===== TestDeregister START ===== //
	// Success to deregister
	resp = serveHTTP("DELETE", "/v1/registrations/default_service/han@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 204)

	// Deregister again
	resp = serveHTTP("DELETE", "/v1/registrations/default_service/han@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 404)

	// It's fine to delete an non-existed resource
	resp = serveHTTP("DELETE", "/v1/registrations/default_service/not_existed@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 404)

	// Fail to delete the registration when service is not existed
	resp = serveHTTP("POST", "/v1/registrations/non_existed_service/nickhsine@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 404)
	// ===== TestDeregister END ===== //
}

func TestGetRegisterUser(t *testing.T) {
	resp := serveHTTP("GET", fmt.Sprintf("/v1/registrations/default_service/%v", Globs.Defaults.Account), "", "", "")
	assert.Equal(t, resp.Code, 200)

	resp = serveHTTP("GET", "/v1/registrations/default_service/non_existed@twreporter.org", "", "", "")
	assert.Equal(t, resp.Code, 404)
}

func TestGetRegisterUsers(t *testing.T) {
	resp := serveHTTP("GET", "/v1/registrations/default_service", "", "", "")
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := registrationResponse{}
	json.Unmarshal(body, &res)
	assert.NotEqual(t, res.Count, 0)
	assert.Equal(t, res.Registrations[0].Service.Name, Globs.Defaults.Service)

	// Fail
	resp = serveHTTP("GET", "/v1/registrations/non_existed_service", "", "", "")
	assert.Equal(t, resp.Code, 404)
}

func TestActivateRegistration(t *testing.T) {
	utils.Cfg.ConsumerSettings.Host = "www.twreporter.org"
	utils.Cfg.ConsumerSettings.Protocol = "https"
	resp := serveHTTP("GET", fmt.Sprintf("/v1/activation/default_service/%v?activeToken=default_token", Globs.Defaults.Account), "", "", "")
	assert.Equal(t, resp.Code, 307)

	loc, _ := resp.Result().Location()
	assert.Equal(t, loc.String(), "https://www.twreporter.org/activate")

	ms := storage.NewGormStorage(Globs.GormDB)
	reg, _ := ms.GetRegistration(Globs.Defaults.Account, "default_service")
	assert.Equal(t, reg.Active, true)
}
