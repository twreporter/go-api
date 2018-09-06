package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

var Defaults = struct {
	WebPushEndpoint string
}{
	WebPushEndpoint: "https://fcm.googleapis.com/fcm/send/f4Stnx6WC5s:APA91bFGo-JD8bDwezv1fx3RRyBVq6XxOkYIo8_7vCAJ3HFHLppKAV6GNmOIZLH0YeC2lM_Ifs9GkLK8Vi_8ASEYLBC1aU9nJy2rZSUfH7DE0AqIIbLrs93SdEdkwr5uL6skLPMjJsRQ",
}

// setDefaultWebPushSubscription - set up default records in web_push_subscriptions table
func setDefaultWebPushSubscription() {
	var path = "/v1/web-push/subscriptions"
	var webPush = WebPushSubscriptionPostBody{
		Endpoint:       Defaults.WebPushEndpoint,
		Keys:           "{\"p256dh\":\"BDmY8OGe-LfW0ENPIADvmdZMo3GfX2J2yqURpsDOn5tT8lQV-VVHyhRUgzjnmx_RRoobwdLULdBr26oULtLML3w\",\"auth\":\"P_AJ9QSqcgM-KJi_GRN3fQ\"}",
		ExpirationTime: "1526959900",
		UserID:         "1",
	}

	webPushJSON, _ := json.Marshal(webPush)
	ServeHTTP("POST", path, string(webPushJSON), "application/json", "")
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestMain(m *testing.M) {
	var err error

	fmt.Println("load default config")
	if globals.Conf, err = configs.LoadDefaultConf(); err != nil {
		panic(fmt.Sprintf("Can not load default config, but got err=%+v", err))
	}

	// Create DB connections
	if DB, err = OpenGormConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
	}

	DB.SetJoinTableHandler(&models.User{}, globals.TableBookmarks, &models.UsersBookmarks{})

	// Create Mongo DB connections
	if MgoDB, err = OpenMgoConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to mongo database, but got err=%+v", err))
	}
	// Set up database, including drop existing tables and create tables
	RunMigration()

	// Set up default records in tables
	SetDefaultRecords()

	SetupGinServer()

	// add default records in web_push_subscriptions table
	setDefaultWebPushSubscription()

	retCode := m.Run()
	os.Exit(retCode)
}
