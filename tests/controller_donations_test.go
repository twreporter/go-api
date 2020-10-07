package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

type (
	linePayResultURL struct {
		FrontendRedirectURL string `json:"frontend_redirect_url"`
		BackendRedirectURL  string `json:"backend_redirect_url"`
	}
	donationRecord struct {
		Amount        uint              `json:"amount"`
		Cardholder    models.Cardholder `json:"cardholder"`
		Currency      string            `json:"currency"`
		Details       string            `json:"details"`
		Frequency     string            `json:"frequency"`
		ID            uint              `json:"id"`
		Notes         string            `json:"notes"`
		OrderNumber   string            `json:"order_number"`
		PayMethod     string            `json:"pay_method"`
		SendReceipt   string            `json:"send_receipt"`
		ToFeedback    bool              `json:"to_feedback"`
		IsAnonymous   bool              `json:"is_anonymous"`
		PaymentUrl    string            `json:"payment_url"`
		ReceiptHeader null.String       `json:"receipt_header"`
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

	reqHeader struct {
		Cookie        *http.Cookie
		Authorization string
	}

	tapPayRequestBody struct {
		models.PayInfo    `json:"pay_info"`
		RecTradeID        string `json:"rec_trade_id"`
		BankTransactionID string `json:"bank_transaction_id"`
		OrderNumber       string `json:"order_number"`
		Amount            uint   `json:"amount"`
		Status            int    `json:"status"`
		TransactionTime   int64  `json:"transaction_time_millis"`
		Acquirer          string `json:"acquirer"`
		BankResultCode    string `json:"bank_result_code"`
		BankResultMsg     string `json:"bank_result_msg"`
	}
)

const (
	testCreditCardPrime = "test_3a2fb2b7e892b914a03c95dd4dd5dc7970c908df67a49527c0a648b2bc9"
	testLinePrime       = "ln_test_utigjeyfutj5867uyjhuty47rythfjru485768tigjfheufhtu5i6ojk"

	testCreditCardMerchant = "GlobalTesting_CTBC"
	testLineMerchant       = "GlobalTesting_LINEPAY"

	testDetails       = "報導者小額捐款"
	testAmount   uint = 500
	testCurrency      = "TWD"
	testFeedback      = true

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

	oneTimeOrderPathPrefix  = "/v1/donations/prime/orders/"
	periodicOrderPathPrefix = "/v1/periodic-donations/orders/"

	statusPaying = "paying"
	statusPaid   = "paid"
	statusFail   = "fail"
)

var methodToPrime = map[string]string{
	creditCardPayMethod: testCreditCardPrime,
	linePayMethod:       testLinePrime,
}

var methodToMerchant = map[string]string{
	creditCardPayMethod: testCreditCardMerchant,
	linePayMethod:       testLineMerchant,
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

func helperSetupAuth(user models.User) (authorization string, cookie http.Cookie) {
	jwt := generateIDToken(user)
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

func testDonationCreateServerError(t *testing.T, path string, userID uint, frequency string, payMethod string, authorization string, cookie http.Cookie) {
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
				Prime:  "test_prime_which_will_occurs_error",
				UserID: userID,
			},
			resultCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if nil != c.reqBody {
				switch frequency {
				case oneTimeFrequency:
					c.reqBody.PayMethod = payMethod
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

func testDonationCreateClientError(t *testing.T, path string, userID uint, frequency string, payMethod string, authorization string, cookie http.Cookie) {
	var resp *httptest.ResponseRecorder
	var reqBodyInBytes []byte

	header := reqHeader{
		Cookie:        &cookie,
		Authorization: authorization,
	}

	cases := []struct {
		reqHeader
		name       string
		reqBody    *requestBody
		resultCode int
	}{
		{
			name: "StatusCode=StatusBadRequest,Lack of UserID",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     methodToPrime[payMethod],
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
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
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of Cardholder",
			reqBody: &requestBody{
				Amount:    testAmount,
				PayMethod: creditCardPayMethod,
				Prime:     methodToPrime[payMethod],
				UserID:    userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
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
				Prime:     methodToPrime[payMethod],
				UserID:    userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of Amount",
			reqBody: &requestBody{
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     methodToPrime[payMethod],
				UserID:    userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Lack of PayMethod",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				Prime:  methodToPrime[payMethod],
				UserID: userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Malformed email",
			reqBody: &requestBody{
				Amount: testAmount,
				Cardholder: models.Cardholder{
					Email: "developer-twreporter,org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     methodToPrime[payMethod],
				UserID:    userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusBadRequest,Amount less than 1",
			reqBody: &requestBody{
				Amount: 0,
				Cardholder: models.Cardholder{
					Email: "developer@twreporter.org",
				},
				PayMethod: creditCardPayMethod,
				Prime:     methodToPrime[payMethod],
				UserID:    userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
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
				Prime:     methodToPrime[payMethod],
				UserID:    userID,
			},
			reqHeader:  header,
			resultCode: http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusUnauthorized,Without Cookie",
			reqHeader: reqHeader{
				Authorization: authorization,
			},
			resultCode: http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusUnauthorized,Without Authorization Header",
			reqHeader: reqHeader{
				Cookie: &cookie,
			},
			resultCode: http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusForbidden,Unauthorized Resource",
			reqBody: &requestBody{
				UserID: 1000,
			},
			reqHeader:  header,
			resultCode: http.StatusForbidden,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if nil != c.reqBody {
				reqBodyInBytes, _ = json.Marshal(c.reqBody)
			}
			if nil != c.Cookie {
				resp = serveHTTPWithCookies("POST", path, string(reqBodyInBytes), "application/json", c.Authorization, *c.Cookie)
			} else {
				resp = serveHTTP("POST", path, "", "application/json", c.Authorization)
			}
			assert.Equal(t, c.resultCode, resp.Code)

		})
	}
}

func testDonationCreateSuccess(t *testing.T, path string, userID uint, frequency string, payMethod string, authorization string, cookie http.Cookie) {
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
				MerchantID: methodToMerchant[payMethod],
				Prime:      methodToPrime[payMethod],
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
				MerchantID: methodToMerchant[payMethod],
				Prime:      methodToPrime[payMethod],
				UserID:     userID,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			const defaultAnonymity = false
			if frequency == oneTimeFrequency {
				c.reqBody.PayMethod = payMethod
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

			if payMethod == linePayMethod {
				assert.NotEmpty(t, resBody.Data.PaymentUrl)
			}

		})
	}

}

func TestCreateAOneTimeDonation(t *testing.T) {
	var path = "/v1/donations/prime"
	var user models.User

	user = getUser(Globs.Defaults.Account)
	authorization, cookie := helperSetupAuth(user)

	payMethods := []string{
		creditCardPayMethod,
		linePayMethod,
	}

	for _, p := range payMethods {

		testDonationCreateClientError(t, path, user.ID, oneTimeFrequency, p, authorization, cookie)
		testDonationCreateServerError(t, path, user.ID, oneTimeFrequency, p, authorization, cookie)
		testDonationCreateSuccess(t, path, user.ID, oneTimeFrequency, p, authorization, cookie)
	}
}

func TestCreateAPeriodicDonation(t *testing.T) {
	var path = "/v1/periodic-donations"
	var user models.User

	user = getUser(Globs.Defaults.Account)
	authorization, cookie := helperSetupAuth(user)

	frequency := []string{
		monthlyFrequency,
		yearlyFrequency,
	}

	for _, f := range frequency {
		testDonationCreateClientError(t, path, user.ID, f, creditCardPayMethod, authorization, cookie)
		testDonationCreateServerError(t, path, user.ID, f, creditCardPayMethod, authorization, cookie)
		testDonationCreateSuccess(t, path, user.ID, f, creditCardPayMethod, authorization, cookie)
	}
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

func createDefaultPeriodicDonationRecord(user models.User, frequency string) responseBody {
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
		Frequency:  frequency,
		MerchantID: testCreditCardMerchant,
		Prime:      testCreditCardPrime,
		UserID:     user.ID,
	}

	return createDefaultDonationRecord(reqBody, path, user)
}

func createDefaultPrimeDonationRecord(user models.User, payMethod string) responseBody {
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
		MerchantID: methodToMerchant[payMethod],
		PayMethod:  payMethod,
		Prime:      methodToPrime[payMethod],
		UserID:     user.ID,
	}

	return createDefaultDonationRecord(reqBody, path, user)
}

func testDonationPatchClientError(t *testing.T, userID uint, pathPrefix, authorization string, cookie http.Cookie) {
	var reqBodyInBytes []byte
	var resp *httptest.ResponseRecorder
	const invalidUserID = 1000

	const orderNumber = "ValidOrderNumber"
	header := reqHeader{
		Cookie:        &cookie,
		Authorization: authorization,
	}

	cases := []struct {
		reqHeader
		name        string
		reqBody     *map[string]interface{}
		orderNumber string
		resultCode  int
	}{
		{
			name: "StatusCode=StatusUnauthorized,Lack of Cookie",
			reqHeader: reqHeader{
				Authorization: authorization,
			},
			orderNumber: orderNumber,
			resultCode:  http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusUnauthorized,Lack of Authorization Header",
			reqHeader: reqHeader{
				Cookie: &cookie,
			},
			orderNumber: orderNumber,
			resultCode:  http.StatusUnauthorized,
		},
		{
			name:        "StatusCode=StatusForbidden,Unauthorized Resource",
			reqBody:     &map[string]interface{}{"user_id": invalidUserID},
			reqHeader:   header,
			orderNumber: orderNumber,
			resultCode:  http.StatusForbidden,
		},
		{
			name: "StatusCode=StatusBadRequest,Incorrect request body format",
			reqBody: &map[string]interface{}{
				"user_id": userID,
				// to_feedback should be boolean
				"to_feedback": "false",
				// national_id should be string
				"donor": map[string]interface{}{
					"national_id": true,
				},
			},
			reqHeader:   header,
			orderNumber: orderNumber,
			resultCode:  http.StatusBadRequest,
		},
		{
			name: "StatusCode=StatusNotFound,Invalid Order Number",
			reqBody: &map[string]interface{}{
				"send_receipt": "no",
				"to_feedback":  !testFeedback,
				"user_id":      userID,
			},
			reqHeader:   header,
			orderNumber: "InvalidOrderNumber",
			resultCode:  http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := pathPrefix + c.orderNumber

			if nil != c.reqBody {
				reqBodyInBytes, _ = json.Marshal(c.reqBody)
			}
			if nil != c.Cookie {
				resp = serveHTTPWithCookies("PATCH", path, string(reqBodyInBytes), "application/json", c.Authorization, *c.Cookie)
			} else {
				resp = serveHTTP("PATCH", path, "", "application/json", c.Authorization)
			}
			assert.Equal(t, c.resultCode, resp.Code)

		})
	}
}

func TestPatchAPeriodicDonation(t *testing.T) {
	// setup before test
	// create a new user
	user := createUser("periodic-donor@twreporter.org")
	defer func() { deleteUser(user) }()
	authorization, cookie := helperSetupAuth(user)

	frequency := []string{
		monthlyFrequency,
		yearlyFrequency,
	}

	for _, f := range frequency {
		testDonationPatchClientError(t, user.ID, periodicOrderPathPrefix, authorization, cookie)
		testPeriodicDonationPatchSuccess(t, f, user, authorization, cookie)
	}
}
func testPeriodicDonationPatchSuccess(t *testing.T, frequency string, user models.User, authorization string, cookie http.Cookie) {
	for _, tc := range []struct {
		name        string
		existRecord models.PeriodicDonation
		reqBody     map[string]interface{}
	}{
		{
			name: "StatusCode=StatusNoContent,Patch fields with changes",
			existRecord: models.PeriodicDonation{
				Amount:      500,
				Currency:    "TWD",
				Details:     "test donation",
				OrderNumber: "twrepoter-test-order-number",
				Status:      "paid",
				UserID:      user.ID,
				Cardholder: models.Cardholder{
					Email: user.Email.String,
				},
			},
			reqBody: map[string]interface{}{
				"donor": map[string]string{
					"name":    "test-name",
					"address": "test-addres",
				},
				"send_receipt":   "no",
				"is_anonymous":   null.BoolFrom(true),
				"to_feedback":    false,
				"user_id":        user.ID,
				"receipt_header": "mock header",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Create existing Records
			tc.existRecord.Frequency = frequency
			Globs.GormDB.Create(&tc.existRecord)
			defer func() {
				Globs.GormDB.Unscoped().Delete(&tc.existRecord)
			}()

			reqBodyInBytes, _ := json.Marshal(tc.reqBody)

			path := periodicOrderPathPrefix + tc.existRecord.OrderNumber
			resp := serveHTTPWithCookies("PATCH", path, string(reqBodyInBytes), "application/json", authorization, cookie)
			assert.Equal(t, http.StatusNoContent, resp.Code)

			var dataAfterPatch models.PeriodicDonation
			Globs.GormDB.Where("id = ?", tc.existRecord.ID).Find(&dataAfterPatch)
			assert.Equal(t, tc.reqBody["send_receipt"], dataAfterPatch.SendReceipt)
			assert.Equal(t, tc.reqBody["to_feedback"], dataAfterPatch.ToFeedback.ValueOrZero())
			assert.Equal(t, tc.reqBody["is_anonymous"], dataAfterPatch.IsAnonymous)
			assert.Equal(t, tc.reqBody["donor"].(map[string]string)["address"], dataAfterPatch.Cardholder.Address.ValueOrZero())
			assert.Equal(t, tc.reqBody["donor"].(map[string]string)["name"], dataAfterPatch.Cardholder.Name.ValueOrZero())
			assert.Equal(t, tc.reqBody["receipt_header"], dataAfterPatch.ReceiptHeader)
		})
	}
}
func testPrimeDonationPatchSuccess(t *testing.T, payMethod string, user models.User, authorization string, cookie http.Cookie) {

	for _, tc := range []struct {
		name        string
		existRecord models.PayByPrimeDonation
		reqBody     map[string]interface{}
	}{
		{
			name: "StatusCode=StatusNoContent,Patch fields with changes",
			existRecord: models.PayByPrimeDonation{
				Amount:      500,
				Currency:    "TWD",
				Details:     "test donation",
				MerchantID:  "test merchant",
				OrderNumber: "twrepoter-test-order-number",
				Status:      "paid",
				UserID:      user.ID,
				Cardholder: models.Cardholder{
					Email: user.Email.String,
				},
			},
			reqBody: map[string]interface{}{
				"donor": map[string]string{
					"name":    "test-name",
					"address": "test-addres",
				},
				"send_receipt":   "no",
				"is_anonymous":   null.BoolFrom(true),
				"user_id":        user.ID,
				"receipt_header": "mock header",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Create existing Records
			tc.existRecord.PayMethod = payMethod
			Globs.GormDB.Create(&tc.existRecord)
			defer func() {
				Globs.GormDB.Unscoped().Delete(&tc.existRecord)
			}()

			reqBodyInBytes, _ := json.Marshal(tc.reqBody)

			path := oneTimeOrderPathPrefix + tc.existRecord.OrderNumber
			resp := serveHTTPWithCookies("PATCH", path, string(reqBodyInBytes), "application/json", authorization, cookie)
			assert.Equal(t, http.StatusNoContent, resp.Code)

			var dataAfterPatch models.PayByPrimeDonation
			Globs.GormDB.Where("id = ?", tc.existRecord.ID).Find(&dataAfterPatch)
			assert.Equal(t, tc.reqBody["send_receipt"], dataAfterPatch.SendReceipt)
			assert.Equal(t, tc.reqBody["is_anonymous"], dataAfterPatch.IsAnonymous)
			assert.Equal(t, tc.reqBody["donor"].(map[string]string)["address"], dataAfterPatch.Cardholder.Address.ValueOrZero())
			assert.Equal(t, tc.reqBody["donor"].(map[string]string)["name"], dataAfterPatch.Cardholder.Name.ValueOrZero())
			assert.Equal(t, tc.reqBody["receipt_header"], dataAfterPatch.ReceiptHeader)
		})
	}
}

func TestPatchAPrimeDonation(t *testing.T) {
	// setup before test
	// create a new user
	user := createUser("prim-donor@twreporter.org")
	defer func() { deleteUser(user) }()
	authorization, cookie := helperSetupAuth(user)

	payMethods := []string{
		creditCardPayMethod,
		linePayMethod,
	}

	for _, p := range payMethods {
		testDonationPatchClientError(t, user.ID, oneTimeOrderPathPrefix, authorization, cookie)
		testPrimeDonationPatchSuccess(t, p, user, authorization, cookie)
	}
}

func testDonationGetClientError(t *testing.T, pathPrefix, orderNumber, authorization string, cookie http.Cookie) {
	var reqBodyInBytes []byte
	var resp *httptest.ResponseRecorder

	maliciousUser := createUser("malicious-donor@twreporter.org")
	defer func() { deleteUser(maliciousUser) }()
	maliciousAuthorization, maliciousCookie := helperSetupAuth(maliciousUser)

	cases := []struct {
		reqHeader
		name        string
		orderNumber string
		resultCode  int
	}{
		{
			name: "StatusCode=StatusUnauthorized,Lack of Authorization Header",
			reqHeader: reqHeader{
				Cookie: &cookie,
			},
			orderNumber: orderNumber,
			resultCode:  http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusForbidden,Unauthorized Resource",
			reqHeader: reqHeader{
				Cookie:        &maliciousCookie,
				Authorization: maliciousAuthorization,
			},
			orderNumber: orderNumber,
			resultCode:  http.StatusForbidden,
		},
		{
			name: "StatusCode=StatusNotFound,Invalid Order Number",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			orderNumber: "InvalidOrderNumber",
			resultCode:  http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := pathPrefix + c.orderNumber

			resp = serveHTTPWithCookies("GET", path, string(reqBodyInBytes), "application/json", c.Authorization, *c.Cookie)
			assert.Equal(t, c.resultCode, resp.Code)

		})
	}

}

func TestGetAPrimeDonationOfAUser(t *testing.T) {
	// setup before test
	// create a new user
	donorEmail := "get-prime-donor@twreporter.org"
	user := createUser(donorEmail)
	defer func() { deleteUser(user) }()
	authorization, cookie := helperSetupAuth(user)
	payMethods := []string{
		creditCardPayMethod,
		linePayMethod,
	}

	for _, p := range payMethods {

		primeRes := createDefaultPrimeDonationRecord(user, p)

		testDonationGetClientError(t, oneTimeOrderPathPrefix, primeRes.Data.OrderNumber, authorization, cookie)
		path := oneTimeOrderPathPrefix + primeRes.Data.OrderNumber

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
			assert.Equal(t, testCurrency, resBody.Data.Currency)
			assert.Equal(t, p, resBody.Data.PayMethod)
			assert.Equal(t, "no", resBody.Data.SendReceipt)
			assert.Empty(t, resBody.Data.Notes)
			assert.NotEmpty(t, resBody.Data.OrderNumber)
			if assert.False(t, resBody.Data.ReceiptHeader.IsZero()) {
				assert.Empty(t, resBody.Data.ReceiptHeader.String)
			}
		})
	}
}

func TestGetAPeriodicDonationOfAUser(t *testing.T) {
	// setup before test
	donorEmail := "get-periodic-donor@twreporter.org"
	// create a new user
	user := createUser(donorEmail)
	defer func() { deleteUser(user) }()
	authorization, cookie := helperSetupAuth(user)

	frequency := []string{
		monthlyFrequency,
		yearlyFrequency,
	}

	for _, f := range frequency {
		periodicRes := createDefaultPeriodicDonationRecord(user, f)

		testDonationGetClientError(t, periodicOrderPathPrefix, periodicRes.Data.OrderNumber, authorization, cookie)

		path := periodicOrderPathPrefix + periodicRes.Data.OrderNumber

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
			assert.Equal(t, testCurrency, resBody.Data.Currency)
			assert.Equal(t, "no", resBody.Data.SendReceipt)
			assert.Equal(t, true, resBody.Data.ToFeedback)
			assert.Equal(t, f, resBody.Data.Frequency)
			assert.Empty(t, resBody.Data.Notes)
			assert.NotEmpty(t, resBody.Data.OrderNumber)
			if assert.False(t, resBody.Data.ReceiptHeader.IsZero()) {
				assert.Empty(t, resBody.Data.ReceiptHeader.String)
			}
		})
	}
}

