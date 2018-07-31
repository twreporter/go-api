package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

type (
	clientPrimeReq struct {
		Prime       string            `json:"prime" form:"prime" binding:"required"`
		Amount      uint              `json:"amount" form:"amount" binding:"required"`
		Currency    string            `json:"currency" form:"currency"`
		Details     string            `json:"details" form:"details"`
		Cardholder  models.Cardholder `json:"cardholder" form:"cardholder" binding:"required,dive"`
		OrderNumber string            `json:"order_number" form:"order_number"`
		MerchantID  string            `json:"merchant_id" form:"merchant_id"`
		ResultUrl   linePayResultUrl  `json:"result_url" form:"result_url"`
	}

	clientPrimeResp struct {
		IsPeriodic  bool              `json:"is_periodic"`
		PayMethod   string            `json:"pay_method"`
		CardInfo    models.CardInfo   `json:"card_info"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Amount      uint              `json:"amount"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		OrderNumber string            `json:"order_number"`
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
		Status                int                 `json:"status"`
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
		BankResultCode        string              `json:"bank_result_code"`
		BankResultMsg         string              `json:"bank_result_msg"`
	}

	tapPayMinTransactionResp struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	}

	payType int
)

const (
	defaultDetails        = "報導者小額捐款"
	defaultCurrency       = "TWD"
	defaultMerchantID     = "twreporter_CTBC"
	defaultRequestTimeout = 45 * time.Second

	invalidPayMethodID = -1

	orderPrefix = "twreporter"

	statusToPay  = "to_pay"
	statusPaying = "paying"
	statusPaid   = "paid"
	statusError  = "error"

	primeTableName = "pay_by_prime_donations"

	tapPayRespStatusSuccess = 0
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

//todo
func (mc *MembershipController) CreateAPeriodicDonationOfAUser(c *gin.Context) (int, gin.H, error) {
	return 0, gin.H{}, nil
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

	var reqBody clientPrimeReq

	// Validate request body
	if err := c.Bind(&reqBody); nil != err {
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
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	userID, _ := strconv.ParseUint(c.Param("userID"), 10, strconv.IntSize)

	// Start Tappay transaction

	// Build Tappay pay by prime request
	tapPayReq := buildTapPayPrimeReq(payMethod, reqBody)

	draftRecord := buildPrimeDraftRecord(uint(userID), payMethod, tapPayReq)

	if err := mc.Storage.CreateAPayByPrimeDonation(draftRecord); nil != err {
		switch appErr := err.(type) {
		case *models.AppError:
			return 0, gin.H{}, models.NewAppError(errorWhere, "Fails to create a draft prime record", appErr.Error(), appErr.StatusCode)
		default:
			return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Unknown Error Type. Fails to create a draft prime record. %v", err.Error())}, nil
		}
	}

	tapPayReqJson, _ := json.Marshal(tapPayReq)

	// Update the transaction status to 'paying'
	if err := mc.Storage.UpdateTransactionStatus(tapPayReq.OrderNumber, statusPaying, primeTableName); nil != err {
		log.Error(err.Error())
		// Proceed even if the status update failed
	}

	tapPayResp, err := serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		if tapPayRespStatusSuccess != tapPayResp.Status {
			// If tappay error occurs, update the transaction status to 'error'
			mc.Storage.UpdateAPayByPrimeDonation(tapPayReq.OrderNumber, models.PayByPrimeDonation{
				ThirdPartyStatus: tapPayResp.Status,
				Msg:              tapPayResp.Msg,
				Status:           statusError,
			})
		}
		return 0, gin.H{}, models.NewAppError(errorWhere, err.Error(), "", http.StatusInternalServerError)
	}

	// Update the transaction status to 'paid' if transaction succeeds
	updateRecord := buildPrimeSuccessRecord(tapPayResp)

	if err := mc.Storage.UpdateAPayByPrimeDonation(tapPayReq.OrderNumber, updateRecord); nil != err {
		log.Error(err.Error())
	}

	clientResp := buildClientPrimeResp(payMethod, tapPayReq, tapPayResp)

	return http.StatusCreated, gin.H{"status": "success", "data": clientResp}, nil
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

func buildClientPrimeResp(payMethod string, req tapPayPrimeReq, resp tapPayTransactionResp) clientPrimeResp {
	c := clientPrimeResp{}
	c.IsPeriodic = false
	c.PayMethod = payMethod

	c.CardInfo = resp.CardInfo
	c.Cardholder = req.Cardholder
	c.OrderNumber = req.OrderNumber

	c.Amount = resp.Amount
	c.Currency = resp.Currency
	c.Details = req.Details

	return c
}

func buildPrimeDraftRecord(userID uint, payMethod string, req tapPayPrimeReq) models.PayByPrimeDonation {
	m := models.PayByPrimeDonation{}

	m.UserID = userID
	m.PayMethod = payMethod

	m.CardholderEmail = req.Cardholder.Email
	m.CardholderPhoneNumber = &req.Cardholder.PhoneNumber
	m.CardholderName = &req.Cardholder.Name
	m.CardholderZipCode = &req.Cardholder.ZipCode
	m.CardholderAddress = &req.Cardholder.Address
	m.CardholderNationalID = &req.Cardholder.NationalID

	m.Details = req.Details
	m.MerchantID = req.MerchantID
	m.OrderNumber = req.OrderNumber

	m.Status = statusToPay

	return m
}

func buildPrimeSuccessRecord(resp tapPayTransactionResp) models.PayByPrimeDonation {
	m := models.PayByPrimeDonation{}

	m.ThirdPartyStatus = resp.Status
	m.Msg = resp.Msg
	m.RecTradeID = resp.RecTradeID
	m.BankTransactionID = resp.BankTransactionID
	m.AuthCode = resp.AuthCode
	m.Acquirer = resp.Acquirer
	m.BankResultCode = &resp.BankResultCode
	m.BankResultMsg = &resp.BankResultMsg

	ttm := time.Unix(resp.TransactionTimeMillis/1000, resp.TransactionTimeMillis%1000)
	m.TransactionTime = &ttm

	t, err := strconv.ParseInt(resp.BankTransactionTime.StartTimeMillis, 10, strconv.IntSize)
	if nil == err {
		stm := time.Unix(t/1000, t%1000)
		m.BankTransactionStartTime = &stm
	}

	t, err = strconv.ParseInt(resp.BankTransactionTime.EndTimeMillis, 10, strconv.IntSize)
	if nil == err {
		etm := time.Unix(t/1000, t%1000)
		m.BankTransactionEndTime = &etm
	}

	m.CardInfoBinCode = &resp.CardInfo.BinCode
	m.CardInfoLastFour = &resp.CardInfo.LastFour
	m.CardInfoIssuer = &resp.CardInfo.Issuer
	m.CardInfoFunding = &resp.CardInfo.Funding
	m.CardInfoType = &resp.CardInfo.Type
	m.CardInfoLevel = &resp.CardInfo.Level
	m.CardInfoCountry = &resp.CardInfo.Country
	m.CardInfoCountryCode = &resp.CardInfo.CountryCode
	m.CardInfoExpiryDate = &resp.CardInfo.ExpiryDate

	m.Status = "paid"

	return m
}

func buildTapPayPrimeReq(payMethod string, req clientPrimeReq) tapPayPrimeReq {
	t := &tapPayPrimeReq{}

	t.setDefault()

	// Fill up required fields
	t.Prime = req.Prime
	t.Amount = req.Amount
	t.Cardholder = req.Cardholder

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
		t.OrderNumber = generateOrderNumber(oneTime, getPayMethodID(payMethod))
	}
	t.PartnerKey = utils.Cfg.DonationSettings.TapPayPartnerKey

	return *t
}

func generateOrderNumber(t payType, payMethodID int) string {
	timestamp := time.Now().Unix()
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

	if err := json.Unmarshal(body, &minResp); nil != err {
		return tapPayTransactionResp{}, errors.New("Cannot unmarshal json response from tap pay server")
	}

	resp.Status = minResp.Status
	resp.Msg = minResp.Msg

	return resp, nil
}

func serveHttp(key string, reqBodyJson []byte) (tapPayTransactionResp, error) {
	// Setup HTTP client with timeout
	client := &http.Client{Timeout: defaultRequestTimeout}

	req, _ := http.NewRequest("POST", utils.Cfg.DonationSettings.TapPayUrl, bytes.NewBuffer(reqBodyJson))
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

	var resp tapPayTransactionResp
	err = json.Unmarshal(body, &resp)

	switch {
	case nil != err:
		log.Error(err.Error())
		return handleTapPayBodyParseError(body)
	case 0 != resp.Status:
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
