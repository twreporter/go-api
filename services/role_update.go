package services

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
)

// RoleUpdateService handles role update related operations
type RoleUpdateService struct {
	pubSubService *PubSubService
	topicName     string
}

// NewRoleUpdateService creates a new role update service instance
func NewRoleUpdateService(pubSubService *PubSubService) *RoleUpdateService {
	return &RoleUpdateService{
		pubSubService: pubSubService,
		topicName:     globals.Conf.PubSub.TopicName,
	}
}

// PublishRoleUpdate publishes a role update message to the appropriate topic
func (rus *RoleUpdateService) PublishRoleUpdate(email string) error {
	if rus.pubSubService == nil {
		return errors.New("pub/sub service is not available")
	}

	// Create the message
	message := models.RoleUpdateMessage{
		Email: email,
	}

	// Publish the message using the generic pub/sub service
	err := rus.pubSubService.Publish(rus.topicName, message, nil)
	if err != nil {
		return errors.Wrap(err, "failed to publish role update message")
	}

	log.WithFields(log.Fields{
		"topic": rus.topicName,
		"email": email,
	}).Info("Successfully published role update message")

	return nil
}
