package controllers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"
	"gopkg.in/go-playground/validator.v8"
	"gopkg.in/guregu/null.v3"

	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/storage"
)

const (
	defaultCurrency           = "TWD"
	defaultCreditCardMerchant = "GlobalTesting_CTBC"
	defaultLineMerchant       = "GlobalTesting_LINEPAY" // TODO: Need to revise after the application is done

	invalidPayMethodID = -1

	orderPrefix = "twreporter"

	statusPaying  = "paying"
	statusPaid    = "paid"
	statusFail    = "fail"
	statusStopped = "stopped"
	statusInvalid = "invalid"

	tapPayRespStatusSuccess = 0

	defaultPeriodicPayMethod = "credit_card"

	secToMsec     = 1000
	msecToNanosec = 1000000

	monthlyFrequency = "monthly"
	yearlyFrequency  = "yearly"
	oneTimeFrequency = "one_time"

	payMethodCreditCard = "credit_card"
	payMethodLine       = "line"
	payMethodGoogle     = "google"
	payMethodApple      = "apple"
	payMethodSamsung    = "samsung"

	linePayMethodCreditCard = "CREDIT_CARD"
	linePayMethodBalance    = "BALANCE"
	linePayMethodPoint      = "POINT"
)

// pay type Enum
const (
	prime payType = iota
	token
	periodic
)

// pay method collections
var payMethodCollections = []string{
	payMethodCreditCard,
	payMethodLine,
	payMethodGoogle,
	payMethodApple,
	payMethodSamsung,
}

var payMethodMap = map[string]string{
	payMethodCreditCard: "信用卡支付",
	payMethodLine:       "Line Pay",
	payMethodGoogle:     "Google Pay",
	payMethodApple:      "Apple Pay",
	payMethodSamsung:    "Samsung Pay",
}

var methodToMerchant = map[string]string{
	payMethodCreditCard: defaultCreditCardMerchant,
	payMethodLine:       defaultLineMerchant,
}

var cardInfoTypes = map[int64]string{
	1: "VISA",
	2: "MasterCard",
	3: "JCB",
	4: "Union Pay",
	5: "AMEX",
}

var linePayMethods = []string{
	linePayMethodCreditCard,
	linePayMethodBalance,
	linePayMethodPoint,
}

