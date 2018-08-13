package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
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
		CountryCode string `json:"country_code"`
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
	linePayResultURL struct {
		FrontendRedirectURL string `json:"frontend_redirect_url"`
		BackendRedirectURL  string `json:"backend_redirect_url"`
	}
	donationRecord struct {
		ID          uint       `json:"id"`
		IsPeriodic  bool       `json:"is_periodic"`
		PayMethod   string     `json:"pay_method"`
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
	responseBodyForList struct {
		Status string `json:"status"`
		Data   struct {
			Records []donationRecord `json:"records"`
			Meta    struct {
				Total  uint `json:"total"`
				Offset uint `json:"offset"`
				Limit  uint `json:"limit"`
			}
		} `json:"data"`
	}
	requestBody struct {
		Prime       string           `json:"prime"`
		Amount      uint             `json:"amount"`
		Currency    string           `json:"currency"`
		Details     string           `json:"details"`
		Cardholder  cardholder       `json:"cardholder"`
		OrderNumber string           `json:"order_number"`
		MerchantID  string           `json:"merchant_id"`
		ResultURL   linePayResultURL `json:"result_url"` // Line pay needed only
	}
)

const (
	testPrime            = "test_3a2fb2b7e892b914a03c95dd4dd5dc7970c908df67a49527c0a648b2bc9"
	testDetails          = "報導者小額捐款"
	testAmount      uint = 500
	testCurrency         = "TWD"
	testOrderNumber      = "otd:developer@twreporter.org:1531966435"
	testMerchantID       = "GlobalTesting_CTBC"
	donatorEmail         = "donater@twreporter.org"
)

var testCardholder = cardholder{
	PhoneNumber: "+886912345678",
	Name:        "王小明",
	Email:       "developer@twreporter.org",
	ZipCode:     "104",
	Address:     "台北市中山區南京東路X巷X號X樓",
	NationalID:  "A123456789",
}

var defaults = struct {
	Total      uint
	Offset     uint
	Limit      uint
	CreditCard string
}{
	Total:      0,
	Offset:     0,
	Limit:      10,
	CreditCard: "credit_card",
}

func setUpBeforeDonationsTest() {
	user := createUser(donatorEmail)

	reqBody := requestBody{
		Prime:  testPrime,
		Amount: testAmount,
		Cardholder: cardholder{
			Email: donatorEmail,
		},
	}

	// insert three default records
	// first, insert a periodic donation
	path := fmt.Sprintf("/v1/users/%d/periodic_donations", user.ID)
	reqBodyInBytes, _ := json.Marshal(reqBody)
	serveHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	// second, insert a one time donation
	path = fmt.Sprintf("/v1/users/%d/donations/credit_card", user.ID)
	serveHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	// third, insert a periodic donation
	path = fmt.Sprintf("/v1/users/%d/periodic_donations", user.ID)
	serveHTTP("POST", path, string(reqBodyInBytes), "application/json", "")

	// set default total to 3
	defaults.Total = 3
}

func testCardholderWithDefaultValue(t *testing.T, ch cardholder) {
	assert.Equal(t, testCardholder.PhoneNumber, ch.PhoneNumber)
	assert.Equal(t, testCardholder.Name, ch.Name)
	assert.Equal(t, testCardholder.Email, ch.Email)
	assert.Equal(t, testCardholder.ZipCode, ch.ZipCode)
	assert.Equal(t, testCardholder.NationalID, ch.NationalID)
	assert.Equal(t, testCardholder.Address, ch.Address)
}

func testDonationDataValidation(t *testing.T, path string, authorization string) {
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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

	assert.Equal(t, 400, resp.Code)
}

