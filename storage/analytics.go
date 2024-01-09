package storage

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/twreporter/go-api/models"
)

type AnalyticsGormStorage interface {
	UpdateUserReadingPostCount(string, string) (bool, error)
	UpdateUserReadingPostTime(string, string, int) (error)
}

type AnalyticsMongoStorage interface {

}

type mongoDB struct {
	*mongo.Client
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