type (
	clientReq struct {
		Amount       uint              `json:"amount" binding:"required"`
		Cardholder   models.Cardholder `json:"donor" binding:"required,dive"`
		Currency     string            `json:"currency"`
		Details      string            `json:"details"`
		Frequency    string            `json:"frequency"`
		MerchantID   string            `json:"merchant_id"`
		PayMethod    string            `json:"pay_method"`
		Prime        string            `json:"prime" binding:"required"`
		UserID       uint              `json:"user_id" binding:"required"`
		MaxPaidTimes uint              `json:"max_paid_times"`
	}

	clientResp struct {
		Amount        uint              `json:"amount"`
		CardInfo      models.CardInfo   `json:"card_info"`
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

	bankTransactionTime struct {
		StartTimeMillis string `json:"start_time_millis"`
		EndTimeMillis   string `json:"end_time_millis"`
	}

	cardSecret struct {
		CardToken string `json:"card_token"`
		CardKey   string `json:"card_key"`
	}

	linePayResultUrl struct {
		FrontendRedirectUrl string `json:"frontend_redirect_url"`
		BackendNotifyUrl    string `json:"backend_notify_url"`
	}

	tapPayTransactionReq struct {
		Amount                 uint              `json:"amount"`
		Cardholder             models.Cardholder `json:"cardholder"`
		Currency               string            `json:"currency"`
		Details                string            `json:"details"`
		MerchantID             string            `json:"merchant_id"`
		OrderNumber            string            `json:"order_number"`
		PartnerKey             string            `json:"partner_key"`
		Prime                  string            `json:"prime"`
		Remember               bool              `json:"remember"`
		ResultUrl              linePayResultUrl  `json:"result_url"`
		LinePayProductImageUrl null.String       `json:"line_pay_product_image_url"`
	}

	tapPayTransactionResp struct {
		models.TappayResp
		models.PayInfo        `json:"pay_info"`
		BankTransactionTime   bankTransactionTime `json:"bank_transaction_time"`
		CardInfo              models.CardInfo     `json:"card_info"`
		CardSecret            cardSecret          `json:"card_secret"`
		Status                int64               `json:"status"`
		TransactionTimeMillis int64               `json:"transaction_time_millis"`
		PaymentUrl            string              `json:"payment_url"`
		Amount                int                 `json:"amount"`
		OrderNumber           string              `json:"order_number"`
	}

	tapPayMinTransactionResp struct {
		Status int64  `json:"status"`
		Msg    string `json:"msg"`
	}

	tradeRecord struct {
		RecordStatus int `json:"record_status"`
	}

	tapPayTransactionRecordResp struct {
		Status       int           `json:"status"`
		Msg          string        `json:"msg"`
		TradeRecords []tradeRecord `json:"trade_records"`
	}

	tapPayTransactionRecordReq struct {
		PartnerKey     string      `json:"partner_key"`
		RecordsPerPage uint        `json:"records_per_page"`
		Filters        queryFilter `json:"filters"`
	}

	payType int

	patchBody struct {
		Donor         models.Cardholder `json:"donor"`
		Notes         string            `json:"notes"`
		SendReceipt   string            `json:"send_receipt"`
		ToFeedback    bool              `json:"to_feedback"`
		UserID        uint              `json:"user_id" binding:"required"`
		IsAnonymous   bool              `json:"is_anonymous"`
		ReceiptHeader null.String       `json:"receipt_header"`
	}

	queryFilterTime struct {
		StartTime null.Int `json:"start_time"`
		EndTime   null.Int `json:"end_time"`
	}

	queryFilter struct {
		OrderNumber       string           `json:"order_number" binding:"required"`
		BankTransactionID null.String      `json:"bank_transaction_id"`
		RecTradeID        null.String      `json:"rec_trade_id"`
		Time              *queryFilterTime `json:"time"`
	}

	queryReq struct {
		RecordsPerPage uint        `json:"records_per_page"`
		Filters        queryFilter `json:"filters" binding:"required,dive"`
	}
)

func (p *patchBody) BuildPeriodicDonation() models.PeriodicDonation {
	m := new(models.PeriodicDonation)
	m.Cardholder = p.Donor
	m.Notes = p.Notes
	m.SendReceipt = p.SendReceipt
	m.ToFeedback = null.BoolFrom(p.ToFeedback)
	m.UserID = p.UserID
	m.IsAnonymous = null.BoolFrom(p.IsAnonymous)
	m.ReceiptHeader = p.ReceiptHeader
	return *m
}

func (p *patchBody) BuildPrimeDonation() models.PayByPrimeDonation {
	m := new(models.PayByPrimeDonation)
	m.Cardholder = p.Donor
	m.Notes = p.Notes
	m.SendReceipt = p.SendReceipt
	m.UserID = p.UserID
	m.IsAnonymous = null.BoolFrom(p.IsAnonymous)
	m.ReceiptHeader = p.ReceiptHeader
	return *m
}

func (req clientReq) BuildTapPayReq(orderNumber, details, payMethod string) tapPayTransactionReq {
	primeReq := new(tapPayTransactionReq)
	primeReq.Prime = req.Prime
	primeReq.OrderNumber = orderNumber
	primeReq.Amount = req.Amount

	if req.Currency != "" {
		primeReq.Currency = req.Currency
	} else {
		primeReq.Currency = defaultCurrency
	}

	primeReq.Details = details

	if req.MerchantID != "" {
		primeReq.MerchantID = req.MerchantID
	} else {
		primeReq.MerchantID = methodToMerchant[payMethod]
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

	f := ""
	if req.Frequency == monthlyFrequency || req.Frequency == yearlyFrequency {
		primeReq.Remember = true
		f = req.Frequency
	} else {
		f = "one_time"
	}

	// Only build resultUrl and linePayProductImageUrl during linepay transaction
	if payMethod == payMethodLine {
		frontendRedirectUrl := "https://" + globals.Conf.Donation.FrontendHost + "/contribute/line/" + f + "/" + orderNumber

		// Tappay server will validate the hosts provided in the result_url
		// Wrap the backendHost to be test.twreporter.org if not in the staging or production environment
		backendHost := ""
		if globals.Conf.Environment == "production" || globals.Conf.Environment == "staging" {
			backendHost = globals.Conf.App.Host
		} else {
			backendHost = "test.twreporter.org"
		}

		primeReq.ResultUrl = linePayResultUrl{
			FrontendRedirectUrl: frontendRedirectUrl,
			BackendNotifyUrl:    "https://" + backendHost + "/v1/donations/prime/line-notify",
		}

		primeReq.LinePayProductImageUrl = null.StringFrom(globals.Conf.Donation.LinePayProductImageUrl)
	}

	return *primeReq
}

func (req clientReq) BuildDraftPeriodicDonation(orderNumber string) models.PeriodicDonation {
	const defaultDetails = "一般線上定期定額捐款"
	const defaultMaxPaidTimes = 2147483647

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

	// If MaxPaidTimes is not specified or zero value, set it to default maximum paid times.
	if req.MaxPaidTimes != 0 {
		m.MaxPaidTimes = req.MaxPaidTimes
	} else {
		m.MaxPaidTimes = defaultMaxPaidTimes
	}

	return *m
}

func (req clientReq) BuildPrimeDraftRecord(orderNumber string, payMethod string) models.PayByPrimeDonation {
	const defaultDetails = "一般線上單筆捐款"
	m := new(models.PayByPrimeDonation)

	m.Amount = req.Amount
	m.Cardholder = req.Cardholder
	m.Currency = req.Currency

	if req.Details != "" {
		m.Details = req.Details
	} else {
		m.Details = defaultDetails
	}

	m.MerchantID = req.MerchantID
	m.UserID = req.UserID
	m.PayMethod = payMethod
	m.OrderNumber = orderNumber
	m.Status = statusPaying

	return *m
}

func (req clientReq) BuildTokenDraftRecord(orderNumber string) models.PayByCardTokenDonation {
	const defaultDetails = "一般線上定期定額捐款"
	m := new(models.PayByCardTokenDonation)

	m.Amount = req.Amount
	m.Currency = req.Currency

	if req.Details != "" {
		m.Details = req.Details
	} else {
		m.Details = defaultDetails
	}

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
	cr.PayMethod = payMethodCreditCard
	cr.IsAnonymous = d.IsAnonymous.ValueOrZero()
	cr.ReceiptHeader = d.ReceiptHeader
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
	cr.IsAnonymous = d.IsAnonymous.ValueOrZero()
	cr.ReceiptHeader = d.ReceiptHeader
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

func bindRequestJSONBody(c *gin.Context, reqBody interface{}) (gin.H, bool) {
	var err error
	// Validate request body
	// gin.Context.Bind does not support to bind `JSON` body multiple times
	// the alternative is to use gin.Context.ShouldBindBodyWith function to bind
	if err = c.ShouldBindBodyWith(reqBody, binding.JSON); nil != err {
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

func (mc *MembershipController) sendDonationThankYouMail(body clientResp) {
	var origin string
	switch globals.Conf.Environment {
	case globals.DevelopmentEnvironment:
		origin = globals.SupportSiteDevOrigin
	case globals.StagingEnvironment:
		origin = globals.SupportSiteStagingOrigin
	case globals.ProductionEnvironment:
		origin = globals.SupportSiteOrigin
	default:
		origin = globals.SupportSiteOrigin
	}

	var donationLink string = origin + "/contribute/" + body.Frequency + "/" + body.OrderNumber + "?utm_source=supportsuccess&utm_medium=email"

	var donationType string
	switch body.Frequency {
	case oneTimeFrequency:
		donationType = "單筆捐款"
	case monthlyFrequency:
		donationType = "定期定額"
	case yearlyFrequency:
		donationType = "定期定額"
	default:
		donationType = "捐款"
	}

	reqBody := donationSuccessReqBody{
		Amount:         body.Amount,
		Currency:       body.Currency,
		DonationMethod: payMethodMap[body.PayMethod],
		DonationType:   donationType,
		DonationLink:   donationLink,
		Email:          body.Cardholder.Email,
		Name:           body.Cardholder.Name.ValueOrZero(),
		OrderNumber:    body.OrderNumber,
	}

	if err := postMailServiceEndpoint(reqBody, fmt.Sprintf("http://localhost:%s/v1/%s", globals.LocalhostPort, globals.SendSuccessDonationRoutePath)); err != nil {

		err = errors.Wrap(err, fmt.Sprintf("fail to send %s donation(order_number: %s) mail", donationType, body.OrderNumber))

		if globals.Conf.Environment == "development" {
			log.Errorf("%+v", err)
		} else {
			log.WithField("detail", err).Errorf("%s", f.FormatStack(err))
		}
	}

}

// Handler for an authenticated user to create a periodic donation
func (mc *MembershipController) CreateAPeriodicDonationOfAUser(c *gin.Context) (int, gin.H, error) {
	// Validate client request
	var err error
	var reqBody clientReq
	var tapPayResp tapPayTransactionResp

	if failData, valid := bindRequestJSONBody(c, &reqBody); valid == false {
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
	tapPayReq := reqBody.BuildTapPayReq(dOrderNumber, tokenDonation.Details, payMethodCreditCard)

	// Create a draft periodic donation along with the first token donation record of that periodic donation
	err = mc.Storage.CreateAPeriodicDonation(&periodicDonation, &tokenDonation)
	if nil != err {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to create a draft periodic donation and the first card token transaction record"}, err
	}

	// Start Tappay transaction
	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err = serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'fail' and mark the periodic donation as 'invalid'.
			td := models.PayByCardTokenDonation{}
			tapPayResp.AppendRespOnTokenDonation(&td, statusFail)

			pd := models.PeriodicDonation{}
			pd.Status = statusInvalid
			pd.CardInfo = tapPayResp.CardInfo

			mc.Storage.UpdatePeriodicAndCardTokenDonationInTRX(periodicDonation.ID, pd, td)
		}
		return http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()}, err
	}

	// append tappay response onto donation model
	tapPayResp.AppendRespOnPerodicDonation(&periodicDonation)
	tapPayResp.AppendRespOnTokenDonation(&tokenDonation, statusPaid)

	// since the donation already succeeded, return transaction success even if the information patch fails
	if err = mc.Storage.UpdatePeriodicAndCardTokenDonationInTRX(periodicDonation.ID, periodicDonation, tokenDonation); nil != err {
		log.Infof("%v", err)
	}

	// build response for clients
	resp := new(clientResp)
	resp.BuildFromPeriodicDonationModel(periodicDonation)

	// send success mail asynchronously
	go mc.sendDonationThankYouMail(*resp)

	return http.StatusCreated, gin.H{"status": "success", "data": resp}, nil
}

// Handler for an authenticated user to create an one-time donation
func (mc *MembershipController) CreateADonationOfAUser(c *gin.Context) (int, gin.H, error) {
	var err error
	var reqBody clientReq
	var tapPayResp tapPayTransactionResp

	// Validate client request
	if failData, valid := bindRequestJSONBody(c, &reqBody); valid == false {
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
	tapPayReq := reqBody.BuildTapPayReq(dOrderNumber, primeDonation.Details, payMethod)

	if err = mc.Storage.Create(&primeDonation); nil != err {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "Fails to create a draft prime record"}, err
	}

	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err = serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'fail'
			d := models.PayByPrimeDonation{}
			tapPayResp.AppendRespOnPrimeDonation(&d, statusFail)

			mc.Storage.UpdateByConditions(map[string]interface{}{
				"id": primeDonation.ID,
			}, d)
		}
		return http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()}, err
	}

	// Append tappay response onto donation model
	// Since linepay requires extra transaction process,
	// wait for the line-notify endpoint to update the final transaction status
	if primeDonation.PayMethod == payMethodLine {
		tapPayResp.AppendRespOnPrimeDonation(&primeDonation, statusPaying)
	} else {
		tapPayResp.AppendRespOnPrimeDonation(&primeDonation, statusPaid)
	}

	// since the donation already succeeded, return transaction success even if the information patch fails
	if err, _ = mc.Storage.UpdateByConditions(map[string]interface{}{
		"id": primeDonation.ID,
	}, primeDonation); nil != err {
		log.Infof("%v", err)
	}

	// build response for clients
	resp := new(clientResp)
	resp.BuildFromPrimeDonationModel(primeDonation)
	resp.PaymentUrl = tapPayResp.PaymentUrl

	// only send mail if the transaction completed.
	// send success mail asynchronously
	if primeDonation.Status == statusPaid {
		go mc.sendDonationThankYouMail(*resp)
	}

	return http.StatusCreated, gin.H{"status": "success", "data": resp}, nil
}