func testCreateADonationRecord(t *testing.T, path string, isPeriodic bool, authorization string) {
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
		MerchantID:  testMerchantID,
		Cardholder:  testCardholder,
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)

	assert.Equal(t, 201, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, testAmount, resBody.Data.Amount)
	assert.Equal(t, isPeriodic, resBody.Data.IsPeriodic)
	assert.Equal(t, testCurrency, resBody.Data.Currency)
	assert.Equal(t, testDetails, resBody.Data.Details)
	assert.Equal(t, testOrderNumber, resBody.Data.OrderNumber)
	testCardholderWithDefaultValue(t, resBody.Data.Cardholder)

	// ===========================================
	// Success
	// - Create a Donation by Credit Card
	// - Provide minimun required fields
	// ===========================================
	reqBody = requestBody{
		Prime:      testPrime,
		Amount:     testAmount,
		MerchantID: testMerchantID,
		Cardholder: cardholder{
			Email: "developer@twreporter.org",
		},
	}

	reqBodyInBytes, _ = json.Marshal(reqBody)
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)
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
	resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

	assert.Equal(t, 500, resp.Code)

	// ===========================================
	// Failures (Client Error)
	// - Create a Donation by Credit Card
	// - Request Body Data Validation Error
	// ===========================================
	testDonationDataValidation(t, path, authorization)
}

func TestCreateADonation(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path string

	user := getUser(Globs.Defaults.Account)
	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	// ===========================================
	// Failure (Client Error)
	// - Create a Donation Without Authorization Header
	// - 401 Unauthorized
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/donations/credit_card", user.ID)
	resp = serveHTTP("POST", path, "", "application/json", "")
	assert.Equal(t, 401, resp.Code)

	// ===========================================
	// Failure (Client Error)
	// - Create a Donation on Unauthenticated Resource
	// - 403 Forbidden
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/donations/credit_card", 2)
	resp = serveHTTP("POST", path, "", "application/json", authorization)
	assert.Equal(t, 403, resp.Code)

	// ===========================================
	// Failure (Client Error)
	// - Create a Donation by Unrecognized Pay Method
	// - 404 Not Found Error
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/donations/unknown_pay_method", user.ID)

	resp = serveHTTP("POST", path, "", "application/json", authorization)

	assert.Equal(t, 404, resp.Code)

	// ==========================================
	// Test One Time Donation Creation
	// =========================================
	path = fmt.Sprintf("/v1/users/%d/donations/credit_card", user.ID)
	testCreateADonationRecord(t, path, false, authorization)
}

func TestCreateAPeriodicDonation(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path string

	user := getUser(Globs.Defaults.Account)
	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	// ===========================================
	// Failure (Client Error)
	// - Create a Periodic Donation Without Authorization Header
	// - 401 Unauthorized
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/periodic_donations", user.ID)
	resp = serveHTTP("POST", path, "", "application/json", "")
	assert.Equal(t, 401, resp.Code)

	// ===========================================
	// Failure (Client Error)
	// - Create a Periodic Donation on Unauthenticated Resource
	// - 403 Forbidden
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/periodic_donations", 2)
	resp = serveHTTP("POST", path, "", "application/json", authorization)
	assert.Equal(t, 403, resp.Code)

	// ==========================================
	// Test Periodic Donation Creation
	// =========================================
	path = "/v1/users/1/periodic_donations"
	testCreateADonationRecord(t, path, true, authorization)
}

