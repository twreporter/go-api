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

const defaultWebPushEndpoint = "https://fcm.googleapis.com/fcm/send/f4Stnx6WC5s:APA91bFGo-JD8bDwezv1fx3RRyBVq6XxOkYIo8_7vCAJ3HFHLppKAV6GNmOIZLH0YeC2lM_Ifs9GkLK8Vi_8ASEYLBC1aU9nJy2rZSUfH7DE0AqIIbLrs93SdEdkwr5uL6skLPMjJsRQ"

var isSetUp = false

// setDefaultWebPushSubscription - set up default records in web_push_subscriptions table
func setDefaultWebPushSubscription() {
	var path = "/v1/web-push/subscriptions"
	var webPush = webPushSubscriptionPostBody{
		Endpoint:       defaultWebPushEndpoint,
		Keys:           "{\"p256dh\":\"BDmY8OGe-LfW0ENPIADvmdZMo3GfX2J2yqURpsDOn5tT8lQV-VVHyhRUgzjnmx_RRoobwdLULdBr26oULtLML3w\",\"auth\":\"P_AJ9QSqcgM-KJi_GRN3fQ\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushJSON, _ := json.Marshal(webPush)
	serveHTTP("POST", path, string(webPushJSON), "application/json", "")
}

func setUpBeforeSubcriptionsTest() {
	if !isSetUp {
		setDefaultWebPushSubscription()
	}

	isSetUp = true
}

func TestIsWebPushSubscribed(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path string

	// set up before testing
	setUpBeforeSubcriptionsTest()

	path = fmt.Sprintf("/v1/web-push/subscriptions?endpoint=%v", defaultWebPushEndpoint)

	/** START - Read a web push subscription successfully **/

	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := wpSubResponse{}
	json.Unmarshal(body, &res)

	assert.Equal(t, defaultWebPushEndpoint, res.Data.Endpoint)

	/** END - Read a web push subscription successfully **/

	/** START - Fail to read a web push subscription **/

	// Situation 1: Endpoint query param is not provided
	path = "/v1/web-push/subscriptions?endpoint="
	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 404)

	// Situation 2: Endpoint is provided, but database does not have it
	path = "/v1/web-push/subscriptions?endpoint=http://web-push.subscriptions/endpoint-is-not-in-the-db"
	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 404)

	/** END - Fail to read a web push subscription **/
}

func TestSubscribeWebPush(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path = "/v1/web-push/subscriptions"
	var webPush webPushSubscriptionPostBody
	var webPushByteArray []byte

	// set up before testing
	setUpBeforeSubcriptionsTest()

	/** START - Add a web push subscription successfully **/
	webPush = webPushSubscriptionPostBody{
		Endpoint:       "http://web-push.subscriptions/new.endpoint.to.subscribe",
		Keys:           "{\"p256dh\":\"test-p256dh\",\"auth\":\"test-auth\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushByteArray, _ = json.Marshal(webPush)
	resp = serveHTTP("POST", path, string(webPushByteArray), "application/json", "")
	assert.Equal(t, resp.Code, 201)
	/** END - Add a web push subscription successfully **/

	/** START - Fail to add a web push subscription **/

	// Situation 1: POST Body is not fully provided, lack of `keys`
	webPush = webPushSubscriptionPostBody{
		Endpoint:       "http://web-push.subscriptions/another.endpoint.to.subscribe",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushByteArray, _ = json.Marshal(webPush)
	resp = serveHTTP("POST", path, string(webPushByteArray), "application/json", "")
	assert.Equal(t, resp.Code, 400)

	// Situation 2: Endpoint is already subscribed
	webPush = webPushSubscriptionPostBody{
		Endpoint:       defaultWebPushEndpoint,
		Keys:           "{\"p256dh\":\"test-p256dh\",\"auth\":\"test-auth\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushByteArray, _ = json.Marshal(webPush)
	resp = serveHTTP("POST", path, string(webPushByteArray), "application/json", "")
	assert.Equal(t, resp.Code, 409)

	/** END - Fail to add a web push subscription **/
}
