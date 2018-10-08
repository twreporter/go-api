package controllers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gopkg.in/go-playground/validator.v8"
	"gopkg.in/guregu/null.v3"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

type (
	clientReq struct {
		Prime       string            `json:"prime" form:"prime" binding:"required"`
		Amount      uint              `json:"amount" form:"amount" binding:"required"`
		Currency    string            `json:"currency" form:"currency"`
		Details     string            `json:"details" form:"details"`
		Cardholder  models.Cardholder `json:"cardholder" form:"cardholder" binding:"required,dive"`
		OrderNumber string            `json:"order_number" form:"order_number"`
		MerchantID  string            `json:"merchant_id" form:"merchant_id"`
		ResultUrl   linePayResultUrl  `json:"result_url" form:"result_url"`
	}

	clientResp struct {
		PeriodicID  null.Int          `json:"periodic_id"`
		PayMethod   string            `json:"pay_method"`
		CardInfo    models.CardInfo   `json:"card_info"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Amount      uint              `json:"amount"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		OrderNumber string            `json:"order_number"`
		ID          uint              `json:"id"`
	}

	bankTransactionTime struct {
		StartTimeMillis string `json:"start_time_millis"`
		EndTimeMillis   string `json:"end_time_millis"`
	}

	cardSecret struct {
		CardToken string `json:"card_token"`
		CardKey   string `json:"card_key"`
	}

	linePayResultUrl struct {
		FrontendRedirectUrl string `json:"frontend_redirect_url" form:"frontend_redirect_url"`
		BackendNotifyUrl    string `json:"backend_notify_url" form:"backend_notify_url"`
	}

	tapPayPrimeReq struct {
		Prime       string            `json:"prime"`
		PartnerKey  string            `json:"partner_key"`
		MerchantID  string            `json:"merchant_id"`
		Amount      uint              `json:"amount"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Remember    bool              `json:"remember"`
		OrderNumber string            `json:"order_number"`
		ResultUrl   linePayResultUrl  `json:"result_url"`
	}

	tapPayTransactionResp struct {
		Status                int64               `json:"status"`
		Msg                   string              `json:"msg"`
		RecTradeID            string              `json:"rec_trade_id"`
		BankTransactionID     string              `json:"bank_transaction_id"`
		AuthCode              string              `json:"auth_code"`
		CardSecret            cardSecret          `json:"card_secret"`
		Amount                uint                `json:"amount"`
		Currency              string              `json:"currency"`
		CardInfo              models.CardInfo     `json:"card_info"`
		OrderNumber           string              `json:"order_number"`
		Acquirer              string              `json:"acquirer"`
		TransactionTimeMillis int64               `json:"transaction_time_millis"`
		BankTransactionTime   bankTransactionTime `json:"bank_transaction_time"`
		BankResultCode        null.String         `json:"bank_result_code"`
		BankResultMsg         null.String         `json:"bank_result_msg"`
	}

	tapPayMinTransactionResp struct {
		Status int64  `json:"status"`
		Msg    string `json:"msg"`
	}

	payType int

	patchBody struct {
		Address     null.String `json:"address"`
		Details     null.String `json:"details"`
		Name        null.String `json:"name"`
		NationalID  null.String `json:"national_id"`
		PhoneNumber null.String `json:"phone_number"`
		SendReceipt null.String `json:"send_receipt"`
		ToFeedback  null.Bool   `json:"to_feedback"`
		ZipCode     null.String `json:"zip_code"`
	}
)

// CardholderName method is used by copier
// which will copy `Name` to `CardholderName`
func (p *patchBody) CardholderName() null.String {
	return p.Name
}

// CardholderNationID method is used by copier
// which will copy `NationalID` to `CardholderNationalID`
func (p *patchBody) CardholderNationalID() null.String {
	return p.NationalID
}

// CardholderPhoneNumber method is used by copier
// which will copy `PhoneNumber` to `CardholderPhoneNumber`
func (p *patchBody) CardholderPhoneNumber() null.String {
	return p.PhoneNumber
}

// CardholderZipCode method is used by copier
// which will copy `ZipCode` to `CardholderZipCode`
func (p *patchBody) CardholderZipCode() null.String {
	return p.ZipCode
}

// CardholderAddress method is used by copier
// which will copy `Address` to `CardholderAddress`
func (p *patchBody) CardholderAddress() null.String {
	return p.Address
}

const (
	defaultDetails        = "報導者小額捐款"
	defaultCurrency       = "TWD"
	defaultMerchantID     = "twreporter_CTBC"
	defaultRequestTimeout = 45 * time.Second

	invalidPayMethodID = -1

	orderPrefix = "twreporter"

	statusPaying = "paying"
	statusPaid   = "paid"
	statusFail   = "fail"

	statusPeriodicToPay  = "to_pay"
	statusPeriodicPaying = "paying"
	statusPeriodicPaid   = "paid"
	statusPeriodicFail   = "fail"

	primeTableName = "pay_by_prime_donations"

	tapPayRespStatusSuccess = 0

	defaultPeriodicPayMethod = "credit_card"

	secToMsec     = 1000
	msecToNanosec = 1000000
)

// pay type Enum
const (
	oneTime payType = iota
	periodic
)

// pay method collections
var payMethodCollections = []string{
	"credit_card",
	"line",
	"google",
	"apple",
	"samsung",
}

func bindRequestBody(c *gin.Context, reqBody interface{}) (gin.H, bool) {
	// Validate request body
	if err := c.Bind(reqBody); nil != err {
		log.Error("parse model error: " + err.Error())
		failData := gin.H{}

		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			failData[e.Field] = fmt.Sprintf("Cannot unmarshal %s into %s", e.Value, e.Type)
		case validator.ValidationErrors:
			for _, errField := range e {
				failData[errField.Name] = "cannot be empty"
			}
		default:
			// Omit intentionally
		}
		return failData, false
	}

	return gin.H{}, true
}

// Handler for an authenticated user to create a periodic donation
func (mc *MembershipController) CreateAPeriodicDonationOfAUser(c *gin.Context) (int, gin.H, error) {
	const errWhere = "MembershipController.CreateAPeriodicDonationOfAUser"

	// Validate client request
	var reqBody clientReq

	if failData, valid := bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	userID, _ := strconv.ParseUint(c.Param("userID"), 10, strconv.IntSize)

	// Build a draft periodic donation record
	draftPeriodicDonation := buildDraftPeriodicDonation(uint(userID), reqBody)

	// Build Tappay prime request
	tapPayReq := buildTapPayPrimeReq(periodic, defaultPeriodicPayMethod, reqBody)

	draftRecord := buildTokenDraftRecord(tapPayReq)

	// Create a draft periodic donation along with the first token donation record of that periodic donation
	err := mc.Storage.CreateAPeriodicDonation(&draftPeriodicDonation, &draftRecord)
	if nil != err {
		errMsg := "Unable to create a draft periodic donation and the first card token transaction record"
		log.Error(fmt.Sprintf("%s: %s", errWhere, errMsg))
		return 0, gin.H{}, models.NewAppError(errWhere, errMsg, err.Error(), http.StatusInternalServerError)
	}

	// Start Tappay transaction
	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err := serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'fail' and stop the periodic donation
			failResp := models.PayByCardTokenDonation{
				TappayApiStatus: null.IntFrom(tapPayResp.Status),
				Msg:             tapPayResp.Msg,
				Status:          statusFail,
			}
			// Procceed even if the deletion is failed
			mc.Storage.DeleteAPeriodicDonation(draftPeriodicDonation.ID, failResp)
		}
		errMsg := err.Error()
		log.Error(fmt.Sprintf("%s: %s", errWhere, errMsg))

		return 0, gin.H{}, models.NewAppError(errWhere, errMsg, "", http.StatusInternalServerError)
	}

	updateRecord := buildTokenSuccessRecord(tapPayResp)
	updatePeriodicDonation := buildSuccessPeriodicDonation(tapPayResp)
	updateRecord.ID = draftRecord.ID

	if err = mc.Storage.UpdatePeriodicAndCardTokenDonationInTRX(draftPeriodicDonation.ID, updatePeriodicDonation, updateRecord); nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
	}

	resp := buildClientResp(defaultPeriodicPayMethod, tapPayReq, tapPayResp)
	resp.PeriodicID = null.IntFrom(int64(draftPeriodicDonation.ID))
	resp.ID = draftRecord.ID

	return http.StatusCreated, gin.H{"status": "success", "data": resp}, nil
}

// Handler for an authenticated user to create an one-time donation
func (mc *MembershipController) CreateADonationOfAUser(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.CreateADonationOfAUser"

	// Validate client request

	// Validate pay_method
	payMethod := c.Param("pay_method")
	if err := validatePayMethod(payMethod); nil != err {
		return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"pay_method": err.Error()}}, nil
	}

	var reqBody clientReq

	if failData, valid := bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	userID, _ := strconv.ParseUint(c.Param("userID"), 10, strconv.IntSize)

	// Start Tappay transaction

	// Build Tappay pay by prime request
	tapPayReq := buildTapPayPrimeReq(oneTime, payMethod, reqBody)

	draftRecord := buildPrimeDraftRecord(uint(userID), payMethod, tapPayReq)

	if err := mc.Storage.Create(&draftRecord); nil != err {
		switch appErr := err.(type) {
		case *models.AppError:
			return 0, gin.H{}, models.NewAppError(errorWhere, "Fails to create a draft prime record", appErr.Error(), appErr.StatusCode)
		default:
			return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Unknown Error Type. Fails to create a draft prime record. %v", err.Error())}, nil
		}
	}

	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err := serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'fail'
			mc.Storage.UpdateByConditions(map[string]interface{}{
				"order_number": tapPayReq.OrderNumber,
			}, models.PayByPrimeDonation{
				TappayApiStatus: null.IntFrom(tapPayResp.Status),
				Msg:             tapPayResp.Msg,
				Status:          statusFail,
			})
		}
		return 0, gin.H{}, models.NewAppError(errorWhere, err.Error(), "", http.StatusInternalServerError)
	}

	// Update the transaction status to 'paid' if transaction succeeds
	updateRecord := buildPrimeSuccessRecord(tapPayResp)

	if err, _ := mc.Storage.UpdateByConditions(map[string]interface{}{
		"order_number": tapPayReq.OrderNumber,
	}, updateRecord); nil != err {
		log.Error(err.Error())
	}

	resp := buildClientResp(payMethod, tapPayReq, tapPayResp)
	resp.ID = draftRecord.ID

	return http.StatusCreated, gin.H{"status": "success", "data": resp}, nil
}

// PatchADonationOfAUser method
// Handler for an authenticated user to patch an prime/token/periodic donation
func (mc *MembershipController) PatchADonationOfAUser(c *gin.Context, donationType string) (int, gin.H, error) {
	var d interface{}
	var err error
	var failData gin.H
	var recordID uint64
	var reqBody patchBody
	var rowsEffected int64
	var userID uint64
	var valid bool

	if userID, err = strconv.ParseUint(c.Param("userID"), 10, strconv.IntSize); err != nil {
		return http.StatusNotFound, gin.H{"status": "error", "message": "record not found, user id should be provided in the url"}, err
	}

	if recordID, err = strconv.ParseUint(c.Param("id"), 10, strconv.IntSize); err != nil {
		return http.StatusNotFound, gin.H{"status": "error", "message": "record not found, record id should be provided in the url"}, err
	}

	if failData, valid = bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	switch donationType {
	case globals.PeriodicDonationType:
		pd := models.PeriodicDonation{}
		copier.Copy(&pd, &reqBody)
		d = pd
	case globals.PrimeDonaitionType:
		pd := models.PayByPrimeDonation{}
		copier.Copy(&pd, &reqBody)
		d = pd
	default:
		return http.StatusInternalServerError,
			gin.H{"status": "error", "message": fmt.Sprintf("donation type(%s) not supported", donationType)},
			nil
	}

	if err, rowsEffected = mc.Storage.UpdateByConditions(map[string]interface{}{
		"user_id": userID,
		"id":      recordID,
	}, d); err != nil {
		return 0, gin.H{}, err
	}

	if rowsEffected == 0 {
		return http.StatusNotFound, gin.H{"status": "error", "message": fmt.Sprintf("record(id: %d) not found", recordID)}, nil
	}

	return http.StatusNoContent, gin.H{}, nil
}

//TODO
func (mc *MembershipController) GetDonationsOfAUser(c *gin.Context) (int, gin.H, error) {
	return 0, gin.H{}, nil
}

func (t *tapPayPrimeReq) setDefault() {
	t.Details = defaultDetails
	t.Currency = defaultCurrency
	t.MerchantID = defaultMerchantID
}

func buildClientResp(payMethod string, req tapPayPrimeReq, resp tapPayTransactionResp) clientResp {
	c := clientResp{}
	c.PayMethod = payMethod

	c.CardInfo = resp.CardInfo
	c.Cardholder = req.Cardholder
	c.OrderNumber = req.OrderNumber

	c.Amount = resp.Amount
	c.Currency = resp.Currency
	c.Details = req.Details

	return c
}

func buildDraftPeriodicDonation(userID uint, req clientReq) models.PeriodicDonation {
	m := models.PeriodicDonation{}

	m.CardholderEmail = req.Cardholder.Email
	m.CardholderPhoneNumber = req.Cardholder.PhoneNumber
	m.CardholderName = req.Cardholder.Name
	m.CardholderZipCode = req.Cardholder.ZipCode
	m.CardholderAddress = req.Cardholder.Address
	m.CardholderNationalID = req.Cardholder.NationalID

	m.UserID = userID
	m.Status = statusPeriodicPaying

	return m
}

func buildPrimeDraftRecord(userID uint, payMethod string, req tapPayPrimeReq) models.PayByPrimeDonation {
	m := models.PayByPrimeDonation{}

	m.UserID = userID
	m.PayMethod = payMethod

	m.CardholderEmail = req.Cardholder.Email
	m.CardholderPhoneNumber = req.Cardholder.PhoneNumber
	m.CardholderName = req.Cardholder.Name
	m.CardholderZipCode = req.Cardholder.ZipCode
	m.CardholderAddress = req.Cardholder.Address
	m.CardholderNationalID = req.Cardholder.NationalID

	m.Amount = req.Amount
	m.Details = req.Details
	m.MerchantID = req.MerchantID
	m.OrderNumber = req.OrderNumber

	m.Status = statusPaying

	return m
}

func buildPrimeSuccessRecord(resp tapPayTransactionResp) models.PayByPrimeDonation {
	m := models.PayByPrimeDonation{}

	m.TappayApiStatus = null.IntFrom(resp.Status)
	m.Msg = resp.Msg
	m.RecTradeID = resp.RecTradeID
	m.BankTransactionID = resp.BankTransactionID
	m.AuthCode = resp.AuthCode
	m.Acquirer = resp.Acquirer
	m.BankResultCode = resp.BankResultCode
	m.BankResultMsg = resp.BankResultMsg

	ttm := time.Unix(resp.TransactionTimeMillis/secToMsec, (resp.TransactionTimeMillis%secToMsec)*msecToNanosec)
	m.TransactionTime = null.TimeFrom(ttm)

	t, err := strconv.ParseInt(resp.BankTransactionTime.StartTimeMillis, 10, strconv.IntSize)
	if nil == err {
		stm := time.Unix(t/secToMsec, t%secToMsec)
		m.BankTransactionStartTime = null.TimeFrom(stm)
	}

	t, err = strconv.ParseInt(resp.BankTransactionTime.EndTimeMillis, 10, strconv.IntSize)
	if nil == err {
		etm := time.Unix(t/secToMsec, t%secToMsec)
		m.BankTransactionEndTime = null.TimeFrom(etm)
	}

	m.CardInfoBinCode = resp.CardInfo.BinCode
	m.CardInfoLastFour = resp.CardInfo.LastFour
	m.CardInfoIssuer = resp.CardInfo.Issuer
	m.CardInfoFunding = resp.CardInfo.Funding
	m.CardInfoType = resp.CardInfo.Type
	m.CardInfoLevel = resp.CardInfo.Level
	m.CardInfoCountry = resp.CardInfo.Country
	m.CardInfoCountryCode = resp.CardInfo.CountryCode
	m.CardInfoExpiryDate = resp.CardInfo.ExpiryDate

	m.Status = "paid"

	return m
}

func buildSuccessPeriodicDonation(resp tapPayTransactionResp) models.PeriodicDonation {
	m := models.PeriodicDonation{}

	m.CardInfoBinCode = resp.CardInfo.BinCode
	m.CardInfoLastFour = resp.CardInfo.LastFour
	m.CardInfoIssuer = resp.CardInfo.Issuer
	m.CardInfoFunding = resp.CardInfo.Funding
	m.CardInfoType = resp.CardInfo.Type
	m.CardInfoLevel = resp.CardInfo.Level
	m.CardInfoCountry = resp.CardInfo.Country
	m.CardInfoCountryCode = resp.CardInfo.CountryCode
	m.CardInfoExpiryDate = resp.CardInfo.ExpiryDate

	var ciphertext string

	ciphertext = encrypt(resp.CardSecret.CardToken, globals.Conf.Donation.CardSecretKey)
	m.CardToken = ciphertext

	ciphertext = encrypt(resp.CardSecret.CardKey, globals.Conf.Donation.CardSecretKey)
	m.CardKey = ciphertext

	t := time.Now()
	m.LastSuccessAt = null.TimeFrom(t)
	m.Status = statusPeriodicPaid

	return m
}

func buildTapPayPrimeReq(pt payType, payMethod string, req clientReq) tapPayPrimeReq {
	t := &tapPayPrimeReq{}

	t.setDefault()

	// Fill up required fields
	t.Prime = req.Prime
	t.Amount = req.Amount
	t.Cardholder = req.Cardholder

	// Per required fields (even empty) of cardholder of tappay documents,
	// use empty strings for name and phonenumber fields instead of empty.
	if !t.Cardholder.Name.Valid {
		t.Cardholder.Name = null.StringFrom("")
	}

	if !t.Cardholder.PhoneNumber.Valid {
		t.Cardholder.PhoneNumber = null.StringFrom("")
	}

	// Fill up optional fields
	if "" != req.Currency {
		t.Currency = req.Currency
	}

	if "" != req.Details {
		t.Details = req.Details
	}

	if "" != req.MerchantID {
		t.MerchantID = req.MerchantID
	}

	if (linePayResultUrl{}) != req.ResultUrl {
		t.ResultUrl = req.ResultUrl
	}

	if "" != req.OrderNumber {
		t.OrderNumber = req.OrderNumber
	} else {
		t.OrderNumber = generateOrderNumber(pt, getPayMethodID(payMethod))
	}
	t.PartnerKey = globals.Conf.Donation.TapPayPartnerKey

	if periodic == pt {
		t.Remember = true
	}

	return *t
}

func buildTokenSuccessRecord(resp tapPayTransactionResp) models.PayByCardTokenDonation {
	m := models.PayByCardTokenDonation{}

	m.TappayApiStatus = null.IntFrom(resp.Status)
	m.Msg = resp.Msg
	m.RecTradeID = resp.RecTradeID
	m.BankTransactionID = resp.BankTransactionID
	m.AuthCode = resp.AuthCode
	m.Acquirer = resp.Acquirer
	m.BankResultCode = resp.BankResultCode
	m.BankResultMsg = resp.BankResultMsg

	ttm := time.Unix(resp.TransactionTimeMillis/secToMsec, (resp.TransactionTimeMillis%secToMsec)*msecToNanosec)
	m.TransactionTime = null.TimeFrom(ttm)

	t, err := strconv.ParseInt(resp.BankTransactionTime.StartTimeMillis, 10, strconv.IntSize)
	if nil == err {
		stm := time.Unix(t/secToMsec, t%secToMsec)
		m.BankTransactionStartTime = null.TimeFrom(stm)
	}

	t, err = strconv.ParseInt(resp.BankTransactionTime.EndTimeMillis, 10, strconv.IntSize)
	if nil == err {
		etm := time.Unix(t/secToMsec, t%secToMsec)
		m.BankTransactionEndTime = null.TimeFrom(etm)
	}

	m.Status = statusPaid
	return m
}

func buildTokenDraftRecord(req tapPayPrimeReq) models.PayByCardTokenDonation {
	m := models.PayByCardTokenDonation{}

	m.Details = req.Details
	m.MerchantID = req.MerchantID
	m.OrderNumber = req.OrderNumber
	m.Amount = req.Amount

	m.Status = statusPaying

	return m
}

func createHash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func decrypt(data string, key string) string {
	hashKey := createHash(key)
	block, _ := aes.NewCipher(hashKey)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()

	byteData := []byte(data)
	nonce, ciphertext := byteData[:nonceSize], byteData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if nil != err {
		log.Error(err.Error())
	}

	return string(plaintext)
}

func encrypt(data string, key string) string {
	hashKey := createHash(key)

	// create a aes block cipher by the hash value of our key
	block, _ := aes.NewCipher(hashKey)

	// use Galois Counter Mode for better efficiency
	gcm, err := cipher.NewGCM(block)
	if nil != err {
		//fallback to store plaintext instead
		log.Error("cannot create a block cipher with the given key")
		return data
	}

	nonce := make([]byte, gcm.NonceSize())

	// generate random nonce for encryption
	if _, err := io.ReadFull(rand.Reader, nonce); nil != err {
		//fallback to store plaintext instead
		log.Error("cannot generate a nonce")
		return data
	}

	// prepend the cipher with the random nonce
	return string(gcm.Seal(nonce, nonce, []byte(data), nil))
}

func generateOrderNumber(t payType, payMethodID int) string {
	timestamp := time.Now().UnixNano()
	orderNumber := fmt.Sprintf("%s-%d%d%d", orderPrefix, timestamp, t, payMethodID)
	msg := fmt.Sprintf("OrderNumber: %s", orderNumber)
	log.Info(msg)
	return orderNumber
}

func getPayMethodID(payMethod string) int {
	for ind, v := range payMethodCollections {
		if v == payMethod {
			return ind
		}
	}
	return invalidPayMethodID
}

func handleTapPayBodyParseError(body []byte) (tapPayTransactionResp, error) {
	var minResp tapPayMinTransactionResp
	var resp tapPayTransactionResp
	var err error

	if err = json.Unmarshal(body, &minResp); nil != err {
		return tapPayTransactionResp{}, errors.New("Cannot unmarshal json response from tap pay server")
	}

	if tapPayRespStatusSuccess != minResp.Status {
		log.Error("tap pay msg: " + minResp.Msg)
		err = errors.New("Cannot make success transaction on tap pay")
	}

	resp.Status = minResp.Status
	resp.Msg = minResp.Msg

	return resp, err
}

func serveHttp(key string, reqBodyJson []byte) (tapPayTransactionResp, error) {
	// Setup HTTP client with timeout
	client := &http.Client{Timeout: defaultRequestTimeout}

	req, _ := http.NewRequest("POST", globals.Conf.Donation.TapPayURL, bytes.NewBuffer(reqBodyJson))
	req.Header.Add("x-api-key", key)
	req.Header.Add("Content-Type", "application/json")

	rawResp, err := client.Do(req)

	// If fail to sending request
	if nil != err {
		log.Error(err.Error())
		return tapPayTransactionResp{}, errors.New("cannot request to tap pay server")
	}
	defer rawResp.Body.Close()

	// If timeout or other errors occur during reading the body...
	// TODO: Might require a mechanism to notify users
	body, err := ioutil.ReadAll(rawResp.Body)
	if nil != err {
		log.Error(err.Error())
		return tapPayTransactionResp{}, errors.New("Cannot read response from tap pay server")
	}

	resp := tapPayTransactionResp{}

	err = json.Unmarshal(body, &resp)

	switch {
	case nil != err:
		log.Error(err.Error())
		return handleTapPayBodyParseError(body)
	case tapPayRespStatusSuccess != resp.Status:
		log.Error("tap pay msg: " + resp.Msg)
		return resp, errors.New("Cannot make success transaction on tap pay")
	default:
		// Omit intentionally
	}

	return resp, nil
}

func validatePayMethod(payMethod string) error {
	if invalidPayMethodID != getPayMethodID(payMethod) {
		return nil
	}

	errMsg := fmt.Sprintf("Unsupported pay_method. Only support payment by %s", strings.Join(payMethodCollections, ", "))
	return errors.New(errMsg)
}