// PatchADonationOfAUser method
// Handler for an authenticated user to patch an prime/token/periodic donation
func (mc *MembershipController) PatchADonationOfAUser(c *gin.Context, donationType string) (int, gin.H, error) {
	var d interface{}
	var err error
	var failData gin.H
	var reqBody patchBody
	var rowsAffected int64
	var valid bool
	var orderNumber string

	if orderNumber = c.Param("order"); "" == orderNumber {
		return http.StatusNotFound, gin.H{"status": "error", "message": "record not found, order_number should be provided in the url"}, nil
	}

	if failData, valid = bindRequestJSONBody(c, &reqBody); valid == false {
		log.WithField("payload", reqBody).Infof("cannot patch the personal info of the donor, %v", failData)
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	switch donationType {
	case globals.PeriodicDonationType:
		d = reqBody.BuildPeriodicDonation()
	case globals.PrimeDonationType:
		d = reqBody.BuildPrimeDonation()
	default:
		return http.StatusInternalServerError,
			gin.H{"status": "error", "message": fmt.Sprintf("donation type(%s) not supported", donationType)},
			nil
	}

	if err, rowsAffected = mc.Storage.UpdateByConditions(map[string]interface{}{
		"user_id":      reqBody.UserID,
		"order_number": orderNumber,
	}, d); err != nil {
		return http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "unable to patch the record",
		}, err
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
	var (
		err         error
		_userID     uint
		resp        = new(clientResp)
		authUserID  interface{}
		orderNumber string
	)

	orderNumber = c.Param("order")

	switch donationType {
	case globals.PeriodicDonationType:
		d := models.PeriodicDonation{}
		err = mc.Storage.GetByConditions(map[string]interface{}{
			"order_number": orderNumber,
		}, &d)
		resp.BuildFromPeriodicDonationModel(d)
		_userID = uint(d.UserID)
		break
	case globals.PrimeDonationType:
		d := models.PayByPrimeDonation{}
		err = mc.Storage.GetByConditions(map[string]interface{}{
			"order_number": orderNumber,
		}, &d)
		resp.BuildFromPrimeDonationModel(d)
		_userID = uint(d.UserID)
		break
	case globals.OthersDonationType:
		d := models.PayByOtherMethodDonation{}
		err = mc.Storage.GetByConditions(map[string]interface{}{
			"order_number": orderNumber,
		}, &d)
		resp.BuildFromOtherMethodDonationModel(d)
		_userID = uint(d.UserID)
		break

	// TODO
	// case globasl.PayByCardTokenDonation:

	default:
		return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("donation type %s is not supported", donationType)}, nil
	}

	if err != nil {
		if storage.IsNotFound(err) {
			return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{
				"req.URL": fmt.Sprintf("%s cannot address a found resource", c.Request.RequestURI),
			}}, nil
		}
		return toResponse(err)
	}

	// Compare with the auth-user-id in context extracted from access_token
	authUserID = c.Request.Context().Value(globals.AuthUserIDProperty)

	if fmt.Sprint(_userID) != fmt.Sprint(authUserID) {
		return http.StatusForbidden, gin.H{"status": "fail", "data": gin.H{
			"req.Headers.Authorization": fmt.Sprintf("%s is forbidden to access", c.Request.RequestURI),
		}}, nil
	}

	return http.StatusOK, gin.H{"status": "success", "data": resp}, nil
}

