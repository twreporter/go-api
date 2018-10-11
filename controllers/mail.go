package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"twreporter.org/go-api/services"
)

type activationReqBody struct {
	Email        string `json:"email" binding:"required"`
	ActivateLink string `json:"activate_link" binding:"required"`
}

type donationSuccessReqBody struct {
	Address          string `json:"address"`
	Amount           uint   `json:"amount" binding:"required"`
	CardInfoLastFour string `json:"card_info_four_number"`
	CardInfoType     string `json:"card_info_type"`
	Currency         string `json:"currency"`
	DonationDatetime string `json:"donation_datetime"`
	DonationLink     string `json:"donation_link"`
	DonationMethod   string `json:"donation_method" binding:"required"`
	DonationType     string `json:"donation_type" binding:"required"`
	Email            string `json:"email" binding:"required"`
	Name             string `json:"name"`
	NationalID       string `json:"national_id"`
	OrderNumber      string `json:"order_number" binding:"required"`
	PhoneNumber      string `json:"phone_number"`
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
	const subject = "登入報導者"
	var err error
	var failData gin.H
	var mailBody string
	var out bytes.Buffer
	var reqBody activationReqBody
	var valid bool

	if failData, valid = bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	if err = contrl.HTMLTemplate.ExecuteTemplate(&out, "signin.tmpl", struct {
		Href string
	}{
		reqBody.ActivateLink,
	}); err != nil {
		log.Error(err)
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "can not create activate mail body"}, nil
	}

	mailBody = out.String()

	if err = contrl.MailService.Send(reqBody.Email, subject, mailBody); err != nil {
		log.Error(err)
		return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("can not send activate mail to %s", reqBody.Email)}, nil
	}

	return http.StatusNoContent, gin.H{}, nil
}

func (contrl *MailController) SendDonationSuccessMail(c *gin.Context) (int, gin.H, error) {
	const subject = "感謝您成為報導者的夥伴"
	var err error
	var failData gin.H
	var mailBody string
	var out bytes.Buffer
	var reqBody donationSuccessReqBody
	var valid bool

	// parse requst JSON into struct
	if failData, valid = bindRequestBody(c, &reqBody); valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": failData}, nil
	}

	if reqBody.Currency == "" {
		// give default Currency
		reqBody.Currency = "TWD"
	}

	if reqBody.DonationDatetime == "" {
		// give default DonationDatetime
		reqBody.DonationDatetime = time.Now().Format("2006-01-02 15:04:05")
	}

	if err = contrl.HTMLTemplate.ExecuteTemplate(&out, "success-donation.tmpl", reqBody); err != nil {
		log.Error(err)
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "can not create donation success mail body"}, nil
	}

	mailBody = out.String()

	// send email through mail service
	if err = contrl.MailService.Send(reqBody.Email, subject, mailBody); err != nil {
		log.Error(err)
		return http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("can not send donation success mail to %s", reqBody.Email)}, nil
	}

	return http.StatusNoContent, gin.H{}, nil
}

func postMailServiceEndpoint(reqBody interface{}, endpoint string) error {
	var body []byte
	var err error
	var rawResp *http.Response
	var timeout = 10 * time.Second

	if body, err = json.Marshal(reqBody); err != nil {
		return err
	}

	// Setup HTTP client with timeout
	client := &http.Client{Timeout: timeout}

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	if rawResp, err = client.Do(req); err != nil {
		return err
	}

	defer rawResp.Body.Close()

	if rawResp.StatusCode != http.StatusNoContent {
		if body, err = ioutil.ReadAll(rawResp.Body); err != nil {
			return err
		}
		errMsg := fmt.Sprintf("receive error status code from %s. error response: %s", endpoint, string(body))
		return errors.New(errMsg)
	}

	return nil
}
