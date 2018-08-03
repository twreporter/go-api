package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"
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
	Controllers map[string]Controller
}

// GetController ...
func (cf *ControllerFactory) GetController(cn string) Controller {
	return cf.Controllers[cn]
}

// GetControllers returns an array of controllers
func (cf *ControllerFactory) GetControllers() []Controller {
	var cons []Controller

	for _, con := range cf.Controllers {
		cons = append(cons, con)
	}
	return cons
}

// SetController ...
func (cf *ControllerFactory) SetController(cn string, c Controller) {
	cf.Controllers[cn] = c
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
	for _, v := range cf.Controllers {
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
	fc := Facebook{Storage: gs}
	gc := Google{Storage: gs}

	ms := storage.NewMongoStorage(session)
	nc := NewNewsController(ms)

	cf := &ControllerFactory{
		Controllers: make(map[string]Controller),
	}
	cf.SetController(constants.MembershipController, mc)
	cf.SetController(constants.FacebookController, fc)
	cf.SetController(constants.GoogleController, gc)
	cf.SetController(constants.NewsController, nc)

	return cf, nil
}

func appErrorTypeAssertion(err error) *models.AppError {
	switch appErr := err.(type) {
	case *models.AppError:
		return appErr
	default:
		return models.NewAppError("AppErrorTypeAssertion", "unknown error type", err.Error(), http.StatusInternalServerError)
	}
}

type wrappedFn func(c *gin.Context) (int, gin.H, error)

func ginResponseWrapper(fn wrappedFn) func(c *gin.Context) {
	return func(c *gin.Context) {
		statusCode, obj, err := fn(c)
		if err != nil {
			appErr := err.(*models.AppError)
			log.Error(appErr.Error())
			c.JSON(appErr.StatusCode, gin.H{"status": "error", "message": appErr.Message})
			return
		}
		c.JSON(statusCode, obj)
	}
}