func (mc *MembershipController) GetVerificationInfoOfADonation(c *gin.Context) (int, gin.H, error) {
	var d models.PayByPrimeDonation

	orderNumber := c.Param("order")
	err := mc.Storage.GetByConditions(map[string]interface{}{
		"order_number": orderNumber,
	}, &d)
	_userID := uint(d.UserID)

	if err != nil {
		if storage.IsNotFound(err) {
			return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{
				"req.URL": fmt.Sprintf("%s cannot address a found resource", c.Request.RequestURI),
			}}, nil
		}
		return toResponse(err)
	}

	// Compare with the auth-user-id in context extracted from access_token
	authUserID := c.Request.Context().Value(globals.AuthUserIDProperty)

	if fmt.Sprint(_userID) != fmt.Sprint(authUserID) {
		return http.StatusForbidden, gin.H{"status": "fail", "data": gin.H{
			"req.Headers.Authorization": fmt.Sprintf("%s is forbidden to access", c.Request.RequestURI),
		}}, nil
	}

	return http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"rec_trade_id":        d.RecTradeID,
		"bank_transaction_id": d.BankTransactionID,
		"status":              d.Status,
	}}, nil
}

func (mc *MembershipController) PatchLinePayOfAUser(c *gin.Context) (int, gin.H, error) {
	var callbackPayload tapPayTransactionResp

	if failData, valid := bindRequestJSONBody(c, &callbackPayload); valid == false {
		log.Infof("Fail to bind callback payload, %v", failData)
		return http.StatusBadRequest, gin.H{}, nil
	}

	// Validate Line Pay Method if PayInfo is set
	if callbackPayload.PayInfo.Method.IsZero() == false {
		if valid := validateLinePayMethod(callbackPayload.PayInfo.Method.String); valid == false {
			log.Infof("Invalid line pay method %s, should be %s", callbackPayload.PayInfo.Method.String, strings.Join(linePayMethods, ","))
			return http.StatusBadRequest, gin.H{}, nil
		}
	}

	if linePayMethodCreditCard == callbackPayload.PayInfo.Method.String {
		// Validate Line Pay Masked Credit Card Number format
		// sample: ************1234
		re := regexp.MustCompile("^[\\*]{12}[\\d]{4}$")

		if re.MatchString(callbackPayload.PayInfo.MaskedCreditCardNumber.String) == false {
			log.Infof("Invalid line pay credit number format: %s", callbackPayload.PayInfo.MaskedCreditCardNumber.String)
			return http.StatusBadRequest, gin.H{}, nil
		}
	}

	updateData := models.PayByPrimeDonation{}
	if tapPayRespStatusSuccess == callbackPayload.Status {
		callbackPayload.AppendLinePayOnPrimeDonation(&updateData, statusPaid)
	} else {
		callbackPayload.AppendLinePayOnPrimeDonation(&updateData, statusFail)
	}
	conditions := map[string]interface{}{
		"order_number":        callbackPayload.OrderNumber,
		"rec_trade_id":        callbackPayload.TappayResp.RecTradeID,
		"bank_transaction_id": callbackPayload.TappayResp.BankTransactionID,
		"amount":              callbackPayload.Amount,
	}
	err, rowsAffected := mc.Storage.UpdateByConditions(conditions, updateData)

	switch {
	case err != nil:
		return http.StatusInternalServerError, gin.H{}, err
	case rowsAffected == 0:
		log.Infof("No corresponding record to patch, condition: %v", conditions)
		return http.StatusUnprocessableEntity, gin.H{}, nil
	}

	if updateData.Status == statusPaid {
		var d models.PayByPrimeDonation
		mail := new(clientResp)

		mc.Storage.GetByConditions(map[string]interface{}{
			"order_number": callbackPayload.OrderNumber,
		}, &d)
		mail.BuildFromPrimeDonationModel(d)

		go mc.sendDonationThankYouMail(*mail)
	}

	return http.StatusNoContent, gin.H{}, nil
}

