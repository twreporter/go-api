package controllers

import (
	"net/http"

	// "gopkg.in/mgo.v2/bson"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
	//log "github.com/Sirupsen/logrus"
)

// ControllerFactory ...
type ControllerFactory struct {
	gormDB     *gorm.DB
	mgoSession *mgo.Session
	mailSender *utils.EmailContext
}

func (cf *ControllerFactory) GetGoogleController() Google {
	gs := storage.NewGormStorage(cf.gormDB)
	return Google{Storage: gs}
}

func (cf *ControllerFactory) GetFacebookController() Facebook {
	gs := storage.NewGormStorage(cf.gormDB)
	return Facebook{Storage: gs}
}

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

func (cf *ControllerFactory) GetMembershipController() *MembershipController {
	gs := storage.NewGormStorage(cf.gormDB)
	return NewMembershipController(gs)
}

func (cf *ControllerFactory) GetNewsController() *NewsController {
	ms := storage.NewMongoStorage(cf.mgoSession)
	return NewNewsController(ms)
}

func (cf *ControllerFactory) GetMailSender() *utils.EmailContext {
	return cf.mailSender
}

func (cf *ControllerFactory) GetMgoSession() *mgo.Session {
	return cf.mgoSession
}

func (cf *ControllerFactory) GetGormDB() *gorm.DB {
	return cf.gormDB
}

// NewControllerFactory ...
func NewControllerFactory(gormDB *gorm.DB, mgoSession *mgo.Session, mailSender *utils.EmailContext) *ControllerFactory {
	return &ControllerFactory{
		gormDB:     gormDB,
		mgoSession: mgoSession,
		mailSender: mailSender,
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
