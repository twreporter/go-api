package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/utils"
)

func TestSendActivation(t *testing.T) {
	const expire int = 100
	var authorization string
	var reqBody = make(map[string]interface{})
	var bodyBytes []byte
	var resp *httptest.ResponseRecorder

	authorization, _ = utils.RetrieveMailServiceAccessToken(expire)
	authorization = fmt.Sprintf("Bearer %s", authorization)

	// successful case
	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		reqBody["email"] = Globs.Defaults.Account
		reqBody["activate_link"] = "test-activate-link"
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		// =====================================
		// Error situation:
		// activate_link is empty
		// =====================================
		reqBody = make(map[string]interface{})
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// email is empty
		// =====================================
		reqBody = map[string]interface{}{
			"activate_link": "test-activate-link",
		}
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// activate_link is not string
		// =====================================
		reqBody = map[string]interface{}{
			"activate_link": 123,
		}
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("StatusCode=StatusInternalServerError", func(t *testing.T) {
		reqBody["email"] = Globs.Defaults.ErrorEmailAddress
		reqBody["activate_link"] = "test-activate-link"
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSendDonationSuccessMail(t *testing.T) {
	const expire int = 100
	var authorization string
	var reqBody = make(map[string]interface{})
	var bodyBytes []byte
	var resp *httptest.ResponseRecorder
	var getDefaultReqBody = func() map[string]interface{} {
		return map[string]interface{}{
			"email":              Globs.Defaults.Account,
			"order_number":       "test-order-number",
			"amount":             300,
			"donation_method":    "信用卡捐款",
			"donation_type":      "定期定額",
			"donation_timestamp": 1541671797,
		}
	}

	authorization, _ = utils.RetrieveMailServiceAccessToken(expire)
	authorization = fmt.Sprintf("Bearer %s", authorization)

	// successful case
	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		reqBody = getDefaultReqBody()
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		// =====================================
		// Error situation:
		// Amount is not number or empty
		// =====================================

		// Copy reqBody from defaultReqBody
		reqBody = getDefaultReqBody()
		reqBody["amount"] = "wrong-amount-type"
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		reqBody["amount"] = 0
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// donation_method is empty
		// =====================================
		reqBody = getDefaultReqBody()
		reqBody["donation_method"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// donation_type is empty
		// =====================================
		reqBody = getDefaultReqBody()
		reqBody["donation_type"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// email is empty
		// =====================================
		reqBody = getDefaultReqBody()
		reqBody["email"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// order_number is empty
		// =====================================
		reqBody = getDefaultReqBody()
		reqBody["order_number"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("StatusCode=StatusInternalServerError", func(t *testing.T) {
		reqBody = getDefaultReqBody()
		reqBody["email"] = Globs.Defaults.ErrorEmailAddress
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", authorization)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
