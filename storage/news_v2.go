package storage

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/internal/news"
	"github.com/twreporter/go-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"
)

type fetchResult struct {
	Content interface{}
	Error   error
}

type mongoStorage struct {
	*mongo.Client
}

type gormStorage struct {
	db *gorm.DB
}

func NewMongoV2Storage(client *mongo.Client) *mongoStorage {
	return &mongoStorage{client}
}

func NewNewsV2SqlStorage(db *gorm.DB) *gormStorage {
	return &gormStorage{db}
}

func (gs *gormStorage) GetBookmarksOfPosts(ctx context.Context, userID string, posts []news.MetaOfPost) ([]news.MetaOfPost, error) {
	if userID == "" {
		log.Error("userID is required")
		return posts, nil
	}
	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-gs.GetBookmarksOfPostsTask(ctx, userID, posts):
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

func (gs *gormStorage) GetBookmarksOfPostsTask(ctx context.Context, userID string, posts []news.MetaOfPost) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, userID string, posts []news.MetaOfPost) {
		slugs := make([]string, len(posts))
		for index, post := range posts {
			slugs[index] = post.Slug
		}

		var bookmarks []models.Bookmark
		err := gs.db.Where("id IN (?)", gs.db.Table("users_bookmarks").Select("bookmark_id").Where("user_id = ?", userID).QueryExpr()).Where("slug IN (?)", slugs).Where("deleted_at IS NULL").Find(&bookmarks).Error

		if err != nil {
			log.WithField("detail", err).Errorf("%s", f.FormatStack(err))
		}

		slugBookmarkMap := map[string]string{}
		for _, bookmark := range bookmarks {
			slugBookmarkMap[bookmark.Slug] = strconv.Itoa(int(bookmark.ID))
		}

		for index, post := range posts {
			if slugBookmarkMap[post.Slug] == "" {
				continue
			}
			posts[index].BookmarkID = slugBookmarkMap[post.Slug]
		}

		result <- fetchResult{Content: posts}
	}(ctx, userID, posts)
	return result
}

func (m *mongoStorage) GetFullPosts(ctx context.Context, q *news.Query) ([]news.Post, error) {
	var posts []news.Post

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupFullPost)...)
	// build filter stages for related documents
	stages = append(stages, news.BuildFilterRelatedPost()...)
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
	// build sorting stages from query again to sort grouped data
	stages = append(stages, news.BuildSortQueryStatements(mq)...)

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
	// build filter stages for related documents
	stages = append(stages, news.BuildFilterRelatedPost()...)

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
	// build filter stages for related documents
	stages = append(stages, news.BuildFilterRelatedPost()...)

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

func (m *mongoStorage) GetAuthors(ctx context.Context, q *news.Query) ([]news.Author, error) {
	var authors []news.Author

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupAuthor)...)
	// rewrite bio field with markdown format only
	stages = append(stages, news.BuildBioMarkdownOnlyStatement())

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getAuthors(ctx, stages):
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		authors = result.Content.([]news.Author)
	}

	return authors, nil
}

func (m *mongoStorage) getAuthors(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColContacts).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		defer cursor.Close(ctx)

		var authors []news.Author
		for cursor.Next(ctx) {
			var author news.Author
			err := cursor.Decode(&author)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			authors = append(authors, author)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: authors}
	}(ctx, stages)
	return result
}
func (m *mongoStorage) GetPostCount(ctx context.Context, q *news.Query) (int64, error) {
	return m.getCount(ctx, q, news.ColPosts)
}

func (m *mongoStorage) GetTopicCount(ctx context.Context, q *news.Query) (int64, error) {
	return m.getCount(ctx, q, news.ColTopics)
}

func (m *mongoStorage) GetAuthorCount(ctx context.Context, q *news.Query) (int64, error) {
	return m.getCount(ctx, q, news.ColContacts)
}

func (m *mongoStorage) CheckCategorySetValid(ctx context.Context, q *news.Query) (bool, error) {
	// if no subcategory then always true
	if q.Filter.CategorySet.Subcategory == "" {
		return true, nil
	}
	// if has subcateogry but no category then error
	if q.Filter.CategorySet.Category == "" {
		return false, nil
	}

	result := make(chan fetchResult)
	go func(ctx context.Context) {
		var categoryId interface{} = q.Filter.CategorySet.Category
		objectID, err := primitive.ObjectIDFromHex(q.Filter.CategorySet.Category)
		if err == nil {
			categoryId = objectID
		}

		query := bson.M{
			"_id":         categoryId,
			"subcategory": q.Filter.CategorySet.Subcategory,
		}

		count, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection("postcategories").CountDocuments(ctx, query)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: count}
	}(ctx)

	var count int64
	select {
	case <-ctx.Done():
		return false, errors.WithStack(ctx.Err())
	case res := <-result:
		if res.Error != nil {
			return false, res.Error
		}
		count = res.Content.(int64)
	}

	return (count > 0), nil
}

func (m *mongoStorage) getCount(ctx context.Context, q *news.Query, collection string) (int64, error) {
	result := make(chan fetchResult)
	go func(ctx context.Context) {
		// During mongo count document operation, empty array should be specified instead of nil(NULL).
		// Thus, rather than declare stage through var (i.e. zero value = nil)
		// use bson.D{} instead to start with empty stage([])
		document := bson.D{}
		mq := news.NewMongoQuery(q)

		// CountDocument will prepend $match key on the filter.
		// Thus, only build elements array here rather than full match stage.
		if elements := mq.GetFilter().BuildElements(); len(elements) > 0 {
			document = bson.D(elements)
		}
		count, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(collection).CountDocuments(ctx, document)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: count}
	}(ctx)
	var count int64
	select {
	case <-ctx.Done():
		return 0, errors.WithStack(ctx.Err())
	case res, ok := <-result:
		switch {
		case !ok:
			return 0, errors.WithStack(ctx.Err())
		case res.Error != nil:
			return 0, res.Error
		}
		count = res.Content.(int64)
	}

	return count, nil
}

func (m *mongoStorage) GetTags(ctx context.Context, q *news.Query) ([]news.Tag, error) {
	var tags []news.Tag

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getTags(ctx, stages):
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		tags = result.Content.([]news.Tag)
	}

	return tags, nil
}

func (m *mongoStorage) getTags(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColTags).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		defer cursor.Close(ctx)

		var tags []news.Tag
		for cursor.Next(ctx) {
			var tag news.Tag
			err := cursor.Decode(&tag)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			tags = append(tags, tag)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: tags}
	}(ctx, stages)
	return result
}

func (m *mongoStorage) GetPostReviewData(ctx context.Context, q *news.Query) ([]news.Review, error) {
	var reviews []news.Review

	mq := news.NewMongoQuery(q)

	// build aggregate stages from query
	stages := news.BuildQueryStatements(mq)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupReview)...)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-m.getReviews(ctx, stages):
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		reviews = result.Content.([]news.Review)
	}

	return reviews, nil
}

func (m *mongoStorage) getReviews(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColReviews).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		defer cursor.Close(ctx)

		var reviews []news.Review
		for cursor.Next(ctx) {
			var review news.Review
			err := cursor.Decode(&review)
			if err != nil {
				result <- fetchResult{Error: errors.WithStack(err)}
				return
			}
			reviews = append(reviews, review)
		}
		if err := cursor.Err(); err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		result <- fetchResult{Content: reviews}
	}(ctx, stages)
	return result
}
