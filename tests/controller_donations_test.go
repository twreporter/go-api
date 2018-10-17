package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/guregu/null.v3"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
)

type (
	linePayResultURL struct {
		FrontendRedirectURL string `json:"frontend_redirect_url"`
		BackendRedirectURL  string `json:"backend_redirect_url"`
	}
	donationRecord struct {
		Amount      uint              `json:"amount"`
		CardInfo    models.CardInfo   `json:"card_info"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		ID          uint              `json:"id"`
		OrderNumber string            `json:"order_number"`
		PayMethod   string            `json:"pay_method"`
		PeriodicID  null.Int          `json:"periodic_id"`
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
		Amount      uint              `json:"amount"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		MerchantID  string            `json:"merchant_id"`
		OrderNumber string            `json:"order_number"`
		Prime       string            `json:"prime"`
		ResultURL   linePayResultURL  `json:"result_url"` // Line pay needed only
		UserID      uint              `json:"user_id"`
	}
)

const (
	testPrime            = "test_3a2fb2b7e892b914a03c95dd4dd5dc7970c908df67a49527c0a648b2bc9"
	testDetails          = "報導者小額捐款"
	testAmount      uint = 500
	testCurrency         = "TWD"
	testOrderNumber      = "otd:developer@twreporter.org:1531966435"
	testMerchantID       = "GlobalTesting_CTBC"
)

