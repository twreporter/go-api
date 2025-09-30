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

// NewPubSubService creates a new Pub/Sub service instance
func NewPubSubService() (*PubSubService, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, globals.Conf.PubSub.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pub/sub client")
	}

	return &PubSubService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Publish publishes a message to the specified topic
func (ps *PubSubService) Publish(topicName string, data interface{}, attributes map[string]string) error {
	topic := ps.client.Topic(topicName)

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message data")
	}

	// Set default attributes if none provided
	if attributes == nil {
		attributes = make(map[string]string)
	}

	// Add default timestamp and source if not already set
	if _, exists := attributes["timestamp"]; !exists {
		attributes["timestamp"] = time.Now().Format(time.RFC3339)
	}
	if _, exists := attributes["source"]; !exists {
		attributes["source"] = "go-api"
	}

	// Create the pub/sub message
	msg := &pubsub.Message{
		Data:       jsonData,
		Attributes: attributes,
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
	}).Info("Successfully published message")

	return nil
}

// Close closes the Pub/Sub client
func (ps *PubSubService) Close() error {
	return ps.client.Close()
}
