package controllers

import (
	"fmt"
	"go/build"
	"net/http"
	"os"
	"path"

	// "gopkg.in/mgo.v2/bson"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/services"
	"twreporter.org/go-api/storage"
	//log "github.com/Sirupsen/logrus"
)

// ControllerFactory generates controlloers by given persistent storage connection
// and mail service
type ControllerFactory struct {
	gormDB      *gorm.DB
	mgoSession  *mgo.Session
	mailService services.MailService
}

// GetGoogleController returns Google struct
func (cf *ControllerFactory) GetGoogleController() Google {
	gs := storage.NewGormStorage(cf.gormDB)
	return Google{Storage: gs}
}

// GetFacebookController returns Facebook struct
func (cf *ControllerFactory) GetFacebookController() Facebook {
	gs := storage.NewGormStorage(cf.gormDB)
	return Facebook{Storage: gs}
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

// GetMailController returns *MailController struct
func (cf *ControllerFactory) GetMailController() *MailController {
	var gopath string
	var filepath string
	var contrl *MailController

	contrl = NewMailController(cf.mailService, nil)

	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	filepath = path.Join(gopath, "src/twreporter.org/go-api/template")

	contrl.LoadTemplateFiles(fmt.Sprintf("%s/signin.tmpl", filepath), fmt.Sprintf("%s/success-donation.tmpl", filepath))

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
func NewControllerFactory(gormDB *gorm.DB, mgoSession *mgo.Session, mailSvc services.MailService) *ControllerFactory {
	return &ControllerFactory{
		gormDB:      gormDB,
		mgoSession:  mgoSession,
		mailService: mailSvc,
	}
}

func appErrorTypeAssertion(err error) *models.AppError {
	switch appErr := err.(type) {
	case *models.AppError:
		return appErr
	default:
		return models.NewAppError("AppErrorTypeAssertion", "unknown error type", err.Error(), http.StatusInternalServerError)
	}
}