func (mc *MembershipController) QueryTappayServer(c *gin.Context) (int, gin.H, error) {
	var (
		reqBody queryReq
		d       models.PayByPrimeDonation
	)

	if failData, valid := bindRequestJSONBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	err := mc.Storage.GetByConditions(map[string]interface{}{
		"order_number": reqBody.Filters.OrderNumber,
	}, &d)
	recordUserID := uint(d.UserID)

	if err != nil {
		if storage.IsNotFound(err) {
			return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{
				"filters": fmt.Sprintf("%v cannot address a found resource", reqBody.Filters),
			}}, nil
		}
		return toResponse(err)
	}

	// Compare with the auth-user-id in context extracted from access_token
	authUserID := c.Request.Context().Value(globals.AuthUserIDProperty)

	if fmt.Sprint(recordUserID) != fmt.Sprint(authUserID) {
		return http.StatusForbidden, gin.H{"status": "fail", "data": gin.H{
			"req.Headers.Authorization": fmt.Sprintf("%s is forbidden to access", c.Request.RequestURI),
		}}, nil
	}

	// If the required fields `bank_transaction_id` or `rec_trade_id` are not specified, append the time filter of the 90 days interval from now
	if reqBody.Filters.RecTradeID.IsZero() && reqBody.Filters.BankTransactionID.IsZero() {
		end := time.Now()
		start := end.AddDate(0, 0, -90)
		reqBody.Filters.Time.EndTime = null.IntFrom(end.Unix() * 1000)
		reqBody.Filters.Time.StartTime = null.IntFrom(start.Unix() * 1000)
	}
	tapPayResp, err := reqBody.QueryServer()

	if err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()}, nil
	}

	return http.StatusOK, gin.H{"status": "success", "data": tapPayResp}, nil
}

