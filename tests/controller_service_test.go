package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	// ===== START - Fail to pass Authorization ===== //
	// email(of DefaultID2) is not in the admin white list
	req := RequestWithBody("POST", "/v1/services", "")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID2))))
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 401)

	// without Authorization header
	req, _ = http.NewRequest("GET", "/v1/services/1", nil)
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 401)
	// ===== END - Fail to pass Authorization ===== //

	// Pass Authorization
	req = RequestWithBody("POST", "/v1/services", `{"name":"test_service"}`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)
}

func TestCreateAService(t *testing.T) {
	// Wrong JSON POST body
	req := RequestWithBody("POST", "/v1/services", "")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400)

	// Success to create
	req = RequestWithBody("POST", "/v1/services", `{"name":"test_service_1"}`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)

	// Fail to create the existing service
	req = RequestWithBody("POST", "/v1/services", `{"name":"test_service_1"}`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 409)
}

func TestDeleteAService(t *testing.T) {
	// Create a service to delete
	req := RequestWithBody("POST", "/v1/services", `{"name":"test_service_3"}`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)

	// Get the new one service's name
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := ServiceResponse{}
	json.Unmarshal(body, &res)
	name := res.Service.Name

	// Delete the service successfully
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/services/%v", name), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 204)
}

func TestReadAService(t *testing.T) {
	// Fail to read service due to service not existed
	req, _ := http.NewRequest("GET", "/v1/services/service_not_existed", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 404)

	// Read successfully
	req, _ = http.NewRequest("GET", "/v1/services/default_service", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := ServiceResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Service.ID, 1)
	assert.Equal(t, res.Service.Name, "default_service")
}

func TestUpdateAService(t *testing.T) {
	// Create a service if service is not existed
	req := RequestWithBody("PUT", "/v1/services/service_to_update", `{"name":"service_to_update"}`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	// Update the existing service
	req = RequestWithBody("PUT", "/v1/services/service_to_update", `{"name":"updated_service"}`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	// Cannot read the old service
	req, _ = http.NewRequest("GET", "/v1/services/service_to_update", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 404)

	// Read successfully
	req, _ = http.NewRequest("GET", "/v1/services/updated_service", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	req.Header.Add("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

}
