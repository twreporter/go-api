package tests

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/controllers/oauth/facebook"
	"twreporter.org/go-api/controllers/oauth/google"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

var (
	DefaultID        = "1"
	DefaultAccount   = "nickhsine@twreporter.org"
	DefaultPassword  = "0000"
	DefaultID2       = "2"
	DefaultAccount2  = "turtle@twreporter.org"
	DefaultPassword2 = "1111"
	DefaultService   = "default_service"
	DefaultToken     = "default_token"
	Engine           *gin.Engine
	DB               *gorm.DB
)

func OpenTestConnection() (db *gorm.DB, err error) {
	// CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm';
	// CREATE DATABASE gorm;
	// GRANT ALL ON gorm.* TO 'gorm'@'localhost';
	fmt.Println("testing mysql...")
	dbhost := os.Getenv("GORM_DBADDRESS")
	if dbhost != "" {
		dbhost = fmt.Sprintf("tcp(%v)", dbhost)
	} else {
		dbhost = "tcp(127.0.0.1:3306)"
	}
	db, err = gorm.Open("mysql", fmt.Sprintf("gorm:gorm@%v/gorm?charset=utf8&parseTime=True", dbhost))

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(10)

	return
}

func RunMigration() {
	for _, table := range []string{"users_bookmarks"} {
		DB.Exec(fmt.Sprintf("drop table %v;", table))
	}

	values := []interface{}{&models.User{}, &models.OAuthAccount{}, &models.ReporterAccount{}, &models.Bookmark{}, &models.Registration{}, &models.Service{}}
	for _, value := range values {
		DB.DropTable(value)
	}
	if err := DB.AutoMigrate(values...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}
}

func SetDefaultRecords() {
	// Set an active reporter account
	ms := storage.NewGormStorage(DB)

	key, _ := scrypt.Key([]byte(DefaultPassword), []byte(""), 16384, 8, 1, 32)

	ra := models.ReporterAccount{
		Account:       DefaultAccount,
		Password:      fmt.Sprintf("%x", key),
		Active:        true,
		ActivateToken: "",
	}
	_, _ = ms.InsertUserByReporterAccount(ra)

	key, _ = scrypt.Key([]byte(DefaultPassword2), []byte(""), 16384, 8, 1, 32)

	ra = models.ReporterAccount{
		Account:       DefaultAccount2,
		Password:      fmt.Sprintf("%x", key),
		Active:        true,
		ActivateToken: "",
	}
	_, _ = ms.InsertUserByReporterAccount(ra)

	ms.CreateService(models.ServiceJSON{Name: DefaultService})

	ms.CreateRegistration(DefaultService, models.RegistrationJSON{Email: DefaultAccount, ActivateToken: DefaultToken})
}

func SetupGinServer() {
	// set up data storage
	gs := storage.NewGormStorage(DB)

	// init controllers
	mc := controllers.NewMembershipController(gs)
	fc := facebook.Facebook{Storage: gs}
	gc := google.Google{Storage: gs}

	cf := &controllers.ControllerFactory{
		Controllers: make(map[string]controllers.Controller),
	}
	cf.SetController(constants.MembershipController, mc)
	cf.SetController(constants.FacebookController, fc)
	cf.SetController(constants.GoogleController, gc)

	Engine = gin.Default()
	routerGroup := Engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)
}

func RequestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

func GenerateJWT(user models.User) (jwt string) {
	jwt, _ = utils.RetrieveToken(user.ID, user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
	return
}

func GetUser(userId string) (user models.User) {
	as := storage.NewGormStorage(DB)
	user, _ = as.GetUserByID(userId)
	return
}

func ServeHTTP(method, path, body, contentType, authorization string) (resp *httptest.ResponseRecorder) {
	var req *http.Request

	req = RequestWithBody(method, path, body)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)

	return
}
