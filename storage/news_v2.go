package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"twreporter.org/go-api/internal/news"
)

type mongoStorage struct {
	*mongo.Client
}

func NewMongoV2Storage(client *mongo.Client) *mongoStorage {
	return &mongoStorage{client}
}

func (m *mongoStorage) GetFullPosts(context.Context, *news.Query) ([]news.Post, error) {
	return nil, nil
}
func (m *mongoStorage) GetMetaOfPosts(context.Context, *news.Query) ([]news.MetaOfPost, error) {
	return nil, nil
}
func (m *mongoStorage) GetFullTopics(context.Context, *news.Query) ([]news.Topic, error) {
	return nil, nil
}
func (m *mongoStorage) GetMetaOfTopics(context.Context, *news.Query) ([]news.MetaOfTopic, error) {
	return nil, nil
}
func (m *mongoStorage) GetPostCount(context.Context, *news.Filter) (int, error)  { return 0, nil }
func (m *mongoStorage) GetTopicCount(context.Context, *news.Filter) (int, error) { return 0, nil }