func (resp tapPayTransactionResp) AppendRespOnPrimeDonation(m *models.PayByPrimeDonation, status string) {
	m.CardInfo = resp.CardInfo
	m.TappayResp = resp.TappayResp
	m.TappayApiStatus = null.IntFrom(resp.Status)

	if resp.TransactionTimeMillis > 0 {
		ttm := time.Unix(resp.TransactionTimeMillis/secToMsec, (resp.TransactionTimeMillis%secToMsec)*msecToNanosec)
		m.TransactionTime = null.TimeFrom(ttm)
	}

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

	m.Status = status
}

func (q queryReq) QueryServer() (tapPayTransactionRecordResp, error) {
	client := getProxyHttpClient()

	tapPayReq := q.BuildTapPayRecordReq()
	reqBodyJson, _ := json.Marshal(&tapPayReq)
	req, _ := http.NewRequest(http.MethodPost, globals.Conf.Donation.TapPayRecordURL, bytes.NewBuffer(reqBodyJson))
	req.Header.Add("x-api-key", globals.Conf.Donation.TapPayPartnerKey)
	req.Header.Add("Content-Type", "application/json")

	rawResp, err := client.Do(req)

	if err != nil {
		return tapPayTransactionRecordResp{}, errors.New("cannot request to tap pay server")
	}
	defer rawResp.Body.Close()

	body, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return tapPayTransactionRecordResp{}, errors.New("Cannot read response from tap pay server")
	}

	resp := tapPayTransactionRecordResp{}

	json.Unmarshal(body, &resp)

	return resp, nil

}

