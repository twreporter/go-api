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
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
	"gopkg.in/go-playground/validator.v8"
	"gopkg.in/guregu/null.v3"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

type (
	clientReq struct {
		Amount      uint              `json:"amount" form:"amount" binding:"required"`
		Cardholder  models.Cardholder `json:"cardholder" form:"cardholder" binding:"required,dive"`
		Currency    string            `json:"currency" form:"currency"`
		Details     string            `json:"details" form:"details"`
		MerchantID  string            `json:"merchant_id" form:"merchant_id"`
		OrderNumber string            `json:"order_number" form:"order_number"`
		Prime       string            `json:"prime" form:"prime" binding:"required"`
		ResultUrl   linePayResultUrl  `json:"result_url" form:"result_url"`
		UserID      uint              `json:"user_id" form:"user_id" binding:"required"`
	}

	clientResp struct {
		Amount      uint              `json:"amount"`
		CardInfo    models.CardInfo   `json:"card_info,omitempty"`
		Cardholder  models.Cardholder `json:"cardholder,omitempty"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details,omitempty"`
		ID          uint              `json:"id"`
		Notes       string            `json:"notes"`
		OrderNumber string            `json:"order_number,omitempty"`
		PayMethod   string            `json:"pay_method"`
		PeriodicID  uint              `json:"periodic_id,omitempty"`
		SendReceipt string            `json:"send_receipt"`
		ToFeedback  bool              `json:"to_feedback"`
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
		Notes       null.String `json:"notes"`
		PhoneNumber null.String `json:"phone_number"`
		SendReceipt null.String `json:"send_receipt"`
		ToFeedback  null.Bool   `json:"to_feedback"`
		UserID      uint        `json:"user_id" binding:"required"`
		ZipCode     null.String `json:"zip_code"`
	}
)

func (cr *clientResp) BuildFromPeriodicDonationModel(pd models.PeriodicDonation) {
	cardInfo := models.CardInfo{}
	cardholder := models.Cardholder{}

	copier.Copy(&cardInfo, &pd)
	copier.Copy(&cardholder, &pd)
	copier.Copy(cr, &pd)
	cr.Cardholder = cardholder
	cr.CardInfo = cardInfo
	cr.PayMethod = defaultPeriodicPayMethod
}

func (cr *clientResp) BuildFromPrimeDonationModel(pd models.PayByPrimeDonation) {
	cardInfo := models.CardInfo{}
	cardholder := models.Cardholder{}

	copier.Copy(&cardInfo, &pd)
	copier.Copy(&cardholder, &pd)
	copier.Copy(cr, &pd)
	cr.Cardholder = cardholder
	cr.CardInfo = cardInfo
}

