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
	const errorWhere = "MembershipController.CreateADonationOfAUser"

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

	//userID, _ := strconv.ParseUint(c.Param("userID"), 10, strconv.IntSize)

	// Build Tappay pay by prime request
	tapPayReq := buildTapPayPrimeReq(payMethod, reqBody)

	tapPayReqJson, _ := json.Marshal(tapPayReq)

	tapPayResp, err := serveHttp(tapPayReq.PartnerKey, tapPayReqJson)

	if nil != err {
		return 0, gin.H{}, models.NewAppError(errorWhere, err.Error(), "", http.StatusInternalServerError)
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

func serveHttp(key string, reqBodyJson []byte) (tapPayTransactionResp, error) {
	// Setup HTTP client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", utils.Cfg.DonationSettings.TapPayUrl, bytes.NewBuffer(reqBodyJson))
	req.Header.Add("x-api-key", key)
	req.Header.Add("Content-Type", "application/json")

	rawResp, err := client.Do(req)
	if nil != err {
		log.Error(err.Error())
		return tapPayTransactionResp{}, errors.New("cannot request to tap pay server")
	}
	defer rawResp.Body.Close()
	body, err := ioutil.ReadAll(rawResp.Body)
	if nil != err {
		log.Error(err.Error())
		return tapPayTransactionResp{}, errors.New("Cannot read response from Tap Pay Server")
	}

	var resp tapPayTransactionResp
	err = json.Unmarshal(body, &resp)

	switch {
	case nil != err:
		return tapPayTransactionResp{}, errors.New("Cannot unmarshal json response from Tap Pay Server")
	case 0 != resp.Status:
		log.Error("tap pay msg: " + resp.Msg)
		return tapPayTransactionResp{}, errors.New("Cannot make success transaction on tap pay")
	default:
		// Omit intentionally
	}

	return resp, nil
}

func validatePayMethod(payMethod string) error {
	if invalidPayMethodID != getPayMethodID(payMethod) {
		return nil
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
