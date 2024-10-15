package member_cms

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const receiptEndpoint = "/receipt"

func PostPrimeDonationReceipt(receiptNumber string) error {
	if len(receiptNumber) == 0 {
		return errors.New("order number is required")
	}

	url, err := GetApiBaseUrl()
	if err != nil {
		return err
	}
	url = url + receiptEndpoint

	payload := map[string]string{"receipt_number": receiptNumber}
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
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