func (cr *clientResp) BuildFromOtherMethodDonationModel(d models.PayByOtherMethodDonation) {
	copier.Copy(cr, &d)
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

var cardInfoTypes = map[int64]string{
	1: "VISA",
	2: "MasterCard",
	3: "JCB",
	4: "Union Pay",
	5: "AMEX",
}

func bindRequestBody(c *gin.Context, reqBody interface{}) (gin.H, bool) {
	var err error
	// Validate request body
	// gin.Context.Bind does not support to bind `JSON` body multiple times
	// the alternative is to use gin.Context.ShouldBindBodyWith function to bind
	if err = c.ShouldBindBodyWith(reqBody, binding.JSON); nil != err {

		// bind other format rather than JSON
		if err = c.Bind(reqBody); nil != err {
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
	}

	return gin.H{}, true
}

func (mc *MembershipController) sendDonationThankYouMail(body clientResp, donationType string) {
	reqBody := donationSuccessReqBody{
		Address:          body.Cardholder.Address.ValueOrZero(),
		Amount:           body.Amount,
		CardInfoLastFour: body.CardInfo.LastFour.ValueOrZero(),
		CardInfoType:     cardInfoTypes[body.CardInfo.Type.ValueOrZero()],
		Currency:         body.Currency,
		DonationMethod:   "信用卡支付",
		DonationType:     donationType,
		Email:            body.Cardholder.Email,
		Name:             body.Cardholder.Name.ValueOrZero(),
		OrderNumber:      body.OrderNumber,
		NationalID:       body.Cardholder.NationalID.ValueOrZero(),
		PhoneNumber:      body.Cardholder.PhoneNumber.ValueOrZero(),
	}

	if err := postMailServiceEndpoint(reqBody, fmt.Sprintf("http://localhost:%s/v1/%s", globals.LocalhostPort, globals.SendSuccessDonationRoutePath)); err != nil {
		log.Warnf("fail to send %s donation(order_number: %s) thank you mail due to %s", donationType, body.OrderNumber, err.Error())
	}

}

// Handler for an authenticated user to create a periodic donation
func (mc *MembershipController) CreateAPeriodicDonationOfAUser(c *gin.Context) (int, gin.H, error) {
	const errWhere = "MembershipController.CreateAPeriodicDonationOfAUser"

	// Validate client request
	var reqBody clientReq

	if failData, valid := bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	userID := reqBody.UserID

	// Build a draft periodic donation record
	draftPeriodicDonation := buildDraftPeriodicDonation(userID, reqBody)

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
	resp.PeriodicID = draftPeriodicDonation.ID
	resp.ID = draftRecord.ID

	// send success mail asynchronously
	go mc.sendDonationThankYouMail(resp, "定期定額")

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

	userID := reqBody.UserID

	// Start Tappay transaction

	// Build Tappay pay by prime request
	tapPayReq := buildTapPayPrimeReq(oneTime, payMethod, reqBody)

	draftRecord := buildPrimeDraftRecord(userID, payMethod, tapPayReq)

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

	// send success mail asynchronously
	go mc.sendDonationThankYouMail(resp, "單筆捐款")

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
	var rowsAffected int64
	var valid bool

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

	if err, rowsAffected = mc.Storage.UpdateByConditions(map[string]interface{}{
		"user_id": reqBody.UserID,
		"id":      recordID,
	}, d); err != nil {
		return 0, gin.H{}, err
	}

	if rowsAffected == 0 {
		return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{
			"uri": fmt.Sprintf("%s can not address a resource", c.Request.RequestURI)},
		}, nil
	}

	return http.StatusNoContent, gin.H{}, nil
}

// TODO
// GetDonationsOfAUser returns the donations list of a user
func (mc *MembershipController) GetDonationsOfAUser(c *gin.Context) (int, gin.H, error) {
	return 0, gin.H{}, nil
}

// GetADonationOfAUser returns a donation of a user
func (mc *MembershipController) GetADonationOfAUser(c *gin.Context, donationType string) (int, gin.H, error) {
	var err error
	var recordID uint64
	var userID uint64
	var _userID uint
	var resp = new(clientResp)

	if userID, err = strconv.ParseUint(c.Query("user_id"), 10, strconv.IntSize); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{
			"req.URL.query": "?user_id=:userID, userID should be integer",
		}}, err
	}

	if recordID, err = strconv.ParseUint(c.Param("id"), 10, strconv.IntSize); err != nil {
		return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{
			"url": fmt.Sprintf("%s cannot address a found resource", c.Request.RequestURI),
		}}, err
	}

	switch donationType {
	case globals.PeriodicDonationType:
		pd := models.PeriodicDonation{}
		if err = mc.Storage.Get(uint(recordID), &pd); err != nil {
			return 0, gin.H{}, err
		}
		resp.BuildFromPeriodicDonationModel(pd)
		_userID = uint(pd.UserID)
		break
	case globals.PrimeDonaitionType:
		pd := models.PayByPrimeDonation{}
		if err = mc.Storage.Get(uint(recordID), &pd); err != nil {
			return 0, gin.H{}, err
		}
		resp.BuildFromPrimeDonationModel(pd)
		_userID = uint(pd.UserID)
		break
	case globals.OthersDonationType:
		d := models.PayByOtherMethodDonation{}
		if err = mc.Storage.Get(uint(recordID), &d); err != nil {
			return 0, gin.H{}, err
		}
		resp.BuildFromOtherMethodDonationModel(d)
		_userID = uint(d.UserID)
		break

	// TODO
	// case globasl.PayByCardTokenDonation:

	default:
		return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("donation type %s is not supported", donationType)}, nil
	}

	if err != nil {
		appErr, _ := err.(*models.AppError)
		if appErr.StatusCode == http.StatusNotFound {
			return appErr.StatusCode, gin.H{"status": "fail", "data": gin.H{
				"url": fmt.Sprintf("%s cannot address a found resource", c.Request.RequestURI),
			}}, nil
		}
		return 0, gin.H{}, err
	}

	if _userID != uint(userID) {
		return http.StatusForbidden, gin.H{"status": "fail", "data": gin.H{
			"authorization": fmt.Sprintf("%s is forbidden to access", c.Request.RequestURI),
		}}, nil
	}

	return http.StatusOK, gin.H{"status": "success", "data": resp}, nil
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

	copier.Copy(&m, &req.Cardholder)

	m.Amount = req.Amount
	m.UserID = userID
	m.Status = statusPeriodicPaying

	return m
}

func buildPrimeDraftRecord(userID uint, payMethod string, req tapPayPrimeReq) models.PayByPrimeDonation {
	m := models.PayByPrimeDonation{}

	m.UserID = userID
	m.PayMethod = payMethod

	copier.Copy(&m, &req)
	copier.Copy(&m, &req.Cardholder)

	m.Status = statusPaying

	return m
}

func buildPrimeSuccessRecord(resp tapPayTransactionResp) models.PayByPrimeDonation {
	m := models.PayByPrimeDonation{}

	copier.Copy(&m, &resp)
	copier.Copy(&m, &resp.CardInfo)

	m.TappayApiStatus = null.IntFrom(resp.Status)

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

	m.Status = "paid"

	return m
}

func buildSuccessPeriodicDonation(resp tapPayTransactionResp) models.PeriodicDonation {
	m := models.PeriodicDonation{}

	copier.Copy(&m, &resp.CardInfo)

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

	copier.Copy(&m, &resp)
	m.TappayApiStatus = null.IntFrom(resp.Status)

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

	copier.Copy(&m, &req)

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
