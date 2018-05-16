package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
)

type wpSubResponse struct {
	Status string                     `json:"status"`
	Data   models.WebPushSubscription `json:"data"`
}

type WebPushSubscriptionPostBody struct {
	Endpoint       string `json:"endpoint"`
	Keys           string `json:"keys"`
	ExpirationTime string `json:"expiration_time"`
	UserID         string `json:"user_id"`
}

func TestIsWebPushSubscribed(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path string

	path = fmt.Sprintf("/v1/web-push/subscriptions?endpoint=%v", Defaults.WebPushEndpoint)

	/** START - Read a web push subscription successfully **/

	resp = ServeHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := wpSubResponse{}
	json.Unmarshal(body, &res)

	assert.Equal(t, Defaults.WebPushEndpoint, res.Data.Endpoint)

	/** END - Read a web push subscription successfully **/

	/** START - Fail to read a web push subscription **/

	// Situation 1: Endpoint query param is not provided
	path = "/v1/web-push/subscriptions?endpoint="
	resp = ServeHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 404)

	// Situation 2: Endpoint is provided, but database does not have it
	path = "/v1/web-push/subscriptions?endpoint=http://web-push.subscriptions/endpoint-is-not-in-the-db"
	resp = ServeHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 404)

	/** END - Fail to read a web push subscription **/
}

func TestSubscribeWebPush(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path = "/v1/web-push/subscriptions"
	var webPush WebPushSubscriptionPostBody
	var webPushByteArray []byte

	/** START - Add a web push subscription successfully **/
	webPush = WebPushSubscriptionPostBody{
		Endpoint:       "http://web-push.subscriptions/new.endpoint.to.subscribe",
		Keys:           "{\"p256dh\":\"test-p256dh\",\"auth\":\"test-auth\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushByteArray, _ = json.Marshal(webPush)
	resp = ServeHTTP("POST", path, string(webPushByteArray), "application/json", "")
	assert.Equal(t, resp.Code, 201)
	/** END - Add a web push subscription successfully **/

	/** START - Fail to add a web push subscription **/

	// Situation 1: POST Body is not fully provided, lack of `keys`
	webPush = WebPushSubscriptionPostBody{
		Endpoint:       "http://web-push.subscriptions/another.endpoint.to.subscribe",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushByteArray, _ = json.Marshal(webPush)
	resp = ServeHTTP("POST", path, string(webPushByteArray), "application/json", "")
	assert.Equal(t, resp.Code, 400)

	// Situation 2: Endpoint is already subscribed
	webPush = WebPushSubscriptionPostBody{
		Endpoint:       Defaults.WebPushEndpoint,
		Keys:           "{\"p256dh\":\"test-p256dh\",\"auth\":\"test-auth\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushByteArray, _ = json.Marshal(webPush)
	resp = ServeHTTP("POST", path, string(webPushByteArray), "application/json", "")
	assert.Equal(t, resp.Code, 409)

	/** END - Fail to add a web push subscription **/
}
