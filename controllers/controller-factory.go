package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/controllers/oauth/facebook"
	"twreporter.org/go-api/controllers/oauth/google"
	"twreporter.org/go-api/storage"
)

// Controller ...
type Controller interface {
	SetRoute(*gin.RouterGroup) *gin.RouterGroup
}

// ControllerFactory ...
type ControllerFactory struct {
	controllers map[string]Controller
}

// GetController ...
func (cf *ControllerFactory) GetController(cn string) Controller {
	return cf.controllers[cn]
}

// SetController ...
func (cf *ControllerFactory) SetController(cn string, c Controller) {
	cf.controllers[cn] = c
}

// SetRoute set route by calling the correspoding controllers.
func (cf *ControllerFactory) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	for _, v := range cf.controllers {
		group = v.SetRoute(group)
	}
	return group
}

// NewControllerFactory ...
func NewControllerFactory(db *gorm.DB) *ControllerFactory {
	// set up data storage
	s := storage.NewMembershipStorage(db)

	// init controllers
	mc := NewMembershipController(s)
	fc := facebook.Facebook{s}
	gc := google.Google{s}

	cf := &ControllerFactory{
		controllers: make(map[string]Controller),
	}
	cf.SetController(constants.MembershipController, mc)
	cf.SetController(constants.FacebookController, fc)
	cf.SetController(constants.GoogleController, gc)

	return cf
}
