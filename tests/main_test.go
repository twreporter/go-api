package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/twreporter/go-api/configs"
	"github.com/twreporter/go-api/controllers"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/routers"
	"github.com/twreporter/go-api/storage"
	"github.com/twreporter/go-api/utils"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

var Globs globalVariables

var testMongoClient *mongodriver.Client

func init() {
	var defaults = defaultVariables{
		Account: "developer@twreporter.org",
		Service: "default_service",
		Token:   "default_token",

		ErrorEmailAddress: "error@twreporter.org",
	}

	Globs = globalVariables{
		Defaults: defaults,
	}
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

func generateIDToken(user models.User) (jwt string) {
	jwt, _ = utils.RetrieveV2IDToken(user.ID, user.Email.ValueOrZero(), user.FirstName.ValueOrZero(), user.LastName.ValueOrZero(), 3600)
	return
}

func getReporterAccount(email string) (ra models.ReporterAccount) {
	as := storage.NewGormStorage(Globs.GormDB)
	ra, _ = as.GetReporterAccountData(email)
	return ra
}

func createUser(email string) models.User {
	as := storage.NewGormStorage(Globs.GormDB)

	ra := models.ReporterAccount{
		Email:         email,
		ActivateToken: Globs.Defaults.Token,
		ActExpTime:    time.Now().Add(time.Duration(15) * time.Minute),
	}

	user, _ := as.InsertUserByReporterAccount(ra)

	return user
}

func deleteUser(user models.User) {
	db := Globs.GormDB

	// Remove corresponding reporter account
	db.Unscoped().Delete(user.ReporterAccount)
	db.Unscoped().Delete(user)
}

func getUser(email string) (user models.User) {
	as := storage.NewGormStorage(Globs.GormDB)
	user, _ = as.GetUserByEmail(email)
	return
}

func serveHTTP(method, path, body, contentType, authorization string) (resp *httptest.ResponseRecorder) {
	var req *http.Request

	req = requestWithBody(method, path, body)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	resp = httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)

	return
}

func serveHTTPWithCookies(method, path, body, contentType, authorization string, cookies ...http.Cookie) (resp *httptest.ResponseRecorder) {
	var req *http.Request

	req = requestWithBody(method, path, body)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	resp = httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)

	return
}

type mockMailStrategy struct{}

func (s mockMailStrategy) Send(to, subject, body string) error {
	if to == Globs.Defaults.ErrorEmailAddress {
		return errors.New("mail service works abnormally")
	}
	return nil
}

type mockIndexSearcher struct{}

func (mockIndexSearcher) Search(query string, opts ...interface{}) (res search.QueryRes, err error) {
	return search.QueryRes{}, errors.New("no index search support during test")
}

func setupGinServer(gormDB *gorm.DB, mgoDB *mgo.Session, client *mongodriver.Client) *gin.Engine {
	mailSvc := mockMailStrategy{}
	searcher := mockIndexSearcher{}
	cf := controllers.NewControllerFactory(gormDB, mgoDB, mailSvc, client, searcher)
	engine := routers.SetupRouter(cf)
	return engine
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	resp := httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestMain(m *testing.M) {
	var err error
	var l net.Listener

	fmt.Println("load default config")
	if globals.Conf, err = configs.LoadDefaultConf(); err != nil {
		panic(fmt.Sprintf("Can not load default config, but got err=%+v", err))
	}

	// set up DB environment
	gormDB, mgoDB, client := setUpDBEnvironment()

	Globs.GormDB = gormDB
	Globs.MgoDB = mgoDB
	testMongoClient = client

	// set up gin server
	engine := setupGinServer(gormDB, mgoDB, client)

	Globs.GinEngine = engine

	defer Globs.GormDB.Close()
	defer Globs.MgoDB.Close()
	defer func() { testMongoClient.Disconnect(context.Background()) }()

	// start server for testing
	// the reason why we start the server
	// is because we send HTTP request internally between controllers
	ts := httptest.NewUnstartedServer(engine)
	if l, err = net.Listen("tcp", "127.0.0.1:8080"); err != nil {
		panic(err)
	}
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	retCode := m.Run()
	os.Exit(retCode)
}
