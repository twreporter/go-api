package storage

import (
	"fmt"
	"strconv"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/internal/news"
)

type AnalyticsGormStorage interface {
	UpdateUserReadingPostCount(string, string) (bool, error)
	UpdateUserReadingPostTime(string, string, int) (error)
	UpdateUserReadingFootprint(string, string) (bool, error)
	GetFootprintsOfAUser(string, int, int) ([]models.UsersPostsReadingFootprint, int, error)
}

type AnalyticsMongoStorage interface {
	GetPostsOfIDs(context.Context, []string) ([]news.MetaOfFootprint, error)
}

type mongoDB struct {
	db *mongo.Client
}

type gormDB struct {
	db *gorm.DB
}

func NewAnalyticsGormStorage(db *gorm.DB) *gormDB {
	return &gormDB{db}
}

func NewAnalyticsMongoStorage(db *mongo.Client) *mongoDB {
	return &mongoDB{db}
}

// Update user reading posts count
func (gs *gormDB) UpdateUserReadingPostCount(userID string, postID string) (bool, error) {
	tx := gs.db.Begin() // Start the transaction

	// Check if the user exists
	var count int64
	if err := tx.Model(&models.UsersPostsReadingCount{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to check user read posts count existence (user_id: %s, post_id: %s)", userID, postID))
	}
	// Directly return if record exist
	if count != 0 {
		tx.Rollback()
		return true, nil
	}

	// Add the user read post record
	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to parse int from user_id: %s", userID))
	}
	if err := tx.Create(&models.UsersPostsReadingCount{ UserID: userIdInt, PostID: postID }).Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to create user read post count (user_id: %s, post_id: %s)", userID, postID))
	}

	// Update the user read post count on users table
	if err := tx.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("read_posts_count", gorm.Expr("read_posts_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to update user's read post count (user_id: %s)", userID))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, "failed to commit transaction")
	}

	return false, nil
}

// Update user reading posts count
func (gs *gormDB) UpdateUserReadingPostTime(userID string, postID string, second int) (error) {
	tx := gs.db.Begin() // Start the transaction

	// Add the user read post record
	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("failed to parse int from user_id: %s", userID))
	}
	if err := tx.Create(&models.UsersPostsReadingTime{ UserID: userIdInt, PostID: postID, Seconds: second }).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("failed to create user read post time (user_id: %s, post_id: %s, seconds: %d)", userID, postID, second))
	}

	// Update the user read post count on users table
	if err := tx.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("read_posts_sec", gorm.Expr("read_posts_sec + ?", second)).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("failed to update user's read post time (user_id: %s)", userID))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Update user reading footprint: post
func (gs *gormDB) UpdateUserReadingFootprint(userID string, postID string) (bool, error) {
	tx := gs.db.Begin() // Start the transaction

	// Check if the user exists
	var count int64
	if err := tx.Model(&models.UsersPostsReadingFootprint{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to check user reading footprint existence (user_id: %s, post_id: %s)", userID, postID))
	}
	// Directly return if record exist
	if count != 0 {
		tx.Rollback()
		return true, nil
	}

	// Add the user read post record
	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to parse int from user_id: %s", userID))
	}
	if err := tx.Create(&models.UsersPostsReadingFootprint{ UserID: userIdInt, PostID: postID }).Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, fmt.Sprintf("failed to create user reading footprint (user_id: %s, post_id: %s)", userID, postID))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, errors.Wrap(err, "failed to commit transaction")
	}

	return false, nil
}

func (gs *gormDB) GetFootprintsOfAUser(userID string, limit int, offset int) ([]models.UsersPostsReadingFootprint, int, error) {
	var err error
	var total int
	var footprints []models.UsersPostsReadingFootprint

	statement := gs.db.Model(&models.UsersPostsReadingFootprint{}).Where("user_id = ?", userID)
	if err = statement.Limit(limit).Offset(offset).Order("updated_at desc").Find(&footprints).Error; err != nil {
		return nil, 0, err
	}
	if err = statement.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return footprints, total, nil
}

func (ms *mongoDB) GetPostsOfIDs(ctx context.Context, postIDs []string) ([]news.MetaOfFootprint, error) {
	var posts []news.MetaOfFootprint
	if len(postIDs) == 0 {
		return posts, nil
	}

	// build _id filter
	stages := news.BuildFilterIDs(postIDs)
	// build lookup(join) stages according to required fields
	stages = append(stages, news.BuildLookupStatements(news.LookupMetaOfFootprint)...)

	select {
	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	case result, ok := <-ms.getMetaOfFootprints(ctx, stages):
		switch {
		case !ok:
			return nil, errors.WithStack(ctx.Err())
		case result.Error != nil:
			return nil, result.Error
		}
		posts = result.Content.([]news.MetaOfFootprint)
	}

	return posts, nil
}

func (ms *mongoDB) getMetaOfFootprints(ctx context.Context, stages []bson.D) <-chan fetchResult {
	result := make(chan fetchResult)
	go func(ctx context.Context, stages []bson.D) {
		defer close(result)
		cursor, err := ms.db.Database(globals.Conf.DB.Mongo.DBname).Collection(news.ColPosts).Aggregate(ctx, stages)
		if err != nil {
			result <- fetchResult{Error: errors.WithStack(err)}
			return
		}
		defer cursor.Close(ctx)

		var posts []news.MetaOfFootprint
		for cursor.Next(ctx) {
			var post news.MetaOfFootprint
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
