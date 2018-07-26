package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

type (
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

	tapPayByPrimeReqBody struct {
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

	tapPayByPrimeResp struct {
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

	payType int
)

const (
	defaultDetails    = "報導者小額捐款"
	defaultCurrency   = "TWD"
	defaultMerchantID = "twreporter_CTBC"

	invalidPayMethodID = -1

	orderPrefix = "twreporter"
)

// pay type Enum
const (
	oneTime payType = iota + 1
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

func (mc *MembershipController) CreateADonationOfAUser(c *gin.Context) (int, gin.H, error) {
	type primeReqBody struct {
		Prime       string            `json:"prime" form:"prime" binding:"required"`
		Amount      uint              `json:"amount" form:"amount" binding:"required"`
		Currency    string            `json:"currency" form:"currency"`
		Details     string            `json:"details" form:"details"`
		Cardholder  models.Cardholder `json:"cardholder" form:"cardholder" binding:"required,dive"`
		OrderNumber string            `json:"order_number" form:"order_number"`
		MerchantID  string            `json:"merchant_id" form:"merchant_id"`
		ResultUrl   linePayResultUrl  `json:"result_url" form:"result_url"`
	}

	type primeRespBody struct {
		IsPeriodic  bool              `json:"is_periodic"`
		PayMethod   string            `json:"pay_method"`
		CardInfo    models.CardInfo   `json:"card_info"`
		Cardholder  models.Cardholder `json:"cardholder"`
		Amount      uint              `json:"amount"`
		Currency    string            `json:"currency"`
		Details     string            `json:"details"`
		OrderNumber string            `json:"order_number"`
	}
	const errorWhere = "MembershipController.CreateADonationOfAUser"

	// Validate pay_method
	payMethod := c.Param("pay_method")
	if err := validatePayMethod(payMethod); nil != err {
		return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"pay_method": err.Error()}}, nil
	}

	var reqBody primeReqBody

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

	/*userIDStr := c.Param("userID")
	if userID, err := strconv.ParseUint(userIDStr, 10, strconv.IntSize); nil != err {
		return http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"user_id": "Invalid user id: " + userIDStr}}, nil
	}*/

	// Build Tappay Pay By Prime Request
	tapPayReq := newTapPayByPrimeReq()

	// Fill up required fields
	tapPayReq.Prime = reqBody.Prime
	tapPayReq.Amount = reqBody.Amount
	tapPayReq.Cardholder = reqBody.Cardholder

	// Fill up optional fields
	if "" != reqBody.Currency {
		tapPayReq.Currency = reqBody.Currency
	}

	if "" != reqBody.Details {
		tapPayReq.Details = reqBody.Details
	}

	if "" != reqBody.MerchantID {
		tapPayReq.MerchantID = reqBody.MerchantID
	}

	if (linePayResultUrl{}) != reqBody.ResultUrl {
		tapPayReq.ResultUrl = reqBody.ResultUrl
	}

	if "" != reqBody.OrderNumber {
		tapPayReq.OrderNumber = reqBody.OrderNumber
	} else {
		tapPayReq.OrderNumber = generateOrderNumber(oneTime, getPayMethodID(payMethod))
	}
	tapPayReq.PartnerKey = utils.Cfg.DonationSettings.TapPayPartnerKey

	// Setup HTTP client
	client := &http.Client{}

	tapPayReqJson, _ := json.Marshal(tapPayReq)
	log.Info("TapPayUrl: " + utils.Cfg.DonationSettings.TapPayUrl)
	req, _ := http.NewRequest("POST", utils.Cfg.DonationSettings.TapPayUrl, bytes.NewBuffer(tapPayReqJson))
	req.Header.Add("x-api-key", tapPayReq.PartnerKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if nil != err {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Cannot request to Tap Pay Server", err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Cannot read response from Tap Pay Server", err.Error(), http.StatusInternalServerError)
	}

	var tapPayResp tapPayByPrimeResp
	err = json.Unmarshal(body, &tapPayResp)

	switch {
	case nil != err:
		return 0, gin.H{}, models.NewAppError(errorWhere, "Cannot unmarshal json response from Tap Pay Server", err.Error(), http.StatusInternalServerError)
	// TODO: Should deal with several Tap Pay error code
	case 0 != tapPayResp.Status:
		return 0, gin.H{}, models.NewAppError(errorWhere, "Cannot make success transaction on tap pay", "", http.StatusInternalServerError)
	default:
		// Omit intentionally
	}
	if nil != err {
		return 0, gin.H{}, models.NewAppError(errorWhere, "Cannot unmarshal json response from Tap Pay Server", err.Error(), http.StatusInternalServerError)
	}

	primeResp := primeRespBody{}
	primeResp.IsPeriodic = false
	primeResp.PayMethod = payMethod
	primeResp.CardInfo = tapPayResp.CardInfo
	primeResp.Cardholder = tapPayReq.Cardholder
	primeResp.Amount = tapPayResp.Amount
	primeResp.Currency = tapPayResp.Currency
	primeResp.Details = tapPayReq.Details
	primeResp.OrderNumber = tapPayReq.OrderNumber

	return http.StatusCreated, gin.H{"status": "success", "data": primeResp}, nil
}

//TODO
func (mc *MembershipController) GetDonationsOfAUser(c *gin.Context) (int, gin.H, error) {
	return 0, gin.H{}, nil
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

func newTapPayByPrimeReq() tapPayByPrimeReqBody {
	req := tapPayByPrimeReqBody{}
	req.Details = defaultDetails
	req.Currency = defaultCurrency
	req.MerchantID = defaultMerchantID
	return req
}

func validatePayMethod(payMethod string) error {
	switch {
	case invalidPayMethodID != getPayMethodID(payMethod):
		return nil
	default:
	}

	var errMsg bytes.Buffer
	errMsg.WriteString("Unsupported pay_method. Only support payment by")

	for k, v := range payMethodCollections {
		if 0 != k {
			errMsg.WriteString(", ")
		}
		errMsg.WriteString("'")
		errMsg.WriteString(v)
		errMsg.WriteString("'")
	}

	return errors.New(errMsg.String())
}
