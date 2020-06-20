package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, opts *options.ClientOptions) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "Establishing a new connection to cluster occurs error:")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "Connection to cluster does not response:")
	}
	return client, nil
}
