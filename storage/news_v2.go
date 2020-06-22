package storage

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/internal/news"
)

type fetchResult struct {
	Content interface{}
	Error   error
}

type mongoStorage struct {
	*mongo.Client
}

func NewMongoV2Storage(client *mongo.Client) *mongoStorage {
	return &mongoStorage{client}
}

func (m *mongoStorage) GetFullPosts(ctx context.Context, q *news.Query) ([]news.Post, error) {
	var posts []news.Post

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupFullPost)...)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getFullPosts(ctx, stages):
		if !ok {
			return nil, errors.WithStack(ctx.Err())
		}
		posts = result.Content.([]news.Post)
	}

	return posts, nil
}

func (m *mongoStorage) getFullPosts(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColPosts).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
		}
		defer cursor.Close(ctx)

		var posts []news.Post
		for cursor.Next(ctx) {
			var post news.Post
			err := cursor.Decode(&post)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
			}
			posts = append(posts, post)
		}
		result <- fetchResult{Content: posts}
	}(ctx, stages)
	return result
}

func (m *mongoStorage) GetMetaOfPosts(ctx context.Context, q *news.Query) ([]news.MetaOfPost, error) {
	var posts []news.MetaOfPost

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupMetaOfPost)...)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getMetaOfPosts(ctx, stages):
		if !ok {
			return nil, errors.WithStack(ctx.Err())
		}
		posts = result.Content.([]news.MetaOfPost)
	}

	return posts, nil
}

func (m *mongoStorage) getMetaOfPosts(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColPosts).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
		}
		defer cursor.Close(ctx)

		var posts []news.MetaOfPost
		for cursor.Next(ctx) {
			var post news.MetaOfPost
			err := cursor.Decode(&post)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
			}
			posts = append(posts, post)
		}
		result <- fetchResult{Content: posts}
	}(ctx, stages)
	return result
}

func (m *mongoStorage) GetFullTopics(context.Context, *news.Query) ([]news.Topic, error) {
	return nil, nil
}
func (m *mongoStorage) GetMetaOfTopics(context.Context, *news.Query) ([]news.MetaOfTopic, error) {
	return nil, nil
}
func (m *mongoStorage) GetPostCount(context.Context, *news.Filter) (int, error)  { return 0, nil }
func (m *mongoStorage) GetTopicCount(context.Context, *news.Filter) (int, error) { return 0, nil }