func TestGetDonations(t *testing.T) {
	var resBody responseBodyForList
	var resp *httptest.ResponseRecorder
	var path string

	// set up default records
	setUpBeforeDonationsTest()

	defaultUser := getUser(Globs.Defaults.Account)
	user := getUser(donatorEmail)
	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	// ===========================================
	// Failure (Client Error)
	// - Get Donations of A Unkonwn User Without Authorization Header
	// - 401 Unauthorized
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/donations", user.ID)
	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, 401, resp.Code)

	// ===========================================
	// Failure (Client Error)
	// - Create a Periodic Donation on Unauthenticated Resource
	// - 403 Forbidden
	// ===========================================
	path = fmt.Sprintf("/v1/users/%d/donations", defaultUser.ID)
	resp = serveHTTP("GET", path, "", "", authorization)
	assert.Equal(t, 403, resp.Code)

	// ===========================================
	// Failure (Client Error)
	// - Get Donations of A Unkonwn User
	// - 404 Not Found Error
	// ===========================================
	path = "/v1/users/1000/donations"
	jwt = generateJWT(models.User{
		ID: 1000,
		Email: sql.NullString{
			String: "unknown@twreporter.org",
			Valid:  true,
		},
	})

	resp = serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %s", jwt))
	assert.Equal(t, 404, resp.Code)

	// ================================================================
	// Success
	// - Get Donations of A User Without `pay_methods` Param
	// - Missing `pay_methods` Param (which means all pay_methods)
	// ================================================================
	path = fmt.Sprintf("/v1/users/%d/donations?offset=%d&limit=%d", user.ID, defaults.Offset, defaults.Limit)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ := ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, defaults.Total, resBody.Data.Meta.Total)
	assert.Equal(t, defaults.Offset, resBody.Data.Meta.Offset)
	assert.Equal(t, defaults.Limit, resBody.Data.Meta.Limit)
	assert.Equal(t, defaults.Total, len(resBody.Data.Records))
	assert.Equal(t, defaults.CreditCard, resBody.Data.Records[0].PayMethod)
	assert.Equal(t, true, resBody.Data.Records[0].IsPeriodic)
	assert.Equal(t, defaults.CreditCard, resBody.Data.Records[1].PayMethod)
	assert.Equal(t, false, resBody.Data.Records[1].IsPeriodic)
	assert.Equal(t, defaults.CreditCard, resBody.Data.Records[2].PayMethod)
	assert.Equal(t, true, resBody.Data.Records[2].IsPeriodic)

	// ===================================================
	// Success
	// - Get Donations of A User Without `offset` Param
	// - Missing `offset` Param (which means offset=0)
	// ===================================================
	path = fmt.Sprintf("/v1/users/%d/donations?pay_methods=credit_card&limit=%d", user.ID, defaults.Limit)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, defaults.Total, resBody.Data.Meta.Total)
	assert.Equal(t, defaults.Offset, resBody.Data.Meta.Offset)

	// =====================================================
	// Success
	// - Get Donations of A User Without `limit` Param
	// - Missing `limit` Param (which means limit=10)
	// =====================================================
	path = fmt.Sprintf("/v1/users/%d/donations?pay_methods=credit_card&offset=%d", user.ID, defaults.Offset)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, defaults.Total, resBody.Data.Meta.Total)
	assert.Equal(t, defaults.Offset, resBody.Data.Meta.Limit)

	// ===================================================
	// Success
	// - Get Donations of A User Without Any Params
	// - Missing `pay_method`, `offset` and `limit` Param
	// ===================================================
	path = fmt.Sprintf("/v1/users/%d/donations", user.ID)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, defaults.Total, resBody.Data.Meta.Total)
	assert.Equal(t, defaults.Limit, resBody.Data.Meta.Limit)
	assert.Equal(t, defaults.Total, len(resBody.Data.Records))

	// ===============================================================
	// Success
	// - Get Donations of A User by Providing `offset=1` and `limit=1`
	// - ?offset=1&limit=1
	// ===============================================================
	path = fmt.Sprintf("/v1/users/%d/donations?offset=1&limit=1", user.ID)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, defaults.Total, resBody.Data.Meta.Total)
	assert.Equal(t, 1, resBody.Data.Meta.Limit)
	assert.Equal(t, 1, len(resBody.Data.Records))
	assert.Equal(t, defaults.CreditCard, resBody.Data.Records[0].PayMethod)
	assert.Equal(t, false, resBody.Data.Records[0].IsPeriodic)

	// ====================================================
	// Success
	// - Get Donations of A User With `offset>total`
	// - ?offset=3&limit=1 (offset is equal to or more than total)
	// ====================================================
	path = fmt.Sprintf("/v1/users/%d/donations?offset=%d&limit=1", user.ID, defaults.Total)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, 0, len(resBody.Data.Records))

	// =========================================================
	// Success
	// - Get Donations of A User With `limit=0`
	// - ?offset=0&limit=0 (limit is 0)
	// =========================================================
	path = fmt.Sprintf("/v1/users/%d/donations?offset=0&limit=0", user.ID)
	resp = serveHTTP("GET", path, "", "", "")
	resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, 0, len(resBody.Data.Records))

	// =========================================================
	// Success
	// - Get Donations of A User
	// - Test offset and limit are not unsigned integer
	// - Test SQL Injection, put statement in pay_methods
	// - ?limit=NaN&offset=-1&pay_methods=;select * from users;
	// =========================================================
	path = fmt.Sprintf("/v1/users/%d/donations?limit=NaN&offset=-1&pay_methods=;select * from users;", user.ID)
	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, 3, len(resBody.Data.Records))
}
