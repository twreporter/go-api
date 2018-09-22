package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

type serviceJSON struct {
	ID        int    `json:"ID"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
	Name      string `json:"Name"`
}

type serviceResponse struct {
	Status  string      `json:"status"`
	Service serviceJSON `json:"record"`
}

func TestServiceAuthorization(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var user models.User

	user = getUser(Globs.Defaults.Account)
	email := user.Email
	// set different email
	user.Email = utils.ToNullString("contact@twreporter.org")

	// ===== START - Fail to pass Authorization ===== //
	// Globs.Defaults.Account is not in the admin white list
	resp = serveHTTP("POST", "/v1/services", "", "", fmt.Sprintf("Bearer %v", generateJWT(user)))
	assert.Equal(t, resp.Code, 401)

	// without Authorization header
	resp = serveHTTP("GET", "/v1/services/1", "", "", "")
	assert.Equal(t, resp.Code, 401)
	// ===== END - Fail to pass Authorization ===== //

	// Pass Authorization
	// reset user email
	user.Email = email
	resp = serveHTTP("POST", "/v1/services", `{"name":"test_service"}`,
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(user)))
	assert.Equal(t, resp.Code, 201)
}

func TestCreateAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Wrong JSON POST body
	resp = serveHTTP("POST", "/v1/services", "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 400)

	// Success to create
	resp = serveHTTP("POST", "/v1/services", `{"name":"test_service_1"}`,
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 201)

	// Fail to create the existing service
	resp = serveHTTP("POST", "/v1/services", `{"name":"test_service_1"}`,
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 409)
}

func TestDeleteAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Create a service to delete
	resp = serveHTTP("POST", "/v1/services", `{"name":"test_service_3"}`,
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))

	// Get the new one service's name
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := serviceResponse{}
	json.Unmarshal(body, &res)
	name := res.Service.Name

	// Delete the service successfully
	resp = serveHTTP("DELETE", fmt.Sprintf("/v1/services/%v", name), "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 204)

	// Delete the service again
	resp = serveHTTP("DELETE", fmt.Sprintf("/v1/services/%v", name), "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 404)
}

func TestReadAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Fail to read service due to service not existed
	resp = serveHTTP("GET", "/v1/services/service_not_existed", "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 404)

	// Read successfully
	resp = serveHTTP("GET", "/v1/services/default_service", "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := serviceResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Service.ID, 1)
	assert.Equal(t, res.Service.Name, "default_service")
}

func TestUpdateAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Create a service if service is not existed
	resp = serveHTTP("PUT", "/v1/services/service_to_update", `{"name":"service_to_update"}`,
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 200)

	// Update the existing service
	resp = serveHTTP("PUT", "/v1/services/service_to_update", `{"name":"updated_service"}`,
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 200)

	// Cannot read the old service
	resp = serveHTTP("GET", "/v1/services/service_to_update", "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 404)

	// Read successfully
	resp = serveHTTP("GET", "/v1/services/updated_service", "",
		"application/json", fmt.Sprintf("Bearer %v", generateJWT(getUser(Globs.Defaults.Account))))
	assert.Equal(t, resp.Code, 200)

}