func (q queryReq) BuildTapPayRecordReq() (t tapPayTransactionRecordReq) {
	const defaultRecordPerPage = 1

	t.PartnerKey = globals.Conf.Donation.TapPayPartnerKey

	if q.RecordsPerPage == 0 {
		t.RecordsPerPage = defaultRecordPerPage
	} else {
		t.RecordsPerPage = q.RecordsPerPage
	}
	t.Filters = q.Filters

	return
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

func (resp tapPayTransactionResp) AppendRespOnTokenDonation(m *models.PayByCardTokenDonation, status string) {
	m.TappayResp = resp.TappayResp
	m.TappayApiStatus = null.IntFrom(resp.Status)

	if resp.TransactionTimeMillis > 0 {
		ttm := time.Unix(resp.TransactionTimeMillis/secToMsec, (resp.TransactionTimeMillis%secToMsec)*msecToNanosec)
		m.TransactionTime = null.TimeFrom(ttm)
	}

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

	m.Status = status
}

func (resp tapPayTransactionResp) AppendLinePayOnPrimeDonation(m *models.PayByPrimeDonation, status string) {
	m.PayInfo = resp.PayInfo
	m.TappayApiStatus = null.IntFrom(resp.Status)

	if resp.PayInfo.Method.String == linePayMethodCreditCard {
		m.CardInfo.LastFour = null.StringFrom(strings.Replace(resp.PayInfo.MaskedCreditCardNumber.String, "*", "", -1))
	}
	m.BankResultMsg = resp.BankResultMsg
	m.BankResultCode = resp.BankResultCode
	m.Status = status
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
		log.Infof("%+v", errors.WithStack(err))
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
		log.Infof("%+v", errors.Wrap(err, "cannot create a block cipher with the given key"))
		return data
	}

	nonce := make([]byte, gcm.NonceSize())

	// generate random nonce for encryption
	if _, err := io.ReadFull(rand.Reader, nonce); nil != err {
		//fallback to store plaintext instead
		log.Infof("%+v", errors.Wrap(err, "cannot generate a nonce"))
		return data
	}

	// prepend the cipher with the random nonce
	return string(gcm.Seal(nonce, nonce, []byte(data), nil))
}

