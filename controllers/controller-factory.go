package controllers

import (
	"github.com/gin-gonic/gin"
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
	Close() error
}

// ControllerFactory ...
type ControllerFactory struct {
	controllers map[string]Controller
}

// GetController ...
func (cf *ControllerFactory) GetController(cn string) Controller {
	return cf.controllers[cn]
}

// GetControllers returns an array of controllers
func (cf *ControllerFactory) GetControllers() []Controller {
	var cons []Controller

	for _, con := range cf.controllers {
		cons = append(cons, con)
	}
	return cons
}

// SetController ...
func (cf *ControllerFactory) SetController(cn string, c Controller) {
	cf.controllers[cn] = c
}

// Close this func releases the resource appropriately
func (cf *ControllerFactory) Close() error {
	var err error
	for _, controller := range cf.GetControllers() {
		err = controller.Close()
		if err != nil {
			return err
		}
	}
	return nil
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
	gs := storage.NewGormStorage(db)

	// init controllers
	mc := NewMembershipController(gs)
	fc := facebook.Facebook{Storage: gs}
	gc := google.Google{Storage: gs}

	ms := storage.NewMongoStorage(session)
	nc := NewNewsController(ms)

	cf := &ControllerFactory{
		controllers: make(map[string]Controller),
	}
	cf.SetController(constants.MembershipController, mc)
	cf.SetController(constants.FacebookController, fc)
	cf.SetController(constants.GoogleController, gc)
	cf.SetController(constants.NewsController, nc)

	return cf, nil
}
