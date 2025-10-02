package controllers

import (
	log "github.com/sirupsen/logrus"
	"github.com/twreporter/go-api/services"
	"github.com/twreporter/go-api/storage"
)

// NewMembershipController ...
func NewMembershipController(s storage.MembershipStorage) *MembershipController {
	pubSubService, err := services.NewPubSubService()
	if err != nil {
		// Log error but don't fail the controller creation
		// The service will handle pub/sub failures gracefully
		log.WithField("error", err).Error("Failed to initialize PubSubService, role update messages will not be sent")
	}

	var roleUpdateService *services.RoleUpdateService
	if pubSubService != nil {
		roleUpdateService = services.NewRoleUpdateService(pubSubService)
	}

	return &MembershipController{
		Storage:           s,
		PubSubService:     pubSubService,
		RoleUpdateService: roleUpdateService,
	}
}

// MembershipController ...
type MembershipController struct {
	Storage           storage.MembershipStorage
	PubSubService     *services.PubSubService
	RoleUpdateService *services.RoleUpdateService
}

// Close is the method of Controller interface
func (mc *MembershipController) Close() error {
	err := mc.Storage.Close()
	if err != nil {
		return err
	}

	// Close pub/sub service if it exists
	if mc.PubSubService != nil {
		if closeErr := mc.PubSubService.Close(); closeErr != nil {
			return closeErr
		}
	}

	return nil
}
