package member_cms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/twreporter/go-api/globals"
)

const receiptEndpoint = "/receipt"
const yearlyReceiptPath = "yearly"

func GetPrimeDonationReceiptRequest(receiptNumber string) (*http.Request, error) {
	if !globals.Conf.Features.MemberCMS {
		return nil, errors.New("disable intergrating with member cms")
	}
	if len(receiptNumber) == 0 {
		return nil, errors.New("receipt numner is required")
	}

	url, err := GetApiBaseUrl()
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("%s%s/%s", url, receiptEndpoint, receiptNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	AppendRequiredHeader(req)

	return req, nil
}

func PostPrimeDonationReceipt(receiptNumber string, orderNumber string) error {
	if !globals.Conf.Features.MemberCMS {
		return errors.New("disable intergrating with member cms")
	}
	if len(receiptNumber) == 0 && len(orderNumber) == 0 {
		return errors.New("one of receipt number or order number should be provided")
	}

	url, err := GetApiBaseUrl()
	if err != nil {
		return err
	}
	url = url + receiptEndpoint

	payload := map[string]string{"receipt_number": receiptNumber, "order_number": orderNumber}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	AppendRequiredHeader(req)

	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("\nresp:\n%+v\n", resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func GetYearlyReceiptRequest(email string, year string) (*http.Request, error) {
	if !globals.Conf.Features.MemberCMS {
		return nil, errors.New("disable intergrating with member cms")
	}
	if len(email) == 0 {
		return nil, errors.New("email is required")
	}

	url, err := GetApiBaseUrl()
	if err != nil {
		return nil, err
	}
	url = fmt.Sprintf("%s%s/%s/%s/%s", url, receiptEndpoint, yearlyReceiptPath, email, year)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	AppendRequiredHeader(req)

	return req, nil
}
