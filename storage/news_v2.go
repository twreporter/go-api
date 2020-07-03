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
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		posts = result.Content.([]news.Post)
		for i := 0; i < len(posts); i++ {
			posts[i].Full = true
		}

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
			return
		}
		defer cursor.Close(ctx)

		var posts []news.Post
		for cursor.Next(ctx) {
			var post news.Post
			err := cursor.Decode(&post)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			posts = append(posts, post)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
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
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
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
			return
		}
		defer cursor.Close(ctx)

		var posts []news.MetaOfPost
		for cursor.Next(ctx) {
			var post news.MetaOfPost
			err := cursor.Decode(&post)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			posts = append(posts, post)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: posts}
	}(ctx, stages)
	return result
}

func (m *mongoStorage) GetFullTopics(ctx context.Context, q *news.Query) ([]news.Topic, error) {
	var topics []news.Topic

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupFullTopic)...)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getFullTopics(ctx, stages):
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		topics = result.Content.([]news.Topic)
		for i := 0; i < len(topics); i++ {
			topics[i].Full = true
		}
	}

	return topics, nil
}

func (m *mongoStorage) getFullTopics(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColTopics).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		defer cursor.Close(ctx)

		var topics []news.Topic
		for cursor.Next(ctx) {
			var topic news.Topic
			err := cursor.Decode(&topic)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			topics = append(topics, topic)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: topics}
	}(ctx, stages)
	return result
}

func (m *mongoStorage) GetMetaOfTopics(ctx context.Context, q *news.Query) ([]news.MetaOfTopic, error) {
	var topics []news.MetaOfTopic

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupMetaOfTopic)...)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getMetaOfTopics(ctx, stages):
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		topics = result.Content.([]news.MetaOfTopic)
	}

	return topics, nil
}

func (m *mongoStorage) getMetaOfTopics(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColTopics).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		defer cursor.Close(ctx)

		var topics []news.MetaOfTopic
		for cursor.Next(ctx) {
			var topic news.MetaOfTopic
			err := cursor.Decode(&topic)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			topics = append(topics, topic)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: topics}
	}(ctx, stages)
	return result
}

func (m *mongoStorage) GetPostCount(ctx context.Context, q *news.Query) (int, error) {
	// During mongo count document operation, empty array should be specified instead of nil(NULL).
	// Thus, rather than declare stage through var (i.e. zero value = nil)
	// use bson.D{} instead to start with empty stage([])
	stage := bson.D{}
	mq := news.NewMongoQuery(q)

	if stages := mq.GetFilter().BuildStage(); len(stages) > 0 {
		stage = stages[0]
	}
	count, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColPosts).CountDocuments(ctx, stage)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return int(count), nil
}

func (m *mongoStorage) GetTopicCount(ctx context.Context, q *news.Query) (int, error) {
	// During mongo count document operation, empty array should be specified instead of nil(NULL).
	// Thus, rather than declare stage through var (i.e. zero value = nil)
	// use bson.D{} instead to start with empty stage([])
	stage := bson.D{}
	mq := news.NewMongoQuery(q)

	if stages := mq.GetFilter().BuildStage(); len(stages) > 0 {
		stage = stages[0]
	}
	count, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColTopics).CountDocuments(ctx, stage)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return int(count), nil
}