var testCardholder = models.Cardholder{
	PhoneNumber: null.StringFrom("+886912345678"),
	Name:        null.StringFrom("王小明"),
	Email:       "developer@twreporter.org",
	ZipCode:     null.StringFrom("104"),
	Address:     null.StringFrom("台北市中山區南京東路X巷X號X樓"),
	NationalID:  null.StringFrom("A123456789"),
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

func testCardholderWithDefaultValue(t *testing.T, ch models.Cardholder) {
	assert.Equal(t, testCardholder.PhoneNumber.ValueOrZero(), ch.PhoneNumber.ValueOrZero())
	assert.Equal(t, testCardholder.Name.ValueOrZero(), ch.Name.ValueOrZero())
	assert.Equal(t, testCardholder.Email, ch.Email)
	assert.Equal(t, testCardholder.ZipCode.ValueOrZero(), ch.ZipCode.ValueOrZero())
	assert.Equal(t, testCardholder.NationalID.ValueOrZero(), ch.NationalID.ValueOrZero())
	assert.Equal(t, testCardholder.Address.ValueOrZero(), ch.Address.ValueOrZero())
}

func testDonationDataValidation(t *testing.T, path string, userID uint, authorization string) {
	var resp *httptest.ResponseRecorder
	var reqBody requestBody
	var reqBodyInBytes []byte

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Lack of UserID
		// ===========================================
		reqBody = requestBody{
			Prime:  testPrime,
			Amount: testAmount,
			Cardholder: models.Cardholder{
				Email: "developer@twreporter.org",
			},
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Lack of Prime
		// ===========================================
		reqBody = requestBody{
			Amount: testAmount,
			Cardholder: models.Cardholder{
				Email: "developer@twreporter.org",
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Lack of Cardholder
		// ===========================================
		reqBody = requestBody{
			Prime:  testPrime,
			Amount: testAmount,
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Lack of Cardholder.Email
		// ===========================================
		reqBody = requestBody{
			Prime:  testPrime,
			Amount: testAmount,
			Cardholder: models.Cardholder{
				Name:        null.StringFrom("王小明"),
				PhoneNumber: null.StringFrom("+886912345678"),
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Lack of Amount
		// ===========================================
		reqBody = requestBody{
			Prime: testPrime,
			Cardholder: models.Cardholder{
				Email: "developer@twreporter.org",
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Malformed Email
		// ===========================================
		reqBody = requestBody{
			Prime:  testPrime,
			Amount: testAmount,
			Cardholder: models.Cardholder{
				Email: "developer-twreporter,org",
			},
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Amount is less than 1(minimum value)
		// ===========================================
		reqBody = requestBody{
			Prime:  testPrime,
			Amount: 0,
			Cardholder: models.Cardholder{
				Email: "developer@twreporter.org",
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		// ===========================================
		// Failure (Data Validation)
		// - Create a Donation by Credit Card
		// - Malformed Cardholder.PhoneNumber (E.164 format)
		// ===========================================
		reqBody = requestBody{
			Prime:  testPrime,
			Amount: 0,
			Cardholder: models.Cardholder{
				Email:       "developer@twreporter.org",
				PhoneNumber: null.StringFrom("0912345678"),
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func testCreateADonationRecord(t *testing.T, path string, userID uint, isPeriodic bool, authorization string) {
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
	t.Run("StatusCode=StatusCreated", func(t *testing.T) {
		reqBody = requestBody{
			Prime:       testPrime,
			Amount:      testAmount,
			Currency:    testCurrency,
			Details:     testDetails,
			OrderNumber: testOrderNumber,
			MerchantID:  testMerchantID,
			Cardholder:  testCardholder,
			UserID:      userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)
		resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
		json.Unmarshal(resBodyInBytes, &resBody)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, "success", resBody.Status)
		assert.Equal(t, testAmount, resBody.Data.Amount)
		assert.Equal(t, isPeriodic, resBody.Data.PeriodicID.Valid)
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
			Cardholder: models.Cardholder{
				Email: "developer@twreporter.org",
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)
		resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
		json.Unmarshal(resBodyInBytes, &resBody)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, "success", resBody.Status)
		assert.Equal(t, testAmount, resBody.Data.Amount)
		assert.Equal(t, isPeriodic, resBody.Data.PeriodicID.Valid)
		assert.Equal(t, testCurrency, resBody.Data.Currency)
		assert.Equal(t, testDetails, resBody.Data.Details)
	})

	// ===========================================
	// Failure (Server Error)
	// - Create a Donation by Credit Card
	// - Invalid Prime
	// ===========================================
	t.Run("StatusCode=StatusInternalServerError", func(t *testing.T) {
		reqBody = requestBody{
			Prime:  "test_prime_which_will_occurs_error",
			Amount: testAmount,
			Cardholder: models.Cardholder{
				Email: "developer@twreporter.org",
			},
			UserID: userID,
		}

		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTP("POST", path, string(reqBodyInBytes), "application/json", authorization)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	// ===========================================
	// Failures (Client Error)
	// - Create a Donation by Credit Card
	// - Request Body Data Validation Error
	// ===========================================
	testDonationDataValidation(t, path, userID, authorization)
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
	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		path = "/v1/donations/credit_card"
		resp = serveHTTP("POST", path, "", "application/json", "")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	// ===========================================
	// Failure (Client Error)
	// - Create a Donation on Unauthorized Resource
	// - 403 Forbidden
	// ===========================================
	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		path = "/v1/donations/credit_card"
		resp = serveHTTP("POST", path, `{"user_id":1000}`, "application/json", authorization)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	// ===========================================
	// Failure (Client Error)
	// - Create a Donation by Unrecognized Pay Method
	// - 404 Not Found Error
	// ===========================================
	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		path = "/v1/donations/unknown_pay_method"

		resp = serveHTTP("POST", path, fmt.Sprintf(`{"user_id":%d}`, user.ID), "application/json", authorization)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	// ==========================================
	// Test One Time Donation Creation
	// =========================================
	path = "/v1/donations/credit_card"
	testCreateADonationRecord(t, path, user.ID, false, authorization)
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
	t.Run("StatusCode=StatusUnauthrorized", func(t *testing.T) {
		path = "/v1/periodic-donations"
		resp = serveHTTP("POST", path, "", "application/json", "")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	// ===========================================
	// Failure (Client Error)
	// - Create a Periodic Donation on Unauthenticated Resource
	// - 403 Forbidden
	// ===========================================
	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		path = "/v1/periodic-donations"
		resp = serveHTTP("POST", path, `{"user_id":1000}`, "application/json", authorization)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	// ==========================================
	// Test Periodic Donation Creation
	// =========================================
	path = "/v1/periodic-donations"
	testCreateADonationRecord(t, path, user.ID, true, authorization)
}

func createDefaultPostBody(userID uint, email string) (reqBody requestBody) {
	reqBody = requestBody{
		Prime:  testPrime,
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Email: email,
		},
		MerchantID: testMerchantID,
		UserID:     userID,
	}

	return reqBody
}

func createDefaultDonationRecord(user models.User, endpoint string) responseBody {
	// create jwt of this user
	jwt := generateJWT(user)
	// prepare jwt authorization string
	authorization := fmt.Sprintf("Bearer %s", jwt)

	// create a donation by HTTP POST request
	reqBody := createDefaultPostBody(user.ID, user.Email.ValueOrZero())
	reqBodyInBytes, _ := json.Marshal(reqBody)
	resp := serveHTTP("POST", endpoint, string(reqBodyInBytes), "application/json", authorization)
	respInBytes, _ := ioutil.ReadAll(resp.Result().Body)
	defer resp.Result().Body.Close()

	// parse response into struct
	resBody := responseBody{}
	json.Unmarshal(respInBytes, &resBody)
	return resBody
}

func createDefaultPeriodicDonationRecord(user models.User) responseBody {
	// create a default periodic donation record
	path := "/v1/periodic-donations"
	return createDefaultDonationRecord(user, path)
}

func createDefaultPrimeDonationRecord(user models.User) responseBody {
	// create a default prime donation record
	path := "/v1/donations/credit_card"
	return createDefaultDonationRecord(user, path)
}

func TestPatchAPeriodicDonation(t *testing.T) {
	// setup before test
	donatorEmail := "prime-donator@twreporter.org"
	// create a new user
	user := createUser(donatorEmail)

	// get record to patch
	res := createDefaultPeriodicDonationRecord(user)

	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		path := fmt.Sprintf("/v1/periodic-donations/%d", res.Data.PeriodicID.ValueOrZero())
		resp := serveHTTP("PATCH", path, "", "application/json", "")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		otherUserID := 100
		path := fmt.Sprintf("/v1/periodic-donations/%d", res.Data.PeriodicID.ValueOrZero())
		resp := serveHTTP("PATCH", path, fmt.Sprintf(`{"user_id": %d}`, otherUserID), "application/json", authorization)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		path := fmt.Sprintf("/v1/periodic-donations/%d", res.Data.PeriodicID.ValueOrZero())
		reqBody := map[string]interface{}{
			"user_id": user.ID,
			// to_feedback should be boolean
			"to_feedback": "true",
			// national_id should be string
			"national_id": true,
		}
		reqBodyInBytes, _ := json.Marshal(reqBody)
		resp := serveHTTP("PATCH", path, string(reqBodyInBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		recordIDNotFound := 1000
		path := fmt.Sprintf("/v1/periodic-donations/%d", recordIDNotFound)
		reqBody := map[string]interface{}{
			"user_id":      user.ID,
			"to_feedback":  true,
			"send_receipt": "no",
		}
		reqBodyInBytes, _ := json.Marshal(reqBody)
		resp := serveHTTP("PATCH", path, string(reqBodyInBytes), "application/json", authorization)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		path := fmt.Sprintf("/v1/periodic-donations/%d", res.Data.PeriodicID.ValueOrZero())
		reqBody := map[string]interface{}{
			"user_id":      user.ID,
			"to_feedback":  true,
			"send_receipt": "no",
		}
		reqBodyInBytes, _ := json.Marshal(reqBody)
		resp := serveHTTP("PATCH", path, string(reqBodyInBytes), "application/json", authorization)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})
}

func TestPatchAPrimeDonation(t *testing.T) {
	// setup before test
	donatorEmail := "periodic-donator@twreporter.org"
	// create a new user
	user := createUser(donatorEmail)

	// get record to patch
	res := createDefaultPrimeDonationRecord(user)

	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		path := fmt.Sprintf("/v1/donations/prime/%d", res.Data.ID)
		resp := serveHTTP("PATCH", path, "", "application/json", "")
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		otherUserID := 100
		path := fmt.Sprintf("/v1/donations/prime/%d", res.Data.ID)
		resp := serveHTTP("PATCH", path, fmt.Sprintf(`{"user_id": %d}`, otherUserID), "application/json", authorization)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("StatusCode=StatusBadRequest", func(t *testing.T) {
		path := fmt.Sprintf("/v1/donations/prime/%d", res.Data.ID)
		reqBody := map[string]interface{}{
			"user_id": user.ID,
			// national_id should be string
			"national_id": true,
		}
		reqBodyInBytes, _ := json.Marshal(reqBody)
		resp := serveHTTP("PATCH", path, string(reqBodyInBytes), "application/json", authorization)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		recordIDNotFound := 1000
		path := fmt.Sprintf("/v1/donations/prime/%d", recordIDNotFound)
		reqBody := map[string]interface{}{
			"user_id":      user.ID,
			"send_receipt": "no",
		}
		reqBodyInBytes, _ := json.Marshal(reqBody)
		resp := serveHTTP("PATCH", path, string(reqBodyInBytes), "application/json", authorization)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		path := fmt.Sprintf("/v1/donations/prime/%d", res.Data.ID)
		reqBody := map[string]interface{}{
			"user_id":      user.ID,
			"send_receipt": "no",
		}
		reqBodyInBytes, _ := json.Marshal(reqBody)
		resp := serveHTTP("PATCH", path, string(reqBodyInBytes), "application/json", authorization)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})
}

func TestGetAPrimeDonationOfAUser(t *testing.T) {
	// setup before test
	donatorEmail := "get-prime-donator@twreporter.org"
	// create a new user
	user := createUser(donatorEmail)

	primeRes := createDefaultPrimeDonationRecord(user)

	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	idToken := generateIDToken(user)
	cookie := http.Cookie{
		Name:     "id_token",
		Value:    idToken,
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
	}

	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		path := fmt.Sprintf("/v1/donations/prime/%d?user_id=%d", primeRes.Data.ID, user.ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", "", cookie)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		path := fmt.Sprintf("/v1/donations/prime/%d?user_id=%d", primeRes.Data.ID, getUser(Globs.Defaults.Account).ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		recordIDNotFound := 1000
		path := fmt.Sprintf("/v1/donations/prime/%d?user_id=%d", recordIDNotFound, user.ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("StatusCode=StatusOK", func(t *testing.T) {
		path := fmt.Sprintf("/v1/donations/prime/%d?user_id=%d", primeRes.Data.ID, user.ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		respInBytes, _ := ioutil.ReadAll(resp.Result().Body)
		defer resp.Result().Body.Close()

		// parse response into struct
		resBody := responseBody{}
		json.Unmarshal(respInBytes, &resBody)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, testAmount, resBody.Data.Amount)
		assert.Equal(t, donatorEmail, resBody.Data.Cardholder.Email)
		assert.Equal(t, "4242", resBody.Data.CardInfo.LastFour.ValueOrZero())
		assert.Equal(t, "424242", resBody.Data.CardInfo.BinCode.ValueOrZero())
		assert.Equal(t, int64(0), resBody.Data.CardInfo.Funding.ValueOrZero())
		assert.Equal(t, testCurrency, resBody.Data.Currency)
		assert.Equal(t, "credit_card", resBody.Data.PayMethod)
		assert.NotEmpty(t, resBody.Data.OrderNumber)
	})
}

func TestGetAPeriodicDonationOfAUser(t *testing.T) {
	// setup before test
	donatorEmail := "get-periodic-donator@twreporter.org"
	// create a new user
	user := createUser(donatorEmail)

	tokenRes := createDefaultPeriodicDonationRecord(user)

	jwt := generateJWT(user)
	authorization := fmt.Sprintf("Bearer %s", jwt)

	idToken := generateIDToken(user)
	cookie := http.Cookie{
		Name:     "id_token",
		Value:    idToken,
		MaxAge:   3600,
		Secure:   false,
		HttpOnly: true,
	}

	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		path := fmt.Sprintf("/v1/periodic-donations/%d?user_id=%d", tokenRes.Data.ID, user.ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", "", cookie)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		path := fmt.Sprintf("/v1/periodic-donations/%d?user_id=%d", tokenRes.Data.ID, getUser(Globs.Defaults.Account).ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		recordIDNotFound := 1000
		path := fmt.Sprintf("/v1/periodic-donations/%d?user_id=%d", recordIDNotFound, user.ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("StatusCode=StatusOK", func(t *testing.T) {
		path := fmt.Sprintf("/v1/periodic-donations/%d?user_id=%d", tokenRes.Data.ID, user.ID)
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusOK, resp.Code)
		respInBytes, _ := ioutil.ReadAll(resp.Result().Body)
		defer resp.Result().Body.Close()

		// parse response into struct
		resBody := responseBody{}
		json.Unmarshal(respInBytes, &resBody)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, testAmount, resBody.Data.Amount)
		assert.Equal(t, donatorEmail, resBody.Data.Cardholder.Email)
		assert.Equal(t, testCurrency, resBody.Data.Currency)
		assert.Empty(t, resBody.Data.OrderNumber)
	})
}

/* GetDonationsOfAUser is not implemented yet
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
		ID:    1000,
		Email: null.StringFrom("unknown@twreporter.org"),
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
*/
