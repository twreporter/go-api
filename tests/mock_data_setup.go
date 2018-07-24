package tests

import (
	"encoding/json"
)

// setDefaultWebPushSubscription - set up default records in web_push_subscriptions table
func setDefaultWebPushSubscription(wpe string) {
	var path = "/v1/web-push/subscriptions"
	var webPush = WebPushSubscriptionPostBody{
		Endpoint:       wpe,
		Keys:           "{\"p256dh\":\"BDmY8OGe-LfW0ENPIADvmdZMo3GfX2J2yqURpsDOn5tT8lQV-VVHyhRUgzjnmx_RRoobwdLULdBr26oULtLML3w\",\"auth\":\"P_AJ9QSqcgM-KJi_GRN3fQ\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushJSON, _ := json.Marshal(webPush)
	ServeHTTP("POST", path, string(webPushJSON), "application/json", "")
}

func CreateMockData() {
	// set up web push subscription default value
	setDefaultWebPushSubscription(Globs.Defaults.WebPushEndpoint)
}