func TestLinePayNotify(t *testing.T) {
	const (
		testDonorEmail        = "test@twreporter.org"
		testOrderNumber       = "ValidOrderNumber"
		testRecTradeID        = "ValidRecTradeID"
		testBankTransactionID = "ValidBankTransactionID"
		testAcquirer          = "Test Acquirer"

		oldBankResultMsg  = "Old"
		newBankResultMsg  = "New"
		oldBankResultCode = "Old"
		newBankResultCode = "New"
		errBankResultMsg  = "Bank Error Msg"
		errBankResultCode = "Bank Error Code"

		notifyPath = "/v1/donations/prime/line-notify"
	)

	startTransactionTime := time.Now()
	endTransactionTime := startTransactionTime.Add(30 * time.Second)

	user := createUser(testDonorEmail)
	defer func() { deleteUser(user) }()
	record := models.PayByPrimeDonation{
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Email: testDonorEmail,
		},
		Currency:    testCurrency,
		UserID:      user.ID,
		OrderNumber: testOrderNumber,
		PayMethod:   linePayMethod,
		Status:      statusPaying,
		TappayResp: models.TappayResp{
			RecTradeID:        testRecTradeID,
			BankTransactionID: testBankTransactionID,
		},
	}

	type resultPatchField struct {
		Method          string `json:"method"`
		LastFour        string `json:"last_four"`
		Point           int64  `json:"point"`
		Status          string `json:"status"`
		BankResultCode  string `json:"bank_result_code"`
		BankResultMsg   string `json:"bank_result_msg"`
		TappayApiStatus int64
	}

	db := Globs.GormDB
	cases := []struct {
		name          string
		preRecord     *models.PayByPrimeDonation
		reqBody       tapPayRequestBody
		resultCode    int
		resultCompare *resultPatchField
	}{
		{
			name: "StatusCode=StatusBadRequest,Invalid Line Pay Method",
			reqBody: tapPayRequestBody{
				PayInfo: models.PayInfo{
					Method: null.StringFrom("Invalid Method"),
				},
			},
			resultCode: http.StatusBadRequest,
		},
		{
			name:      "StatusCode=StatusUnprocessableEntity,Unknown OrderNumber",
			preRecord: &record,
			reqBody: tapPayRequestBody{
				OrderNumber: "UnknownOrderNumber",
				PayInfo: models.PayInfo{
					Method:                 null.StringFrom("CREDIT_CARD"),
					MaskedCreditCardNumber: null.StringFrom("************5566"),
					Point:                  null.IntFrom(0),
				},
			},
			resultCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "StatusCode=StatusUnprocessableEntity,Unknown RecTradeID",
			preRecord: &record,
			reqBody: tapPayRequestBody{
				RecTradeID: "UnknownRecTradeID",
				PayInfo: models.PayInfo{
					Method:                 null.StringFrom("CREDIT_CARD"),
					MaskedCreditCardNumber: null.StringFrom("************5566"),
					Point:                  null.IntFrom(0),
				},
			},
			resultCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "StatusCode=StatusUnprocessableEntity,Unknown BankTransactionID",
			preRecord: &record,
			reqBody: tapPayRequestBody{
				BankTransactionID: "UnknownBankTransactionID",
				PayInfo: models.PayInfo{
					Method:                 null.StringFrom("CREDIT_CARD"),
					MaskedCreditCardNumber: null.StringFrom("************5566"),
					Point:                  null.IntFrom(0),
				},
			},
			resultCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "StatusCode=StatusNoContent,Success linepay info using credit card",
			preRecord: &record,
			reqBody: tapPayRequestBody{
				RecTradeID:        testRecTradeID,
				BankTransactionID: testBankTransactionID,
				OrderNumber:       testOrderNumber,
				Amount:            testAmount,
				Status:            0,
				TransactionTime:   endTransactionTime.Unix() * 1000, //millisecond
				PayInfo: models.PayInfo{
					Method:                 null.StringFrom("CREDIT_CARD"),
					MaskedCreditCardNumber: null.StringFrom("************5566"),
					Point:                  null.IntFrom(0),
				},
				Acquirer:       testAcquirer,
				BankResultMsg:  newBankResultMsg,
				BankResultCode: newBankResultCode,
			},
			resultCode: http.StatusNoContent,
			resultCompare: &resultPatchField{
				Method:          "CREDIT_CARD",
				LastFour:        "5566",
				Point:           0,
				Status:          statusPaid,
				BankResultMsg:   newBankResultMsg,
				BankResultCode:  newBankResultCode,
				TappayApiStatus: 0,
			},
		},
		{
			name:      "StatusCode=StatusNoContent,Success linepay info using balance but with redundant credit card number",
			preRecord: &record,
			reqBody: tapPayRequestBody{
				RecTradeID:        testRecTradeID,
				BankTransactionID: testBankTransactionID,
				OrderNumber:       testOrderNumber,
				Amount:            testAmount,
				Status:            0,
				TransactionTime:   endTransactionTime.Unix() * 1000, //millisecond
				PayInfo: models.PayInfo{
					Method:                 null.StringFrom("BALANCE"),
					MaskedCreditCardNumber: null.StringFrom("************5566"),
					Point:                  null.IntFrom(0),
				},
				Acquirer:       testAcquirer,
				BankResultMsg:  newBankResultMsg,
				BankResultCode: newBankResultCode,
			},
			resultCode: http.StatusNoContent,
			resultCompare: &resultPatchField{
				Method:          "BALANCE",
				Point:           0,
				Status:          statusPaid,
				BankResultMsg:   newBankResultMsg,
				BankResultCode:  newBankResultCode,
				TappayApiStatus: 0,
			},
		},
		{
			name:      "StatusCode=StatusNoContent, Linepay Transaction cancelled",
			preRecord: &record,
			reqBody: tapPayRequestBody{
				RecTradeID:        testRecTradeID,
				BankTransactionID: testBankTransactionID,
				OrderNumber:       testOrderNumber,
				Amount:            testAmount,
				Status:            924,                              // error code for gateway timeout
				TransactionTime:   endTransactionTime.Unix() * 1000, // millisecond
				Acquirer:          testAcquirer,
				BankResultMsg:     errBankResultMsg,
				BankResultCode:    errBankResultCode,
			},
			resultCode: http.StatusNoContent,
			resultCompare: &resultPatchField{
				Status:          statusFail,
				BankResultMsg:   errBankResultMsg,
				BankResultCode:  errBankResultCode,
				TappayApiStatus: 924,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var reqBodyInBytes []byte

			if c.preRecord != nil {
				db.Model(&c.preRecord).Create(c.preRecord)

				defer func() {
					db.Unscoped().Delete(c.preRecord)
				}()
			}

			reqBodyInBytes, _ = json.Marshal(&c.reqBody)

			resp := serveHTTP("POST", notifyPath, string(reqBodyInBytes), "application/json", "")

			assert.Equal(t, c.resultCode, resp.Code)

			// Check if the transaction information updated correctly
			if c.resultCode == http.StatusNoContent {
				m := models.PayByPrimeDonation{}
				db.Where("order_number = ?", c.preRecord.OrderNumber).Find(&m)
				assert.Equal(t, c.resultCompare.Method, m.PayInfo.Method.String)
				assert.Equal(t, c.resultCompare.LastFour, m.CardInfo.LastFour.String)
				assert.Equal(t, c.resultCompare.Point, m.PayInfo.Point.Int64)
				assert.Equal(t, c.resultCompare.Status, m.Status)
				assert.Equal(t, c.resultCompare.BankResultMsg, m.BankResultMsg.String)
				assert.Equal(t, c.resultCompare.BankResultCode, m.BankResultCode.String)
				assert.Equal(t, c.resultCompare.TappayApiStatus, m.TappayApiStatus.Int64)
			}
		})
	}
}

