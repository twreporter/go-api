package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/controllers/oauth/facebook"
	"twreporter.org/go-api/controllers/oauth/google"
	// "twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

// Controller ...
type Controller interface {
	SetRoute(*gin.RouterGroup) *gin.RouterGroup
}

// ControllerFactory ...
type ControllerFactory struct {
	controllers map[string]Controller
	mgoDB       *mgo.Session
	gormDB      *gorm.DB
}

// GetController ...
func (cf *ControllerFactory) GetController(cn string) Controller {
	return cf.controllers[cn]
}

// SetController ...
func (cf *ControllerFactory) SetController(cn string, c Controller) {
	cf.controllers[cn] = c
}

// Close this func releases the resource appropriately
func (cf *ControllerFactory) Close() {
	cf.gormDB.Close()
	cf.mgoDB.Close()
}

func (cf *ControllerFactory) setGormDB(db *gorm.DB) {
	cf.gormDB = db
}

func (cf *ControllerFactory) setMgoDB(session *mgo.Session) {
	cf.mgoDB = session
}

// SetRoute set route by calling the correspoding controllers.
func (cf *ControllerFactory) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	for _, v := range cf.controllers {
		group = v.SetRoute(group)
	}
	return group
}

// NewControllerFactory ...
func NewControllerFactory() (*ControllerFactory, error) {
	// set up database connection
	log.Info("Connecting to MySQL cloud")
	db, err := utils.InitDB(10, 5)
	if err != nil {
		return nil, err
	}

	log.Info("Connecting to MongoDB replica")
	session, err := utils.InitMongoDB()
	if err != nil {
		return nil, err
	}
	// set up data storage
	ms := storage.NewMembershipStorage(db)

	// init controllers
	mc := NewMembershipController(ms)
	fc := facebook.Facebook{Storage: ms}
	gc := google.Google{Storage: ms}

	ns := storage.NewNewsStorage(session)
	nc := NewNewsController(ns)

	cf := &ControllerFactory{
		controllers: make(map[string]Controller),
	}
	cf.SetController(constants.MembershipController, mc)
	cf.SetController(constants.FacebookController, fc)
	cf.SetController(constants.GoogleController, gc)
	cf.SetController(constants.NewsController, nc)

	cf.setGormDB(db)

	return cf, nil
}
