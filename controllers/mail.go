package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"

	"twreporter.org/go-api/services"
	"twreporter.org/go-api/utils"
)

type activationReqBody struct {
	Email        string `json:"email" binding:"required"`
	ActivateLink string `json:"activate_link" binding:"required"`
}

type donationSuccessReqBody struct {
	Address           string   `json:"address"`
	Amount            uint     `json:"amount" binding:"required"`
	CardInfoLastFour  string   `json:"card_info_last_four"`
	CardInfoType      string   `json:"card_info_type"`
	Currency          string   `json:"currency"`
	DonationTimestamp null.Int `json:"donation_timestamp"`
	DonationLink      string   `json:"donation_link" binding:"required"`
	DonationMethod    string   `json:"donation_method" binding:"required"`
	DonationType      string   `json:"donation_type" binding:"required"`
	Email             string   `json:"email" binding:"required"`
	Name              string   `json:"name"`
	NationalID        string   `json:"national_id"`
	OrderNumber       string   `json:"order_number" binding:"required"`
	PhoneNumber       string   `json:"phone_number"`
}

// NewMailController is used to new *MailController
func NewMailController(svc services.MailService, t *template.Template) *MailController {
	return &MailController{
		HTMLTemplate: t,
		MailService:  svc,
	}
}

// MailController is the data structure holding HTML template and mail service
type MailController struct {
	HTMLTemplate *template.Template
	MailService  services.MailService
}

// LoadTemplateFiles is a wrapper function to parse template files
func (contrl *MailController) LoadTemplateFiles(filenames ...string) {
	contrl.HTMLTemplate = template.Must(template.ParseFiles(filenames...))
}

// SendActivation retrieves email and activation link from rqeuest body,
// and invoke MailService to send activation mail
func (contrl *MailController) SendActivation(c *gin.Context) (int, gin.H, error) {
	const subject = "歡迎登入報導者，體驗會員專屬功能"
	var err error
	var failData gin.H
	var mailBody string
	var out bytes.Buffer
	var reqBody activationReqBody
	var valid bool

	if failData, valid = bindRequestJSONBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	if err = contrl.HTMLTemplate.ExecuteTemplate(&out, "signin.tmpl", struct {
		Href string
	}{
		reqBody.ActivateLink,
	}); err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "can not create activate mail body"}, errors.WithStack(err)
	}

	mailBody = out.String()

	if err = contrl.MailService.Send(reqBody.Email, subject, mailBody); err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("can not send activate mail to %s", reqBody.Email)}, err
	}

	return http.StatusNoContent, gin.H{}, nil
}

func (contrl *MailController) SendDonationSuccessMail(c *gin.Context) (int, gin.H, error) {
	const subject = "感謝您成為報導者的贊助夥伴"
	const taipeiLocationName = "Asia/Taipei"
	var donationDatetime time.Time
	var err error
	var failData gin.H
	var location *time.Location
	var mailBody string
	var out bytes.Buffer
	var reqBody donationSuccessReqBody
	var valid bool

	// parse requst JSON into struct
	if failData, valid = bindRequestJSONBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	if reqBody.Currency == "" {
		// give default Currency
		reqBody.Currency = "TWD"
	}

	if reqBody.DonationTimestamp.Valid {
		donationDatetime = time.Unix(reqBody.DonationTimestamp.Int64, 0)
	} else {
		donationDatetime = time.Now()
	}

	location, _ = time.LoadLocation(taipeiLocationName)

	var templateData = struct {
		donationSuccessReqBody
		DonationDatetime string
	}{
		reqBody,
		donationDatetime.In(location).Format("2006-01-02 15:04:05 UTC+8"),
	}

	if err = contrl.HTMLTemplate.ExecuteTemplate(&out, "success-donation.tmpl", templateData); err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "can not create donation success mail body"}, errors.WithStack(err)
	}

	mailBody = out.String()

	// send email through mail service
	if err = contrl.MailService.Send(reqBody.Email, subject, mailBody); err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("can not send donation success mail to %s", reqBody.Email)}, err
	}

	return http.StatusNoContent, gin.H{}, nil
}

func postMailServiceEndpoint(reqBody interface{}, endpoint string) error {
	var body []byte
	var err error
	var rawResp *http.Response
	var timeout = 10 * time.Second
	var expiration int = 60
	var accessToken string

	if body, err = json.Marshal(reqBody); err != nil {
		return errors.WithStack(err)
	}

	// Setup HTTP client with timeout
	client := &http.Client{Timeout: timeout}

	accessToken, _ = utils.RetrieveMailServiceAccessToken(expiration)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	if rawResp, err = client.Do(req); err != nil {
		return errors.WithStack(err)
	}

	defer rawResp.Body.Close()

	if rawResp.StatusCode != http.StatusNoContent {
		if body, err = ioutil.ReadAll(rawResp.Body); err != nil {
			return errors.WithStack(err)
		}

		errMsg := fmt.Sprintf("receive error status code(%d) from %s. error response: %s", rawResp.StatusCode, endpoint, string(body))
		return errors.New(errMsg)
	}

	return nil
}
