package controllers

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/services"
	"github.com/twreporter/go-api/storage"
	"github.com/twreporter/go-api/utils"
)

// ControllerFactory generates controlloers by given persistent storage connection
// and mail service
type ControllerFactory struct {
	gormDB      *gorm.DB
	mgoSession  *mgo.Session
	mailService services.MailService
	mongoClient *mongo.Client
}

// GetOAuthController returns OAuth struct
func (cf *ControllerFactory) GetOAuthController(oauthType string) (oauth *OAuth) {
	gs := storage.NewGormStorage(cf.gormDB)
	oauth = &OAuth{Storage: gs}
	if oauthType == globals.GoogleOAuth {
		oauth.InitGoogleConfig()
	} else {
		oauth.InitFacebookConfig()
	}

	return oauth
}

// GetMembershipController returns *MembershipController struct
func (cf *ControllerFactory) GetMembershipController() *MembershipController {
	gs := storage.NewGormStorage(cf.gormDB)
	return NewMembershipController(gs)
}

// GetNewsController returns *NewsController struct
func (cf *ControllerFactory) GetNewsController() *NewsController {
	ms := storage.NewMongoStorage(cf.mgoSession)
	return NewNewsController(ms)
}

func (cf *ControllerFactory) GetNewsV2Controller() *newsV2Controller {
	return NewNewsV2Controller(storage.NewMongoV2Storage(cf.mongoClient))
}

// GetMailController returns *MailController struct
func (cf *ControllerFactory) GetMailController() *MailController {
	var contrl *MailController

	contrl = NewMailController(cf.mailService, nil)

	templateDir := os.Getenv("GOAPI_HTML_TEMPLATE_DIR")

	if templateDir == "" {
		templateDir = utils.GetProjectRoot() + "/template"
	}

	contrl.LoadTemplateFiles(fmt.Sprintf("%s/signin.tmpl", templateDir), fmt.Sprintf("%s/success-donation.tmpl", templateDir))

	return contrl
}

// GetMailService returns MailService it holds
func (cf *ControllerFactory) GetMailService() services.MailService {
	return cf.mailService
}

// GetMgoSession returns *mgo.Session it holds
func (cf *ControllerFactory) GetMgoSession() *mgo.Session {
	return cf.mgoSession
}

// GetMgoSession returns *gorm.DB it holds
func (cf *ControllerFactory) GetGormDB() *gorm.DB {
	return cf.gormDB
}

// NewControllerFactory generate *ControllerFactory struct
func NewControllerFactory(gormDB *gorm.DB, mgoSession *mgo.Session, mailSvc services.MailService, client *mongo.Client) *ControllerFactory {
	return &ControllerFactory{
		gormDB:      gormDB,
		mgoSession:  mgoSession,
		mailService: mailSvc,
		mongoClient: client,
	}
}