func generateOrderNumber(t payType, payMethodID int) string {
	timestamp := time.Now().UnixNano()
	orderNumber := fmt.Sprintf("%s-%d%d%d", orderPrefix, timestamp, t, payMethodID)
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

func getProxyHttpClient() *http.Client {
	const defaultRequestTimeout = 45 * time.Second

	client := &http.Client{Timeout: defaultRequestTimeout}

	// Prior to route through proxy for http request if a proxy server is configured
	if len(globals.Conf.Donation.ProxyServer) > 0 {
		proxyUrl, _ := url.Parse(globals.Conf.Donation.ProxyServer)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	return client
}

func handleTapPayBodyParseError(body []byte) (tapPayTransactionResp, error) {
	var minResp tapPayMinTransactionResp
	var resp tapPayTransactionResp
	var err error

	if err = json.Unmarshal(body, &minResp); nil != err {
		return tapPayTransactionResp{}, errors.Wrap(err, "Cannot unmarshal json response from tap pay server")
	}

	if tapPayRespStatusSuccess != minResp.Status {
		return tapPayTransactionResp{}, errors.Wrap(err, fmt.Sprintf("Cannot make success transaction on tap pay, msg: %s", minResp.Msg))
	}

	resp.Status = minResp.Status
	resp.Msg = minResp.Msg

	return resp, nil
}

func serveHttp(key string, reqBodyJson []byte) (tapPayTransactionResp, error) {
	client := getProxyHttpClient()

	req, _ := http.NewRequest("POST", globals.Conf.Donation.TapPayURL, bytes.NewBuffer(reqBodyJson))
	req.Header.Add("x-api-key", key)
	req.Header.Add("Content-Type", "application/json")

	rawResp, err := client.Do(req)

	// If fail to sending request
	if nil != err {
		return tapPayTransactionResp{}, errors.Wrap(err, "cannot request to tap pay server")
	}
	defer rawResp.Body.Close()

	// If timeout or other errors occur during reading the body...
	// TODO: Might require a mechanism to notify users
	body, err := ioutil.ReadAll(rawResp.Body)
	if nil != err {
		return tapPayTransactionResp{}, errors.Wrap(err, "Cannot read response from tap pay server")
	}

	resp := tapPayTransactionResp{}

	err = json.Unmarshal(body, &resp)

	switch {
	case nil != err:
		return handleTapPayBodyParseError(body)
	case tapPayRespStatusSuccess != resp.Status:
		return resp, errors.New(fmt.Sprintf("Cannot make success transaction on tap pay, msg: %s", resp.Msg))
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

func validateLinePayMethod(method string) bool {
	valid := false

	for _, m := range linePayMethods {
		if m == method {
			valid = true
			break
		}
	}

	return valid
}
