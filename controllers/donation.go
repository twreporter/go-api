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
	"gopkg.in/go-playground/validator.v8"
	"gopkg.in/guregu/null.v3"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

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

	tapPayRespStatusSuccess = 0

	defaultPeriodicPayMethod = "credit_card"

	secToMsec     = 1000
	msecToNanosec = 1000000

	monthlyFrequency = "monthly"
	yearlyFrequency  = "yearly"
	oneTimeFrequency = "one_time"
)

// pay type Enum
const (
	prime payType = iota
	token
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

type (
	clientReq struct {
		Amount     uint              `json:"amount" form:"amount" binding:"required"`
		Cardholder models.Cardholder `json:"donor" form:"donor" binding:"required,dive"`
		Currency   string            `json:"currency" form:"currency"`
		Details    string            `json:"details" form:"details"`
		Frequency  string            `json:"frequency"`
		MerchantID string            `json:"merchant_id" form:"merchant_id"`
		PayMethod  string            `json:"pay_method" form:"pay_method"`
		Prime      string            `json:"prime" form:"prime" binding:"required"`
		ResultUrl  linePayResultUrl  `json:"result_url" form:"result_url"`
		UserID     uint              `json:"user_id" form:"user_id" binding:"required"`
	}

	clientResp struct {
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

	tapPayTransactionReq struct {
		Amount      uint              `json:"amount"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		MerchantID  string            `json:"merchant_id"`
		OrderNumber string            `json:"order_number"`
		PartnerKey  string            `json:"partner_key"`
		Prime       string            `json:"prime"`
		Remember    bool              `json:"remember"`
		ResultUrl   linePayResultUrl  `json:"result_url"`
	}

	tapPayTransactionResp struct {
		models.TappayResp
		BankTransactionTime   bankTransactionTime `json:"bank_transaction_time"`
		CardInfo              models.CardInfo     `json:"card_info"`
		CardSecret            cardSecret          `json:"card_secret"`
		Status                int64               `json:"status"`
		TransactionTimeMillis int64               `json:"transaction_time_millis"`
	}

	tapPayMinTransactionResp struct {
		Status int64  `json:"status"`
		Msg    string `json:"msg"`
	}

	payType int

	patchBody struct {
		Donor       models.Cardholder `json:"donor"`
		Notes       string            `json:"notes"`
		SendReceipt string            `json:"send_receipt"`
		ToFeedback  bool              `json:"to_feedback"`
		UserID      uint              `json:"user_id" binding:"required"`
	}
)

func (p *patchBody) BuildPeriodicDonation() models.PeriodicDonation {
	m := new(models.PeriodicDonation)
	m.Cardholder = p.Donor
	m.Notes = p.Notes
	m.SendReceipt = p.SendReceipt
	m.ToFeedback = null.BoolFrom(p.ToFeedback)
	m.UserID = p.UserID
	return *m
}

func (p *patchBody) BuildPrimeDonation() models.PayByPrimeDonation {
	m := new(models.PayByPrimeDonation)
	m.Cardholder = p.Donor
	m.Notes = p.Notes
	m.SendReceipt = p.SendReceipt
	m.UserID = p.UserID
	return *m
}

func (req clientReq) BuildTapPayReq(orderNumber string) tapPayTransactionReq {
	primeReq := new(tapPayTransactionReq)
	primeReq.Prime = req.Prime
	primeReq.OrderNumber = orderNumber
	primeReq.Amount = req.Amount

	if req.Currency != "" {
		primeReq.Currency = req.Currency
	} else {
		primeReq.Currency = defaultCurrency
	}

	if req.Details != "" {
		primeReq.Details = req.Details
	} else {
		primeReq.Details = defaultDetails
	}

	if req.MerchantID != "" {
		primeReq.MerchantID = req.MerchantID
	} else {
		primeReq.MerchantID = defaultMerchantID
	}

	primeReq.Cardholder = req.Cardholder
	// Per required fields (even empty) of cardholder of tappay documents,
	// use empty strings for name and phonenumber fields instead of empty.
	if !primeReq.Cardholder.Name.Valid {
		primeReq.Cardholder.Name = null.StringFrom("")
	}

	if !primeReq.Cardholder.PhoneNumber.Valid {
		primeReq.Cardholder.PhoneNumber = null.StringFrom("")
	}

	primeReq.PartnerKey = globals.Conf.Donation.TapPayPartnerKey

	if req.Frequency == monthlyFrequency || req.Frequency == yearlyFrequency {
		primeReq.Remember = true
	}

	primeReq.ResultUrl = req.ResultUrl
	return *primeReq
}

func (req clientReq) BuildDraftPeriodicDonation(orderNumber string) models.PeriodicDonation {
	const defaultDetails = "一般線上定期定額捐款"

	m := new(models.PeriodicDonation)

	m.Amount = req.Amount
	m.Cardholder = req.Cardholder
	m.Currency = req.Currency
	m.UserID = req.UserID

	if req.Frequency != "" {
		m.Frequency = req.Frequency
	} else {
		m.Frequency = monthlyFrequency
	}

	if req.Details != "" {
		m.Details = req.Details
	} else {
		m.Details = defaultDetails
	}

	m.OrderNumber = orderNumber
	m.Status = statusPaying

	return *m
}

func (req clientReq) BuildPrimeDraftRecord(orderNumber string, payMethod string) models.PayByPrimeDonation {
	m := new(models.PayByPrimeDonation)

	m.Amount = req.Amount
	m.Cardholder = req.Cardholder
	m.Currency = req.Currency
	m.Details = req.Details
	m.MerchantID = req.MerchantID
	m.UserID = req.UserID
	m.PayMethod = payMethod
	m.OrderNumber = orderNumber
	m.Status = statusPaying

	return *m
}

func (req clientReq) BuildTokenDraftRecord(orderNumber string) models.PayByCardTokenDonation {
	m := new(models.PayByCardTokenDonation)

	m.Amount = req.Amount
	m.Currency = req.Currency
	m.Details = req.Details
	m.MerchantID = req.MerchantID
	m.OrderNumber = orderNumber
	m.Status = statusPaying

	return *m
}

func (cr *clientResp) BuildFromPeriodicDonationModel(d models.PeriodicDonation) {
	cr.Amount = d.Amount
	cr.Cardholder = d.Cardholder
	cr.CardInfo = d.CardInfo
	cr.Currency = d.Currency
	cr.Details = d.Details
	cr.Frequency = d.Frequency
	cr.ID = d.ID
	cr.Notes = d.Notes
	cr.OrderNumber = d.OrderNumber
	cr.SendReceipt = d.SendReceipt
	cr.ToFeedback = d.ToFeedback.ValueOrZero()
}

func (cr *clientResp) BuildFromPrimeDonationModel(d models.PayByPrimeDonation) {
	cr.Amount = d.Amount
	cr.Cardholder = d.Cardholder
	cr.CardInfo = d.CardInfo
	cr.Currency = d.Currency
	cr.Details = d.Details
	cr.ID = d.ID
	cr.Notes = d.Notes
	cr.OrderNumber = d.OrderNumber
	cr.PayMethod = d.PayMethod
	cr.SendReceipt = d.SendReceipt
	cr.ToFeedback = false
	cr.Frequency = oneTimeFrequency
}

func (cr *clientResp) BuildFromOtherMethodDonationModel(d models.PayByOtherMethodDonation) {
	cr.Amount = d.Amount
	cr.Cardholder = models.Cardholder{
		Name:        null.StringFrom(d.Name),
		Email:       d.Email,
		NationalID:  null.StringFrom(d.NationalID),
		Address:     null.StringFrom(d.Address),
		PhoneNumber: null.StringFrom(d.PhoneNumber),
		ZipCode:     null.StringFrom(d.ZipCode),
	}
	cr.Currency = d.Currency
	cr.Details = d.Details
	cr.ID = d.ID
	cr.Notes = d.Notes
	cr.OrderNumber = d.OrderNumber
	cr.PayMethod = d.PayMethod
	cr.SendReceipt = d.SendReceipt
	cr.ToFeedback = false
	cr.Frequency = oneTimeFrequency
}

func bindRequestBody(c *gin.Context, reqBody interface{}) (gin.H, bool) {
	var err error
	// Validate request body
	// gin.Context.Bind does not support to bind `JSON` body multiple times
	// the alternative is to use gin.Context.ShouldBindBodyWith function to bind
	if err = c.ShouldBindBodyWith(reqBody, binding.JSON); nil == err {
		// omit intentionally
	} else if err = c.Bind(reqBody); nil != err {
		// bind other format rather than JSON
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
	const donationDetails = "第一筆定期定額捐款"

	// Validate client request
	var err error
	var reqBody clientReq
	var tapPayResp tapPayTransactionResp

	if failData, valid := bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	if reqBody.Cardholder.Email == "" {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{
			"req.Body.donor.email": "donor email is not valid",
		}}, nil
	}

	if reqBody.Frequency != monthlyFrequency && reqBody.Frequency != yearlyFrequency {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{
			"req.Body.frequency": "frequency is not supported. should be `monthly` or `yearly`",
		}}, nil
	}

	// generate periodic donation order number
	pdOrderNumber := generateOrderNumber(periodic, getPayMethodID(payMethodCollections[0]))
	// Build a draft periodic donation record
	periodicDonation := reqBody.BuildDraftPeriodicDonation(pdOrderNumber)

	// generate token donation order number
	dOrderNumber := generateOrderNumber(token, getPayMethodID(payMethodCollections[0]))
	// Build a draft card token donation record
	tokenDonation := reqBody.BuildTokenDraftRecord(dOrderNumber)

	// Build Tappay prime request
	tapPayReq := reqBody.BuildTapPayReq(dOrderNumber)
	tapPayReq.Details = fmt.Sprintf("%s;%s", donationDetails, tapPayReq.Details)

	// Create a draft periodic donation along with the first token donation record of that periodic donation
	err = mc.Storage.CreateAPeriodicDonation(&periodicDonation, &tokenDonation)
	if nil != err {
		errMsg := "Unable to create a draft periodic donation and the first card token transaction record"
		log.Error(fmt.Sprintf("%s: %s", errWhere, errMsg))
		return 0, gin.H{}, models.NewAppError(errWhere, errMsg, err.Error(), http.StatusInternalServerError)
	}

	// Start Tappay transaction
	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err = serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'fail' and stop the periodic donation
			failResp := models.PayByCardTokenDonation{}
			failResp.TappayApiStatus = null.IntFrom(tapPayResp.Status)
			failResp.Msg = tapPayResp.Msg
			failResp.Status = statusFail
			// Procceed even if the deletion is failed
			mc.Storage.DeleteAPeriodicDonation(periodicDonation.ID, failResp)
		}
		errMsg := err.Error()
		log.Error(fmt.Sprintf("%s: %s", errWhere, errMsg))

		return 0, gin.H{}, models.NewAppError(errWhere, errMsg, "", http.StatusInternalServerError)
	}

	// append tappay response onto donation model
	tapPayResp.AppendRespOnPerodicDonation(&periodicDonation)
	tapPayResp.AppendRespOnTokenDonation(&tokenDonation)

	if err = mc.Storage.UpdatePeriodicAndCardTokenDonationInTRX(periodicDonation.ID, periodicDonation, tokenDonation); nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
	}

	// build response for clients
	resp := new(clientResp)
	resp.BuildFromPeriodicDonationModel(periodicDonation)

	// send success mail asynchronously
	go mc.sendDonationThankYouMail(*resp, "定期定額")

	return http.StatusCreated, gin.H{"status": "success", "data": resp}, nil
}

// Handler for an authenticated user to create an one-time donation
func (mc *MembershipController) CreateADonationOfAUser(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.CreateADonationOfAUser"
	var err error
	var reqBody clientReq
	var tapPayResp tapPayTransactionResp

	// Validate client request
	if failData, valid := bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	// Validate pay_method
	payMethod := reqBody.PayMethod
	if err = validatePayMethod(reqBody.PayMethod); nil != err {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{"req.Body.pay_method": err.Error()}}, nil
	}

	if reqBody.Cardholder.Email == "" {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": gin.H{
			"req.Body.donor.email": "donor email is not valid",
		}}, nil
	}

	// generate token donation order number
	dOrderNumber := generateOrderNumber(prime, getPayMethodID(payMethod))
	// Build a draft card prime donation record
	primeDonation := reqBody.BuildPrimeDraftRecord(dOrderNumber, payMethod)

	// Start Tappay transaction
	// Build Tappay pay by prime request
	tapPayReq := reqBody.BuildTapPayReq(dOrderNumber)

	if err = mc.Storage.Create(&primeDonation); nil != err {
		switch appErr := err.(type) {
		case *models.AppError:
			return 0, gin.H{}, models.NewAppError(errorWhere, "Fails to create a draft prime record", appErr.Error(), appErr.StatusCode)
		default:
			return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Unknown Error Type. Fails to create a draft prime record. %v", err.Error())}, nil
		}
	}

	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err = serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'fail'
			d := models.PayByPrimeDonation{}
			d.TappayApiStatus = null.IntFrom(tapPayResp.Status)
			d.Msg = tapPayResp.Msg
			d.Status = statusFail
			mc.Storage.UpdateByConditions(map[string]interface{}{
				"id": primeDonation.ID,
			}, d)
		}
		return 0, gin.H{}, models.NewAppError(errorWhere, err.Error(), "", http.StatusInternalServerError)
	}

	// append tappay response onto donation model
	tapPayResp.AppendRespOnPrimeDonation(&primeDonation)

	if err, _ = mc.Storage.UpdateByConditions(map[string]interface{}{
		"id": primeDonation.ID,
	}, primeDonation); nil != err {
		log.Error(err.Error())
	}

	// build response for clients
	resp := new(clientResp)
	resp.BuildFromPrimeDonationModel(primeDonation)

	// send success mail asynchronously
	go mc.sendDonationThankYouMail(*resp, "單筆捐款")

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
		return http.StatusNotFound, gin.H{"status": "error", "message": "record not found, record id should be provided in the url"}, nil
	}

	if failData, valid = bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	switch donationType {
	case globals.PeriodicDonationType:
		d = reqBody.BuildPeriodicDonation()
	case globals.PrimeDonaitionType:
		d = reqBody.BuildPrimeDonation()
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
		}}, nil
	}

	if recordID, err = strconv.ParseUint(c.Param("id"), 10, strconv.IntSize); err != nil {
		return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{
			"url": fmt.Sprintf("%s cannot address a found resource", c.Request.RequestURI),
		}}, nil
	}

	switch donationType {
	case globals.PeriodicDonationType:
		d := models.PeriodicDonation{}
		err = mc.Storage.Get(uint(recordID), &d)
		resp.BuildFromPeriodicDonationModel(d)
		_userID = uint(d.UserID)
		break
	case globals.PrimeDonaitionType:
		d := models.PayByPrimeDonation{}
		err = mc.Storage.Get(uint(recordID), &d)
		resp.BuildFromPrimeDonationModel(d)
		_userID = uint(d.UserID)
		break
	case globals.OthersDonationType:
		d := models.PayByOtherMethodDonation{}
		err = mc.Storage.Get(uint(recordID), &d)
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
				"req.URL": fmt.Sprintf("%s cannot address a found resource", c.Request.RequestURI),
			}}, nil
		}
		return 0, gin.H{}, err
	}

	if _userID != uint(userID) {
		return http.StatusForbidden, gin.H{"status": "fail", "data": gin.H{
			"req.Headers.Authorization": fmt.Sprintf("%s is forbidden to access", c.Request.RequestURI),
		}}, nil
	}

	return http.StatusOK, gin.H{"status": "success", "data": resp}, nil
}

func (resp tapPayTransactionResp) AppendRespOnPrimeDonation(m *models.PayByPrimeDonation) {
	m.CardInfo = resp.CardInfo
	m.TappayResp = resp.TappayResp
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
}

func (resp tapPayTransactionResp) AppendRespOnPerodicDonation(m *models.PeriodicDonation) {
	m.CardInfo = resp.CardInfo

	ciphertext := encrypt(resp.CardSecret.CardToken, globals.Conf.Donation.CardSecretKey)
	m.CardToken = ciphertext

	ciphertext = encrypt(resp.CardSecret.CardKey, globals.Conf.Donation.CardSecretKey)
	m.CardKey = ciphertext

	t := time.Now()
	m.LastSuccessAt = null.TimeFrom(t)
	m.Status = statusPaid
}

func (resp tapPayTransactionResp) AppendRespOnTokenDonation(m *models.PayByCardTokenDonation) {
	m.TappayResp = resp.TappayResp
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