func TestGetVerificationOfATransaction(t *testing.T) {
	user := createUser("testDonor@twreporter.org")
	defer func() { deleteUser(user) }()

	maliciousUser := createUser("testMaliciousDonor@twreporter.org")
	defer func() { deleteUser(maliciousUser) }()

	authorization, cookie := helperSetupAuth(user)
	maliciousAuthorization, maliciousCookie := helperSetupAuth(maliciousUser)

	record := models.PayByPrimeDonation{
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Email: user.Email.String,
		},
		Currency:    testCurrency,
		UserID:      user.ID,
		OrderNumber: "ValidOrderNumber1",
		PayMethod:   linePayMethod,
		Status:      "paid",
		TappayResp: models.TappayResp{
			RecTradeID:        "ValidRecTradeID1",
			BankTransactionID: "ValidBankTransactionID1",
			TappayApiStatus:   null.IntFrom(0),
		},
	}

	failRecord := models.PayByPrimeDonation{
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Email: user.Email.String,
		},
		Currency:    testCurrency,
		UserID:      user.ID,
		OrderNumber: "ValidOrderNumber1",
		PayMethod:   linePayMethod,
		Status:      "fail",
		TappayResp: models.TappayResp{
			RecTradeID:        "ValidRecTradeID1",
			BankTransactionID: "ValidBankTransactionID1",
			TappayApiStatus:   null.IntFrom(421), // Gateway Timeout Error
		},
	}

	type (
		verificationData struct {
			RecTradeID        string `json:"rec_trade_id"`
			BankTransactionID string `json:"bank_transaction_id"`
			Status            string `json:"status"`
		}
		verificationResp struct {
			Status string           `json:"status"`
			Data   verificationData `json:"data"`
		}
	)

	db := Globs.GormDB
	cases := []struct {
		reqHeader
		name          string
		preRecord     *models.PayByPrimeDonation
		orderNumber   string
		resultCode    int
		resultCompare *verificationResp
	}{
		{
			name: "StatusCode=StatusUnauthorized,Lack of Authorization Header",
			reqHeader: reqHeader{
				Cookie: &cookie,
			},
			orderNumber: record.OrderNumber,
			resultCode:  http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusForbidden,Unauthorized Resource",
			reqHeader: reqHeader{
				Cookie:        &maliciousCookie,
				Authorization: maliciousAuthorization,
			},
			preRecord:   &record,
			orderNumber: record.OrderNumber,
			resultCode:  http.StatusForbidden,
		},
		{
			name: "StatusCode=StatusNotFound,Invalid Order Number",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			preRecord:   &record,
			orderNumber: "InvalidOrderNumber",
			resultCode:  http.StatusNotFound,
		},
		{
			name: "StatusCode=StatusOK,Transaction Success",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			preRecord:   &record,
			orderNumber: record.OrderNumber,
			resultCode:  http.StatusOK,
			resultCompare: &verificationResp{
				Status: "success",
				Data: verificationData{
					RecTradeID:        record.RecTradeID,
					BankTransactionID: record.BankTransactionID,
					Status:            statusPaid,
				},
			},
		},
		{
			name: "StatusCode=StatusOK,Transaction fail:gateway timeout ",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			preRecord:   &failRecord,
			orderNumber: failRecord.OrderNumber,
			resultCode:  http.StatusOK,
			resultCompare: &verificationResp{
				Status: "success",
				Data: verificationData{
					RecTradeID:        failRecord.RecTradeID,
					BankTransactionID: failRecord.BankTransactionID,
					Status:            statusFail,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.preRecord != nil {
				db.Model(c.preRecord).Create(c.preRecord)

				defer func() {
					db.Unscoped().Delete(c.preRecord)
				}()
			}

			path := fmt.Sprintf("/v1/donations/prime/orders/%s/transaction_verification", c.orderNumber)
			resp := serveHTTPWithCookies("GET", path, "", "application/json", c.reqHeader.Authorization, *c.reqHeader.Cookie)
			assert.Equal(t, c.resultCode, resp.Code)

			if c.resultCompare != nil {
				expect, _ := json.Marshal(c.resultCompare)
				respBodyInBytes, _ := ioutil.ReadAll(resp.Result().Body)
				assert.JSONEq(t, string(expect), string(respBodyInBytes))
			}
		})
	}

}

