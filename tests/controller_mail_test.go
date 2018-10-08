package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"twreporter.org/go-api/globals"
)

type activationReqBody struct {
	Email        string `json:"email"`
	ActivateLink string `json:"activate_link"`
}

type donationSuccessReqBody struct {
	Address          string `json:"address"`
	Amount           uint   `json:"amount" binding:"required"`
	CardInfoLastFour string `json:"card_info_four_number"`
	CardInfoType     string `json:"card_info_type"`
	Currency         string `json:"currency"`
	DonationDatetime string `json:"donation_datetime"`
	DonationLink     string `json:"donation_link"`
	DonationMethod   string `json:"donation_method" binding:"required"`
	DonationType     string `json:"donation_type" binding:"required"`
	Email            string `json:"email" binding:"required"`
	Name             string `json:"name"`
	NationalID       string `json:"national_id"`
	OrderNumber      string `json:"order_number" binding:"required"`
	PhoneNumber      string `json:"phone_number"`
}

func TestSendActivation(t *testing.T) {
	var reqBody = make(map[string]interface{})
	var bodyBytes []byte
	var resp *httptest.ResponseRecorder

	// successful case
	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		reqBody["email"] = Globs.Defaults.Account
		reqBody["activate_link"] = "test-activate-link"
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		// =====================================
		// Error situation:
		// activate_link is empty
		// =====================================
		reqBody = make(map[string]interface{})
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// email is empty
		// =====================================
		reqBody = map[string]interface{}{
			"activate_link": "test-activate-link",
		}
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// activate_link is not string
		// =====================================
		reqBody = map[string]interface{}{
			"activate_link": 123,
		}
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("StatusCode=StatusInternalServerError", func(t *testing.T) {
		reqBody["email"] = Globs.Defaults.ErrorEmailAddress
		reqBody["activate_link"] = "test-activate-link"
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendActivationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSendDonationSuccessMail(t *testing.T) {
	var reqBody = make(map[string]interface{})
	var bodyBytes []byte
	var resp *httptest.ResponseRecorder
	var defaultReqBody = map[string]interface{}{
		"email":           Globs.Defaults.Account,
		"order_number":    "test-order-number",
		"amount":          300,
		"donation_method": "信用卡捐款",
		"donation_type":   "定期定額",
	}

	// successful case
	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		// =====================================
		// Error situation:
		// Amount is not number or empty
		// =====================================

		// Copy reqBody from defaultReqBody
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		reqBody["amount"] = "wrong-amount-type"
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		reqBody["amount"] = 0
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// donation_method is empty
		// =====================================
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		reqBody["donation_method"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// donation_type is empty
		// =====================================
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		reqBody["donation_type"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// email is empty
		// =====================================
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		reqBody["email"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// =====================================
		// Error situation:
		// order_number is empty
		// =====================================
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		reqBody["order_number"] = ""
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("StatusCode=StatusInternalServerError", func(t *testing.T) {
		for key, value := range defaultReqBody {
			reqBody[key] = value
		}
		reqBody["email"] = Globs.Defaults.ErrorEmailAddress
		bodyBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", fmt.Sprintf("/v1/%s", globals.SendSuccessDonationRoutePath), string(bodyBytes), "application/json", "")
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
