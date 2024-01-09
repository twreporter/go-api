package controllers

import (
	"fmt"
	"os"

	"github.com/twreporter/go-api/internal/news"

	"github.com/globalsign/mgo"
	"github.com/jinzhu/gorm"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/services"
	"github.com/twreporter/go-api/storage"
	"github.com/twreporter/go-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// ControllerFactory generates controlloers by given persistent storage connection
// and mail service
type ControllerFactory struct {
	gormDB      *gorm.DB
	mgoSession  *mgo.Session
	mailService services.MailService
	mongoClient *mongo.Client
	indexClient news.AlgoliaSearcher
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

// GetAnalyticsController returns *AnalyticsController struct
func (cf *ControllerFactory) GetAnalyticsController() *AnalyticsController {
	gs := storage.NewAnalyticsGormStorage(cf.gormDB)
	ms := storage.NewAnalyticsMongoStorage(cf.mongoClient)
	return NewAnalyticsController(gs, ms)
}

func (cf *ControllerFactory) GetNewsV2Controller() *newsV2Controller {
	return NewNewsV2Controller(storage.NewMongoV2Storage(cf.mongoClient), cf.indexClient, storage.NewNewsV2SqlStorage(cf.gormDB))
}

// GetMailController returns *MailController struct
func (cf *ControllerFactory) GetMailController() *MailController {
	var contrl *MailController

	contrl = NewMailController(cf.mailService, nil)

	templateDir := os.Getenv("GOAPI_HTML_TEMPLATE_DIR")

	if templateDir == "" {
		templateDir = utils.GetProjectRoot() + "/template"
	}

	contrl.LoadTemplateFiles(
		fmt.Sprintf("%s/signin.tmpl", templateDir),
		fmt.Sprintf("%s/signin-otp.tmpl", templateDir),
		fmt.Sprintf("%s/success-donation.tmpl", templateDir),
		fmt.Sprintf("%s/authenticate.tmpl", templateDir),
		fmt.Sprintf("%s/role-explorer.tmpl", templateDir),
		fmt.Sprintf("%s/role-actiontaker.tmpl", templateDir),
		fmt.Sprintf("%s/role-trailblazer.tmpl", templateDir),
		fmt.Sprintf("%s/role-downgrade.tmpl", templateDir),
	)

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
func NewControllerFactory(gormDB *gorm.DB, mgoSession *mgo.Session, mailSvc services.MailService, client *mongo.Client, sClient news.AlgoliaSearcher) *ControllerFactory {
	return &ControllerFactory{
		gormDB:      gormDB,
		mgoSession:  mgoSession,
		mailService: mailSvc,
		mongoClient: client,
		indexClient: sClient,
	}
}