func TestQueryTappayServer(t *testing.T) {
	type (
		filter struct {
			OrderNumber string `json:"order_number"`
		}

		requestBody struct {
			RecordsPerPage uint   `json:"records_per_page"`
			Filters        filter `json:"filters"`
		}

		tradeRecord struct {
			RecordStatus int `json:"record_status"`
		}

		tappayRecord struct {
			Status       int           `json:"status"`
			Msg          string        `json:"msg"`
			TradeRecords []tradeRecord `json:"trade_records"`
		}

		responseBody struct {
			Status string       `json:"status"`
			Data   tappayRecord `json:"data"`
		}
	)

	user := createUser("testDonorEmailr@twreporter.org")
	defer func() { deleteUser(user) }()
	maliciousUser := createUser("testMaliciousDonor@twreporter.org")
	defer func() { deleteUser(maliciousUser) }()

	authorization, cookie := helperSetupAuth(user)
	maliciousAuthorization, maliciousCookie := helperSetupAuth(maliciousUser)

	dbRecord := models.PayByPrimeDonation{
		Amount: testAmount,
		Cardholder: models.Cardholder{
			Email: user.Email.String,
		},
		Currency:    testCurrency,
		UserID:      user.ID,
		OrderNumber: "ValidOrderNumber1",
		PayMethod:   linePayMethod,
		Status:      statusPaying,
		TappayResp: models.TappayResp{
			RecTradeID:        "ValidRecTradeID1",
			BankTransactionID: "ValidBankTransactionID1",
			TappayApiStatus:   null.IntFrom(0),
		},
	}

	transactionSuccessRecord := tappayRecord{
		Status: 0,
		Msg:    "",
		TradeRecords: []tradeRecord{
			tradeRecord{
				RecordStatus: 0,
			},
		},
	}

	transactionFailRecord := tappayRecord{
		Status: 0,
		Msg:    "",
		TradeRecords: []tradeRecord{
			tradeRecord{
				RecordStatus: -1,
			},
		},
	}

	queryFailRecord := tappayRecord{
		Status: 421,
		Msg:    "Gateway timeout",
	}

	cases := []struct {
		reqHeader
		name             string
		reqBody          *requestBody
		preRecord        *models.PayByPrimeDonation
		stubTappayServer *httptest.Server
		resultCode       int
		resultCompare    *tappayRecord
	}{
		{
			name: "StatusCode=StatusUnauthorized,Lack of Authorization Header",
			reqHeader: reqHeader{
				Cookie: &cookie,
			},
			resultCode: http.StatusUnauthorized,
		},
		{
			name: "StatusCode=StatusForbidden,Unauthorized Resource",
			reqHeader: reqHeader{
				Cookie:        &maliciousCookie,
				Authorization: maliciousAuthorization,
			},
			reqBody: &requestBody{
				Filters: filter{
					OrderNumber: "ValidOrderNumber1",
				},
			},
			preRecord:  &dbRecord,
			resultCode: http.StatusForbidden,
		},
		{
			name: "StatusCode=StatusNotFound,Invalid Order Number",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			reqBody: &requestBody{
				Filters: filter{
					OrderNumber: "Invalid Order Number",
				},
			},
			preRecord:  &dbRecord,
			resultCode: http.StatusNotFound,
		},
		{
			name: "StatusCode=StatusOK,Query Success&Transaction Success",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			reqBody: &requestBody{
				Filters: filter{
					OrderNumber: "ValidOrderNumber1",
				},
			},
			preRecord: &dbRecord,
			stubTappayServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				resp, _ := json.Marshal(transactionSuccessRecord)
				w.Write(resp)
			})),
			resultCode:    http.StatusOK,
			resultCompare: &transactionSuccessRecord,
		},
		{
			name: "StatusCode=StatusOK,Query Success&Transaction Fail",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},
			reqBody: &requestBody{
				Filters: filter{
					OrderNumber: "ValidOrderNumber1",
				},
			},
			preRecord: &dbRecord,
			stubTappayServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				resp, _ := json.Marshal(transactionFailRecord)
				w.Write(resp)
			})),
			resultCode:    http.StatusOK,
			resultCompare: &transactionFailRecord,
		},
		{
			name: "StatusCode=StatusOK,Query Fail",
			reqHeader: reqHeader{
				Cookie:        &cookie,
				Authorization: authorization,
			},

			reqBody: &requestBody{
				Filters: filter{
					OrderNumber: "ValidOrderNumber1",
				},
			},
			preRecord: &dbRecord,
			stubTappayServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				resp, _ := json.Marshal(queryFailRecord)
				w.Write(resp)
			})),
			resultCode:    http.StatusOK,
			resultCompare: &queryFailRecord,
		},
	}

	db := Globs.GormDB
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.preRecord != nil {
				db.Model(c.preRecord).Create(c.preRecord)

				defer func() {
					db.Unscoped().Delete(c.preRecord)
				}()
			}

			// Stub out tappay server if the request would be sent
			if c.stubTappayServer != nil {
				defer c.stubTappayServer.Close()

				// Overwrite the tappay record server to stub server
				url := globals.Conf.Donation.TapPayRecordURL
				globals.Conf.Donation.TapPayRecordURL = c.stubTappayServer.URL
				defer func() {
					globals.Conf.Donation.TapPayRecordURL = url
				}()
			}

			path := "/v1/tappay_query"

			var body []byte
			if c.reqBody != nil {
				body, _ = json.Marshal(c.reqBody)
			}

			resp := serveHTTPWithCookies(http.MethodPost, path, string(body), "application/json", c.reqHeader.Authorization, *c.reqHeader.Cookie)
			assert.Equal(t, c.resultCode, resp.Code)

			if c.resultCompare != nil {
				bodyJSON, _ := ioutil.ReadAll(resp.Body)
				var body responseBody
				json.Unmarshal(bodyJSON, &body)
				assert.Exactly(t, *c.resultCompare, body.Data)
			}
		})
	}
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
	jwt := generateIDToken(user)
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
	jwt = generateIDToken(models.User{
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
