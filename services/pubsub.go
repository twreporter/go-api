package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/twreporter/go-api/globals"
)

// PubSubService handles Google Cloud Pub/Sub operations
type PubSubService struct {
	client *pubsub.Client
	ctx    context.Context
}

// RoleUpdateMessage represents the message structure for role update notifications
type RoleUpdateMessage struct {
	Email string `json:"email"`
}

// NewPubSubService creates a new Pub/Sub service instance
func NewPubSubService() (*PubSubService, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "coastal-run-106202")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pub/sub client")
	}

	return &PubSubService{
		client: client,
		ctx:    ctx,
	}, nil
}

// getTopicName returns the appropriate topic name based on environment
func (ps *PubSubService) getTopicName() string {
	switch globals.Conf.Environment {
	case globals.DevelopmentEnvironment:
		return "dev-role-update"
	case globals.StagingEnvironment:
		return "staging-role-update"
	case globals.ProductionEnvironment:
		return "role-update"
	default:
		// Default to development for safety
		return "dev-role-update"
	}
}

// PublishRoleUpdate publishes a role update message to the appropriate topic
func (ps *PubSubService) PublishRoleUpdate(email string) error {
	topicName := ps.getTopicName()
	topic := ps.client.Topic(topicName)

	// Create the message
	message := RoleUpdateMessage{
		Email: email,
	}

	// Marshal the message to JSON
	data, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "failed to marshal role update message")
	}

	// Create the pub/sub message
	msg := &pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"source":    "go-api",
		},
	}

	// Publish the message
	result := topic.Publish(ps.ctx, msg)
	messageID, err := result.Get(ps.ctx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to publish message to topic %s", topicName))
	}

	log.WithFields(log.Fields{
		"topic":     topicName,
		"messageID": messageID,
		"email":     email,
	}).Info("Successfully published role update message")

	return nil
}

// Close closes the Pub/Sub client
func (ps *PubSubService) Close() error {
	return ps.client.Close()
}
