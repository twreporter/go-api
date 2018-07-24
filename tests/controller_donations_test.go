package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	cardInfo struct {
		BinCode     string `json:"bin_code"`
		LastFour    string `json:"last_four"`
		Issuer      string `json:"issuer"`
		Funding     int    `json:"funding"`
		Type        int    `json:"type"`
		Level       string `json:"level"`
		Country     string `json:"country"`
		CountryCode string `json:"us"`
		ExpiryDate  string `json:"expiry_date"`
	}
	cardholder struct {
		PhoneNumber string `json:"phone_number"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		ZipCode     string `json:"zip_code"`
		Address     string `json:"address"`
		NationalID  string `json:"national_id"`
	}
	donationRecord struct {
		ID          uint       `json:"id"`
		IsPeriodic  bool       `json:"is_periodic"`
		CardInfo    cardInfo   `json:"card_info"`
		Cardholder  cardholder `json:"cardholder"`
		Amount      uint       `json:"amount"`
		Currency    string     `json:"currency"`
		Details     string     `json:"details"`
		OrderNumber string     `json:"order_number"`
	}
	responseBody struct {
		Status string         `json:"status"`
		Data   donationRecord `json:"data"`
	}
	requestBody struct {
		Prime       string     `json:"prime"`
		Amount      uint       `json:"amount"`
		Currency    string     `json:"currency"`
		Details     string     `json:"details"`
		Cardholder  cardholder `json:"cardholder"`
		OrderNumber string     `json:"order_number"`
		MerchantID  string     `json:"merchat_id"`
		ResultURL   string     `json:"result_url"` // Line pay needed only
	}
)

const (
	testPrime            = "test_3a2fb2b7e892b914a03c95dd4dd5dc7970c908df67a49527c0a648b2bc9"
	testDetails          = "報導者小額捐款"
	testAmount      uint = 500
	testCurrency         = "TWD"
	testOrderNumber      = "otd:developer@twreporter.org:1531966435"
	testMerchatID        = "twreporter_CTBC"
)

var testCardholder = cardholder{
	PhoneNumber: "+886912345678",
	Name:        "王小明",
	Email:       "developer@twreporter.org",
	ZipCode:     "104",
	Address:     "台北市中山區南京東路X巷X號X樓",
	NationalID:  "A123456789",
}

func testCardInfoWithDefaultValue(t *testing.T, ci cardInfo) {
	assert.Equal(t, "424242", ci.BinCode)
	assert.Equal(t, "4242", ci.LastFour)
	assert.Equal(t, "JPMORGAN CHASE BANK NA", ci.Issuer)
	assert.Equal(t, 0, ci.Funding)
	assert.Equal(t, 1, ci.Type)
	assert.Equal(t, "", ci.Level)
	assert.Equal(t, "UNITED STATES", ci.Country)
	assert.Equal(t, "US", ci.CountryCode)
	assert.Equal(t, "202301", ci.ExpiryDate)
}

func testCardholderWithDefaultValue(t *testing.T, ch cardholder) {
	assert.Equal(t, testCardholder.PhoneNumber, ch.PhoneNumber)
	assert.Equal(t, testCardholder.Name, ch.Name)
	assert.Equal(t, testCardholder.Email, ch.Email)
	assert.Equal(t, testCardholder.ZipCode, ch.ZipCode)
	assert.Equal(t, testCardholder.NationalID, ch.NationalID)
	assert.Equal(t, testCardholder.Address, ch.Address)
}

func testDonationDataValidation(t *testing.T, path string) {
	var resp *httptest.ResponseRecorder
	var reqBody requestBody
	var reqBodyInBytes []byte

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Lack of Prime
	// ===========================================
	reqBody = requestBody{
		Amount: testAmount,
		Cardholder: cardholder{
			Email: "developer@twreporter.org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Lack of Cardholder
	// ===========================================
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: testAmount,
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Lack of Cardholder.Email
	// ===========================================
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: testAmount,
		Cardholder: cardholder{
			Name:        "王小明",
			PhoneNumber: "+886912345678",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Lack of Amount
	// ===========================================
	reqBody = requestBody{
		Prime: testPrime,
		Cardholder: cardholder{
			Email: "developer@twreporter.org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Malformed Email
	// ===========================================
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: testAmount,
		Cardholder: cardholder{
			Email: "developer-twreporter,org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Amount is less than 1(minimum value)
	// ===========================================
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: 0,
		Cardholder: cardholder{
			Email: "developer@twreporter.org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)

	// ===========================================
	// Failure (Data Validation)
	// - Create a Donation by Credit Card
	// - Malformed Cardholder.PhoneNumber (E.164 format)
	// ===========================================
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: 0,
		Cardholder: cardholder{
			Email:       "developer@twreporter.org",
			PhoneNumber: "0912345678",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 400, resp.Code)
}

func testCreateADonationRecord(t *testing.T, path string, isPeriodic bool) {
	var resp *httptest.ResponseRecorder
	var reqBody requestBody
	var resBody responseBody
	var reqBodyInBytes []byte
	var resBodyInBytes []byte

	// ===========================================
	// Success
	// - Create a Donation by Credit Card
	// - Provide all the fields except `result_url`
	//   in request body
	// ===========================================

	reqBody = requestBody{
		Prime:       testPrime,
		Amount:      testAmount,
		Currency:    testCurrency,
		Details:     testDetails,
		OrderNumber: testOrderNumber,
		MerchantID:  testMerchatID,
		Cardholder:  testCardholder,
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)

	assert.Equal(t, 201, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, testAmount, resBody.Data.Amount)
	assert.Equal(t, isPeriodic, resBody.Data.IsPeriodic)
	assert.Equal(t, testCurrency, resBody.Data.Currency)
	assert.Equal(t, testDetails, resBody.Data.Details)
	assert.Equal(t, testOrderNumber, resBody.Data.OrderNumber)
	testCardInfoWithDefaultValue(t, resBody.Data.CardInfo)
	testCardholderWithDefaultValue(t, resBody.Data.Cardholder)

	// ===========================================
	// Success
	// - Create a Donation by Credit Card
	// - Provide minimun required fields
	// ===========================================
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: testAmount,
		Cardholder: cardholder{
			Email: "developer@twreporter.org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)

	assert.Equal(t, 201, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, testAmount, resBody.Data.Amount)
	assert.Equal(t, isPeriodic, resBody.Data.IsPeriodic)
	assert.Equal(t, testCurrency, resBody.Data.Currency)
	assert.Equal(t, testDetails, resBody.Data.Details)

	// ===========================================
	// Failure (Server Error)
	// - Create a Donation by Credit Card
	// - Invalid Prime
	// ===========================================
	reqBody = requestBody{
		Prime:  "test_prime_which_will_occurs_error",
		Amount: testAmount,
		Cardholder: cardholder{
			Email: "developer@twreporter.org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = ServeHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	assert.Equal(t, 500, resp.Code)

	// ===========================================
	// Failures (Client Error)
	// - Create a Donation by Credit Card
	// - Request Body Data Validation Error
	// ===========================================
	testDonationDataValidation(t, path)

}

func TestCreateADonation(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path string

	// ===========================================
	// Failure (Client Error)
	// - Create a Donation by Unrecognized Pay Method
	// - 404 Not Found Error
	// ===========================================
	path = "/v1/users/1/donations/unknown_pay_method"

	resp = ServeHTTP("POST", path, "", "application/json", "")

	assert.Equal(t, 404, resp.Code)

	// ==========================================
	// Test One Time Donation Creation
	// =========================================
	path = "/v1/users/1/donations/credit_card"
	testCreateADonationRecord(t, path, false)
}

func TestCreateAPeriodicDonation(t *testing.T) {
	var path string

	// ==========================================
	// Test Periodic Donation Creation
	// =========================================
	path = "/v1/users/1/periodic_donations"
	testCreateADonationRecord(t, path, true)
}

func TestGetDonations(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path string

	// ===========================================
	// Failure (Client Error)
	// - Get Donations of A Unkonwn User
	// - 404 Not Found Error
	// ===========================================
	path = "/v1/users/unknown_user/donations"

	resp = ServeHTTP("GET", path, "", "", "")

	assert.Equal(t, 404, resp.Code)
}
