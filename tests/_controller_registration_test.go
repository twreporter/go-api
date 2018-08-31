package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/storage"
)

type UserJSON struct {
}

type RegistrationJSON struct {
	ID            int    `json:"ID"`
	CreatedAt     string `json:"CreatedAt"`
	UpdatedAt     string `json:"UpdatedAt"`
	DeletedAt     string `json:"DeletedAt"`
	ServiceID     string `json:"ServiceID"`
	Service       ServiceJSON
	UserID        int      `json:"UserID"`
	User          UserJSON `json:"User"`
	Active        bool     `json:"Active"`
	ActivateToken string   `json:"ActivateToken"`
}

type RegistrationResponse struct {
	Status        string             `json:"status"`
	Registrations []RegistrationJSON `json:"records"`
	Count         int                `json:"count"`
}

func init() {
	viper.SetDefault("consumersettings.host", "www.twreporter.org")
	viper.SetDefault("consumersettings.protocol", "https")
}

func TestRegisterAndDeregister(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path = "/v1/registrations/default_service"
	var contentType = "application/json"

	// ===== TestRegister START ===== //
	// Wrong JSON POST body
	resp = ServeHTTP("POST", path, "", contentType, "")
	assert.Equal(t, resp.Code, 400)

	// Success to register
	resp = ServeHTTP("POST", "/v1/registrations/default_service", `{"email":"han@twreporter.org"}`, contentType, "")
	assert.Equal(t, resp.Code, 201)

	// Fail to create the existed registration
	resp = ServeHTTP("POST", "/v1/registrations/default_service", `{"email":"han@twreporter.org"}`, contentType, "")
	assert.Equal(t, resp.Code, 409)

	// Fail to create the registration when service is not existed
	resp = ServeHTTP("POST", "/v1/registrations/non_existed_service", `{"email":"han@twreporter.org"}`, contentType, "")
	assert.Equal(t, resp.Code, 404)
	// ===== TestRegister END ===== //

	// ===== TestDeregister START ===== //
	// Success to deregister
	resp = ServeHTTP("DELETE", "/v1/registrations/default_service/han@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 204)

	// Deregister again
	resp = ServeHTTP("DELETE", "/v1/registrations/default_service/han@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 404)

	// It's fine to delete an non-existed resource
	resp = ServeHTTP("DELETE", "/v1/registrations/default_service/not_existed@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 404)

	// Fail to delete the registration when service is not existed
	resp = ServeHTTP("POST", "/v1/registrations/non_existed_service/nickhsine@twreporter.org", "", contentType, "")
	assert.Equal(t, resp.Code, 404)
	// ===== TestDeregister END ===== //
}

func TestGetRegisterUser(t *testing.T) {
	resp := ServeHTTP("GET", fmt.Sprintf("/v1/registrations/default_service/%v", DefaultAccount), "", "", "")
	assert.Equal(t, resp.Code, 200)

	resp = ServeHTTP("GET", "/v1/registrations/default_service/non_existed@twreporter.org", "", "", "")
	assert.Equal(t, resp.Code, 404)
}

func TestGetRegisterUsers(t *testing.T) {
	resp := ServeHTTP("GET", "/v1/registrations/default_service", "", "", "")
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := RegistrationResponse{}
	json.Unmarshal(body, &res)
	assert.NotEqual(t, res.Count, 0)
	assert.Equal(t, res.Registrations[0].Service.Name, DefaultService)

	// Fail
	resp = ServeHTTP("GET", "/v1/registrations/non_existed_service", "", "", "")
	assert.Equal(t, resp.Code, 404)
}

func TestActivateRegistration(t *testing.T) {
	resp := ServeHTTP("GET", fmt.Sprintf("/v1/activation/default_service/%v?activeToken=default_token", DefaultAccount), "", "", "")
	assert.Equal(t, resp.Code, 307)

	loc, _ := resp.Result().Location()
	assert.Equal(t, loc.String(), "https://www.twreporter.org/activate")

	ms := storage.NewGormStorage(DB)
	reg, _ := ms.GetRegistration(DefaultAccount, "default_service")
	assert.Equal(t, reg.Active, true)
}
