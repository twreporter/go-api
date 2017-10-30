package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	//"twreporter.org/go-api/utils"
)

type ServiceJSON struct {
	ID        int    `json:"ID"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
	Name      string `json:"Name"`
}

type ServiceResponse struct {
	Status  string      `json:"status"`
	Service ServiceJSON `json:"record"`
}

func TestServiceAuthorization(t *testing.T) {
	const userID = "2"
	var resp *httptest.ResponseRecorder

	// ===== START - Fail to pass Authorization ===== //
	// email(of userID) is not in the admin white list
	resp = ServeHTTP("POST", "/v1/services", "", "", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(userID))))
	assert.Equal(t, resp.Code, 401)

	// without Authorization header
	resp = ServeHTTP("GET", "/v1/services/1", "", "", "")
	assert.Equal(t, resp.Code, 401)
	// ===== END - Fail to pass Authorization ===== //

	// Pass Authorization
	resp = ServeHTTP("POST", "/v1/services", `{"name":"test_service"}`,
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 201)
}

func TestCreateAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Wrong JSON POST body
	resp = ServeHTTP("POST", "/v1/services", "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 400)

	// Success to create
	resp = ServeHTTP("POST", "/v1/services", `{"name":"test_service_1"}`,
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 201)

	// Fail to create the existing service
	resp = ServeHTTP("POST", "/v1/services", `{"name":"test_service_1"}`,
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 409)
}

func TestDeleteAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Create a service to delete
	resp = ServeHTTP("POST", "/v1/services", `{"name":"test_service_3"}`,
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))

	// Get the new one service's name
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := ServiceResponse{}
	json.Unmarshal(body, &res)
	name := res.Service.Name

	// Delete the service successfully
	resp = ServeHTTP("DELETE", fmt.Sprintf("/v1/services/%v", name), "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 204)

	// Delete the service again
	resp = ServeHTTP("DELETE", fmt.Sprintf("/v1/services/%v", name), "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 404)
}

func TestReadAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Fail to read service due to service not existed
	resp = ServeHTTP("GET", "/v1/services/service_not_existed", "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 404)

	// Read successfully
	resp = ServeHTTP("GET", "/v1/services/default_service", "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := ServiceResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Service.ID, 1)
	assert.Equal(t, res.Service.Name, "default_service")
}

func TestUpdateAService(t *testing.T) {
	var resp *httptest.ResponseRecorder

	// Create a service if service is not existed
	resp = ServeHTTP("PUT", "/v1/services/service_to_update", `{"name":"service_to_update"}`,
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 200)

	// Update the existing service
	resp = ServeHTTP("PUT", "/v1/services/service_to_update", `{"name":"updated_service"}`,
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 200)

	// Cannot read the old service
	resp = ServeHTTP("GET", "/v1/services/service_to_update", "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 404)

	// Read successfully
	resp = ServeHTTP("GET", "/v1/services/updated_service", "",
		"application/json", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 200)

}
