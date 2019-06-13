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
		Frequency   string            `json:"frequency"`
		ID          uint              `json:"id"`
		Notes       string            `json:"notes"`
		OrderNumber string            `json:"order_number"`
		PayMethod   string            `json:"pay_method"`
		SendReceipt string            `json:"send_receipt"`
		ToFeedback  bool              `json:"to_feedback"`
		IsAnonymous bool              `json:"is_anonymous"`
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
		Cardholder  models.Cardholder `json:"donor"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		Frequency   string            `json:"frequency"`
		MerchantID  string            `json:"merchant_id"`
		PayMethod   string            `json:"pay_method"`
		Prime       string            `json:"prime"`
		UserID      uint              `json:"user_id"`
		ToFeedback  bool              `json:"to_feedback"`
		IsAnonymous bool              `json:"is_anonymous"`
	}
)

const (
	testPrime           = "test_3a2fb2b7e892b914a03c95dd4dd5dc7970c908df67a49527c0a648b2bc9"
	testDetails         = "報導者小額捐款"
	testAmount     uint = 500
	testCurrency        = "TWD"
	testMerchantID      = "GlobalTesting_CTBC"
	testFeedback        = true

	testName        = "報導者測試者"
	testAddress     = "台北市南京東路一段32巷100號10樓"
	testNationalID  = "A12345678"
	testPhoneNumber = "+886912345678"
	testZipCode     = "101"

	monthlyFrequency = "monthly"
	yearlyFrequency  = "yearly"
	oneTimeFrequency = "one_time"

	creditCardPayMethod = "credit_card"
	linePayMethod       = "line"

	defaultPeriodicDetails = "一般線上定期定額捐款"
	defaultOneTimeDetails  = "一般線上單筆捐款"
	defaultCurrency        = "TWD"
)

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

func helperSetupAuth(user models.User) (authorization string, cookie http.Cookie) {
	jwt := generateJWT(user)
	authorization = fmt.Sprintf("Bearer %s", jwt)

	idToken := generateIDToken(user)
	cookie = http.Cookie{
		HttpOnly: true,
		MaxAge:   3600,
		Name:     "id_token",
		Secure:   false,
		Value:    idToken,
	}

	return
}

func testDonationCreateServerError(t *testing.T, path string, userID uint, frequency string, paymethod string, authorization string, cookie http.Cookie) {
	var resp *httptest.ResponseRecorder
	var reqBodyInBytes []byte

	cases := []struct {
		name       string
		reqBody    *requestBody
		resultCode int
	}{
		{
			name: "StatusCode=StatusInternalServerError,Invalid Prime",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				Frequency: frequency,
				Prime:     "test_prime_which_will_occurs_error",
				UserID:    userID,
			},
			resultCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if nil != c.reqBody {
				switch frequency {
				case oneTimeFrequency:
					c.reqBody.PayMethod = paymethod
				case monthlyFrequency, yearlyFrequency:
					c.reqBody.Frequency = frequency
				}
				reqBodyInBytes, _ = json.Marshal(c.reqBody)
			}
			resp = serveHTTPWithCookies("POST", path, string(reqBodyInBytes), "application/json", authorization, cookie)

			assert.Equal(t, c.resultCode, resp.Code)
		})
	}

}

func testDonationCreateClientError(t *testing.T, path string, userID uint, frequency string, paymethod string, authorization string, cookie http.Cookie) {
	var resp *httptest.ResponseRecorder
	var reqBodyInBytes []byte

	cases := []struct {
		name          string
		reqBody       *requestBody
		cookie        *http.Cookie
		authorization string
		resultCode    int
	}{
		{
			name: "StatusCode=StatusBadRequest,Lack of UserID",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of Prime",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of Cardholder",
			reqBody: &requestBody{
				Amount:    testAmount,
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of Cardholder.Email",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Name:        null.StringFrom("王小明"),
					PhoneNumber: null.StringFrom("+886912345678"),
				},
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of Amount",
			reqBody: &requestBody{
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of PayMethod",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				Prime:  testPrime,
				UserID: userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Malformed email",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer-twreporter,org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Amount less than 1",
			reqBody: &requestBody{
				Amount: 0,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Malformed Cardholder.PhoneNumber (E.164 format)",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email:       "developer-twreporter,org",
					PhoneNumber: null.StringFrom("0912345678"),
				},
				PayMethod: creditCardPayMethod,
				Prime:     testPrime,
				UserID:    userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusBadRequest,
		},
		{
			name:       "StatusCode=StatusUnauthorized,Without Cookie",
			resultCode: http.StatusUnauthorized,
		},
		{
			name:       "StatusCode=StatusUnauthorized,Without Authorization Header",
			cookie:     &cookie,
			resultCode: http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusForbidden,Unauthorized Resource",
			reqBody: &requestBody{
				UserID: 1000,
			},
			cookie:        &cookie,
			authorization: authorization,
			resultCode:    http.StatusForbidden,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if nil != c.reqBody {
				reqBodyInBytes, _ = json.Marshal(c.reqBody)
			}
			if nil != c.cookie {
				resp = serveHTTPWithCookies("POST", path, string(reqBodyInBytes), "application/json", c.authorization, *c.cookie)
			} else {
				resp = serveHTTP("POST", path, "", "application/json", c.authorization)
			}
			assert.Equal(t, c.resultCode, resp.Code)

		})
	}
}

func testDonationCreateSuccess(t *testing.T, path string, userID uint, frequency string, paymethod string, authorization string, cookie http.Cookie) {
	var resp *httptest.ResponseRecorder
	var resBody responseBody
	var reqBodyInBytes []byte
	var resBodyInBytes []byte

	testCardholder := models.Cardholder{
		PhoneNumber: null.StringFrom("+886912345678"),
		Name:        null.StringFrom("王小明"),
		Email:       "developer@twreporter.org",
		ZipCode:     null.StringFrom("104"),
		Address:     null.StringFrom("台北市中山區南京東路X巷X號X樓"),
		NationalID:  null.StringFrom("A123456789"),
	}

	cases := []struct {
		name    string
		reqBody requestBody
	}{
		{
			name: "StatusCode=StatusCreated,full fields",
			reqBody: requestBody{
				Amount:     testAmount,
				Cardholder: testCardholder,
				Currency:   testCurrency,
				Details:    testDetails,
				MerchantID: testMerchantID,
				Prime:      testPrime,
				UserID:     userID,
			},
		},
		{
			name: "StatusCode=StatusCreated,minimum fields",
			reqBody: requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				Frequency:  frequency,
				MerchantID: testMerchantID,
				Prime:      testPrime,
				UserID:     userID,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			const defaultAnonymity = false
			if frequency == oneTimeFrequency {
				c.reqBody.PayMethod = paymethod
			} else {
				c.reqBody.Frequency = frequency
			}

			reqBodyInBytes, _ = json.Marshal(c.reqBody)
			resp = serveHTTPWithCookies("POST", path, string(reqBodyInBytes), "application/json", authorization, cookie)
			resBodyInBytes, _ = ioutil.ReadAll(resp.Result().Body)
			json.Unmarshal(resBodyInBytes, &resBody)

			assert.Equal(t, http.StatusCreated, resp.Code)
			assert.Equal(t, "success", resBody.Status)
			assert.Equal(t, testAmount, resBody.Data.Amount)
			assert.Equal(t, frequency, resBody.Data.Frequency)

			if c.reqBody.Currency != "" {
				assert.Equal(t, c.reqBody.Currency, resBody.Data.Currency)
			} else {
				assert.Equal(t, defaultCurrency, resBody.Data.Currency)
			}
			if c.reqBody.Details != "" {
				assert.Equal(t, c.reqBody.Details, resBody.Data.Details)
			} else {
				switch frequency {
				case oneTimeFrequency:
					assert.Equal(t, defaultOneTimeDetails, resBody.Data.Details)
				case monthlyFrequency, yearlyFrequency:
					assert.Equal(t, defaultPeriodicDetails, resBody.Data.Details)
				}
			}

			assert.NotEmpty(t, resBody.Data.OrderNumber)
			assert.Empty(t, resBody.Data.Notes)

			assert.Equal(t, defaultAnonymity, resBody.Data.IsAnonymous)
			assert.Equal(t, c.reqBody.Cardholder.PhoneNumber.ValueOrZero(), resBody.Data.Cardholder.PhoneNumber.ValueOrZero())
			assert.Equal(t, c.reqBody.Cardholder.Name.ValueOrZero(), resBody.Data.Cardholder.Name.ValueOrZero())
			assert.Equal(t, c.reqBody.Cardholder.Email, resBody.Data.Cardholder.Email)
			assert.Equal(t, c.reqBody.Cardholder.ZipCode.ValueOrZero(), resBody.Data.Cardholder.ZipCode.ValueOrZero())
			assert.Equal(t, c.reqBody.Cardholder.NationalID.ValueOrZero(), resBody.Data.Cardholder.NationalID.ValueOrZero())
			assert.Equal(t, c.reqBody.Cardholder.Address.ValueOrZero(), resBody.Data.Cardholder.Address.ValueOrZero())

		})
	}

}

func TestCreateAOneTimeDonation(t *testing.T) {
	var path = "/v1/donations/prime"
	var user models.User

	user = getUser(Globs.Defaults.Account)
	authorization, cookie := helperSetupAuth(user)

	testDonationCreateClientError(t, path, user.ID, oneTimeFrequency, creditCardPayMethod, authorization, cookie)
	testDonationCreateServerError(t, path, user.ID, oneTimeFrequency, creditCardPayMethod, authorization, cookie)
	// ==========================================
	// Test One Time Donation Creation
	// Pay by credit card
	// =========================================
	testDonationCreateSuccess(t, path, user.ID, oneTimeFrequency, creditCardPayMethod, authorization, cookie)
}

func TestCreateAPeriodicDonation(t *testing.T) {
	var path = "/v1/periodic-donations"
	var user models.User

	user = getUser(Globs.Defaults.Account)
	authorization, cookie := helperSetupAuth(user)

	testDonationCreateClientError(t, path, user.ID, monthlyFrequency, creditCardPayMethod, authorization, cookie)
	testDonationCreateServerError(t, path, user.ID, monthlyFrequency, creditCardPayMethod, authorization, cookie)
	// ==========================================
	// Test Periodic Donation Creation
	// =========================================
	testDonationCreateSuccess(t, path, user.ID, monthlyFrequency, creditCardPayMethod, authorization, cookie)
}

func createDefaultDonationRecord(reqBody requestBody, endpoint string, user models.User) responseBody {
	authorization, cookie := helperSetupAuth(user)

	// create a donation by HTTP POST request
	reqBodyInBytes, _ := json.Marshal(reqBody)
	resp := serveHTTPWithCookies("POST", endpoint, string(reqBodyInBytes), "application/json", authorization, cookie)
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

	reqBody := requestBody{
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Address:     null.StringFrom(testAddress),
			Email:       user.Email.ValueOrZero(),
			Name:        null.StringFrom(testName),
			NationalID:  null.StringFrom(testNationalID),
			PhoneNumber: null.StringFrom(testPhoneNumber),
			ZipCode:     null.StringFrom(testZipCode),
		},
		Details:    testDetails,
		Frequency:  monthlyFrequency,
		MerchantID: testMerchantID,
		Prime:      testPrime,
		UserID:     user.ID,
	}

	return createDefaultDonationRecord(reqBody, path, user)
}

func createDefaultPrimeDonationRecord(user models.User) responseBody {
	// create a default prime donation record
	path := "/v1/donations/prime"

	reqBody := requestBody{
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Address:     null.StringFrom(testAddress),
			Email:       user.Email.ValueOrZero(),
			Name:        null.StringFrom(testName),
			NationalID:  null.StringFrom(testNationalID),
			PhoneNumber: null.StringFrom(testPhoneNumber),
			ZipCode:     null.StringFrom(testZipCode),
		},
		Details:    testDetails,
		MerchantID: testMerchantID,
		PayMethod:  creditCardPayMethod,
		Prime:      testPrime,
		UserID:     user.ID,
	}

	return createDefaultDonationRecord(reqBody, path, user)
}

func testDonationPatchClientError(t *testing.T, userID uint, frequency, orderNumber, authorization string, cookie http.Cookie) {
	var reqBodyInBytes []byte
	var resp *httptest.ResponseRecorder
	const invalidUserID = 1000

	cases := []struct {
		name          string
		reqBody       *map[string]interface{}
		cookie        *http.Cookie
		authorization string
		orderNumber   string
		resultCode    int
	}{
		{
			name:          "StatusCode=StatusUnauthorized,Lack of Cookie",
			authorization: authorization,
			orderNumber:   orderNumber,
			resultCode:    http.StatusUnauthorized,
		},
		{
			name:        "StatusCode=StatusUnauthorized,Lack of Authorization Header",
			cookie:      &cookie,
			orderNumber: orderNumber,
			resultCode:  http.StatusUnauthorized,
		},
		{
			name:          "StatusCode=StatusForbidden,Unauthorized Resource",
			reqBody:       &map[string]interface{}{"user_id": invalidUserID},
			cookie:        &cookie,
			authorization: authorization,
			orderNumber:   orderNumber,
			resultCode:    http.StatusForbidden,
		},
		{
			name: "StatusCode=StatusBadRequest,Unauthorized Resource",
			reqBody: &map[string]interface{}{
				"user_id": userID,
				// to_feedback should be boolean
				"to_feedback": "false",
				// national_id should be string
				"donor": map[string]interface{}{
					"national_id": true,
				},
			},
			cookie:        &cookie,
			authorization: authorization,
			orderNumber:   orderNumber,
			resultCode:    http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusNotFound,Invalid Order Number",
			reqBody: &map[string]interface{}{
				"send_receipt": "no",
				"to_feedback":  !testFeedback,
				"user_id":      userID,
			},
			cookie:        &cookie,
			authorization: authorization,
			orderNumber:   "InvalidOrderNumber",
			resultCode:    http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := ""
			switch frequency {
			case oneTimeFrequency:
				path = "/v1/donations/prime/orders/" + c.orderNumber
			case monthlyFrequency, yearlyFrequency:
				path = "/v1/periodic-donations/orders/" + c.orderNumber
			}

			if nil != c.reqBody {
				reqBodyInBytes, _ = json.Marshal(c.reqBody)
			}
			if nil != c.cookie {
				resp = serveHTTPWithCookies("PATCH", path, string(reqBodyInBytes), "application/json", c.authorization, *c.cookie)
			} else {
				resp = serveHTTP("PATCH", path, "", "application/json", c.authorization)
			}
			assert.Equal(t, c.resultCode, resp.Code)

		})
	}
}

func TestPatchAPeriodicDonation(t *testing.T) {
	const donorEmail string = "periodic-donor@twreporter.org"
	var defaultRecordRes responseBody
	var path string
	var reqBody map[string]interface{}
	var reqBodyInBytes []byte
	var resp *httptest.ResponseRecorder
	var user models.User

	// setup before test
	// create a new user
	user = createUser(donorEmail)
	authorization, cookie := helperSetupAuth(user)
	// get record to patch
	defaultRecordRes = createDefaultPeriodicDonationRecord(user)

	testDonationPatchClientError(t, user.ID, defaultRecordRes.Data.Frequency, defaultRecordRes.Data.OrderNumber, authorization, cookie)
	path = fmt.Sprintf("/v1/periodic-donations/orders/%s", defaultRecordRes.Data.OrderNumber)

	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		var dataAfterPatch models.PeriodicDonation
		const testIsAnonymous = true
		reqBody = map[string]interface{}{
			"donor": map[string]string{
				"address": "test-addres",
				"name":    "test-name",
			},
			"send_receipt": "no",
			"to_feedback":  !testFeedback,
			"is_anonymous": null.BoolFrom(testIsAnonymous),
			"user_id":      user.ID,
		}
		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTPWithCookies("PATCH", path, string(reqBodyInBytes), "application/json", authorization, cookie)
		assert.Equal(t, http.StatusNoContent, resp.Code)

		Globs.GormDB.Where("id = ?", defaultRecordRes.Data.ID).Find(&dataAfterPatch)
		assert.Equal(t, reqBody["to_feedback"], dataAfterPatch.ToFeedback.ValueOrZero())
		assert.Equal(t, reqBody["send_receipt"], dataAfterPatch.SendReceipt)
		assert.Equal(t, reqBody["is_anonymous"], dataAfterPatch.IsAnonymous)
		assert.Equal(t, reqBody["donor"].(map[string]string)["address"], dataAfterPatch.Cardholder.Address.ValueOrZero())
		assert.Equal(t, reqBody["donor"].(map[string]string)["name"], dataAfterPatch.Cardholder.Name.ValueOrZero())
	})
}

func TestPatchAPrimeDonation(t *testing.T) {
	const donorEmail string = "prim-donor@twreporter.org"
	var defaultRecordRes responseBody
	var path string
	var reqBody map[string]interface{}
	var reqBodyInBytes []byte
	var resp *httptest.ResponseRecorder
	var user models.User

	// setup before test
	// create a new user
	user = createUser(donorEmail)
	authorization, cookie := helperSetupAuth(user)

	// get record to patch
	defaultRecordRes = createDefaultPrimeDonationRecord(user)

	testDonationPatchClientError(t, user.ID, oneTimeFrequency, defaultRecordRes.Data.OrderNumber, authorization, cookie)

	path = fmt.Sprintf("/v1/donations/prime/orders/%s", defaultRecordRes.Data.OrderNumber)
	t.Run("StatusCode=StatusNoContent", func(t *testing.T) {
		var dataAfterPatch models.PayByPrimeDonation
		const testIsAnonymous = true
		reqBody = map[string]interface{}{
			"donor": map[string]string{
				"name":    "test-name",
				"address": "test-addres",
			},
			"send_receipt": "no",
			"is_anonymous": null.BoolFrom(testIsAnonymous),
			"user_id":      user.ID,
		}
		reqBodyInBytes, _ = json.Marshal(reqBody)
		resp = serveHTTPWithCookies("PATCH", path, string(reqBodyInBytes), "application/json", authorization, cookie)
		assert.Equal(t, http.StatusNoContent, resp.Code)

		Globs.GormDB.Where("id = ?", defaultRecordRes.Data.ID).Find(&dataAfterPatch)
		assert.Equal(t, reqBody["send_receipt"], dataAfterPatch.SendReceipt)
		assert.Equal(t, reqBody["is_anonymous"], dataAfterPatch.IsAnonymous)
		assert.Equal(t, reqBody["donor"].(map[string]string)["address"], dataAfterPatch.Cardholder.Address.ValueOrZero())
		assert.Equal(t, reqBody["donor"].(map[string]string)["name"], dataAfterPatch.Cardholder.Name.ValueOrZero())
	})
}

func TestGetAPrimeDonationOfAUser(t *testing.T) {
	// setup before test
	donorEmail := "get-prime-donor@twreporter.org"
	// create a new user
	user := createUser(donorEmail)
	authorization, cookie := helperSetupAuth(user)

	primeRes := createDefaultPrimeDonationRecord(user)

	maliciousDonorEmail := "get-others-prime-donor@twreporter.org"
	maliciousUser := createUser(maliciousDonorEmail)
	maliciousAuthorization, maliciousCookie := helperSetupAuth(maliciousUser)

	path := fmt.Sprintf("/v1/donations/prime/orders/%s", primeRes.Data.OrderNumber)
	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		resp := serveHTTPWithCookies("GET", path, "", "application/json", "", cookie)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		resp := serveHTTPWithCookies("GET", path, "", "application/json", maliciousAuthorization, maliciousCookie)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		invalidPath := "/v1/donations/prime/orders/INVALID_ORDER"
		resp := serveHTTPWithCookies("GET", invalidPath, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("StatusCode=StatusOK", func(t *testing.T) {
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		respInBytes, _ := ioutil.ReadAll(resp.Result().Body)
		defer resp.Result().Body.Close()

		// parse response into struct
		resBody := responseBody{}
		json.Unmarshal(respInBytes, &resBody)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, testAmount, resBody.Data.Amount)
		assert.Equal(t, testDetails, resBody.Data.Details)
		assert.Equal(t, donorEmail, resBody.Data.Cardholder.Email)
		assert.Equal(t, testAddress, resBody.Data.Cardholder.Address.ValueOrZero())
		assert.Equal(t, testName, resBody.Data.Cardholder.Name.ValueOrZero())
		assert.Equal(t, testNationalID, resBody.Data.Cardholder.NationalID.ValueOrZero())
		assert.Equal(t, testPhoneNumber, resBody.Data.Cardholder.PhoneNumber.ValueOrZero())
		assert.Equal(t, testZipCode, resBody.Data.Cardholder.ZipCode.ValueOrZero())
		assert.Equal(t, "4242", resBody.Data.CardInfo.LastFour.ValueOrZero())
		assert.Equal(t, "424242", resBody.Data.CardInfo.BinCode.ValueOrZero())
		assert.Equal(t, int64(0), resBody.Data.CardInfo.Funding.ValueOrZero())
		assert.Equal(t, testCurrency, resBody.Data.Currency)
		assert.Equal(t, "credit_card", resBody.Data.PayMethod)
		assert.Equal(t, "yearly", resBody.Data.SendReceipt)
		assert.Equal(t, false, resBody.Data.ToFeedback)
		assert.Equal(t, oneTimeFrequency, resBody.Data.Frequency)
		assert.Empty(t, resBody.Data.Notes)
		assert.NotEmpty(t, resBody.Data.OrderNumber)
	})
}

func TestGetAPeriodicDonationOfAUser(t *testing.T) {
	// setup before test
	donorEmail := "get-periodic-donor@twreporter.org"
	// create a new user
	user := createUser(donorEmail)
	authorization, cookie := helperSetupAuth(user)
	periodicRes := createDefaultPeriodicDonationRecord(user)

	maliciousDonorEmail := "get-others-periodic-donor@twreporter.org"
	maliciousUser := createUser(maliciousDonorEmail)
	maliciousAuthorization, maliciousCookie := helperSetupAuth(maliciousUser)

	path := fmt.Sprintf("/v1/periodic-donations/orders/%s", periodicRes.Data.OrderNumber)
	t.Run("StatusCode=StatusUnauthorized", func(t *testing.T) {
		resp := serveHTTPWithCookies("GET", path, "", "application/json", "", cookie)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("StatusCode=StatusForbidden", func(t *testing.T) {
		resp := serveHTTPWithCookies("GET", path, "", "application/json", maliciousAuthorization, maliciousCookie)
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("StatusCode=StatusNotFound", func(t *testing.T) {
		invalidPath := "/v1/periodic-donations/orders/INVALID_ORDER"
		resp := serveHTTPWithCookies("GET", invalidPath, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("StatusCode=StatusOK", func(t *testing.T) {
		resp := serveHTTPWithCookies("GET", path, "", "application/json", authorization, cookie)
		assert.Equal(t, http.StatusOK, resp.Code)
		respInBytes, _ := ioutil.ReadAll(resp.Result().Body)
		defer resp.Result().Body.Close()

		// parse response into struct
		resBody := responseBody{}
		json.Unmarshal(respInBytes, &resBody)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, testAmount, resBody.Data.Amount)
		assert.Equal(t, testDetails, resBody.Data.Details)
		assert.Equal(t, donorEmail, resBody.Data.Cardholder.Email)
		assert.Equal(t, testAddress, resBody.Data.Cardholder.Address.ValueOrZero())
		assert.Equal(t, testName, resBody.Data.Cardholder.Name.ValueOrZero())
		assert.Equal(t, testNationalID, resBody.Data.Cardholder.NationalID.ValueOrZero())
		assert.Equal(t, testPhoneNumber, resBody.Data.Cardholder.PhoneNumber.ValueOrZero())
		assert.Equal(t, testZipCode, resBody.Data.Cardholder.ZipCode.ValueOrZero())
		assert.Equal(t, "4242", resBody.Data.CardInfo.LastFour.ValueOrZero())
		assert.Equal(t, "424242", resBody.Data.CardInfo.BinCode.ValueOrZero())
		assert.Equal(t, int64(0), resBody.Data.CardInfo.Funding.ValueOrZero())
		assert.Equal(t, testCurrency, resBody.Data.Currency)
		assert.Equal(t, "yearly", resBody.Data.SendReceipt)
		assert.Equal(t, true, resBody.Data.ToFeedback)
		assert.Equal(t, monthlyFrequency, resBody.Data.Frequency)
		assert.Empty(t, resBody.Data.Notes)
		assert.NotEmpty(t, resBody.Data.OrderNumber)
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
	user := getUser(donorEmail)
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
