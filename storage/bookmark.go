package storage

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"twreporter.org/go-api/models"
)

var (
	bookmarksStr = "Bookmarks"
)

// GetABookmarkBySlug ...
func (g *GormStorage) GetABookmarkBySlug(slug string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "slug = ?", slug).Error
	if err != nil {
		return bookmark, errors.Wrap(err, fmt.Sprintf("get bookmark(slug: '%s') occurs error", slug))
	}

	return bookmark, nil
}

// GetABookmarkByID ...
func (g *GormStorage) GetABookmarkByID(id string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "id = ?", id).Error
	if err != nil {
		return bookmark, errors.Wrap(err, fmt.Sprintf("get bookmark(id: '%s') occurs error", id))
	}

	return bookmark, nil
}

// GetABookmarkOfAUser get a bookmark of a user
func (g *GormStorage) GetABookmarkOfAUser(userID string, slug string, host string) (models.Bookmark, error) {
	var bookmark models.Bookmark

	err := g.db.Preload("Users", "id = ?", userID).Where("slug = ? and host = ?", slug, host).First(&bookmark).Error

	if err != nil {
		return bookmark, errors.Wrap(err, fmt.Sprintf("get bookmark(slug: '%s', host: '%s') from user(id: '%s') occurs error", slug, host, userID))
	}

	return bookmark, nil
}

// GetBookmarksOfAUser lists bookmarks of the user
func (g *GormStorage) GetBookmarksOfAUser(id string, limit, offset int) ([]models.Bookmark, int, error) {
	var bookmarks []models.Bookmark

	// The reason I write the raw sql statement, not use gorm association(see the following commented code),
	// err = g.db.Model(&user).Limit(limit).Offset(offset).Order("created_at desc").Related(&bookmarks, bookmarksStr).Error
	// is because I need to sort/limit/offset the records occording to `users_bookmarks`.`created_at`.
	err := g.db.Raw("SELECT `users_bookmarks`.created_at AS users_bookmarks_created_at, `bookmarks`.* FROM `bookmarks` INNER JOIN `users_bookmarks` ON `users_bookmarks`.`bookmark_id` = `bookmarks`.`id` WHERE `bookmarks`.deleted_at IS NULL AND ((`users_bookmarks`.`user_id` IN (?))) ORDER BY users_bookmarks_created_at desc LIMIT ? OFFSET ?", id, limit, offset).Scan(&bookmarks).Error

	if err != nil {
		return bookmarks, 0, errors.Wrap(err, fmt.Sprintf("get bookmarks of the user(id: %s) with conditions(limit: %d, offset: %d) occurs error", id, limit, offset))
	}

	userID, _ := strconv.Atoi(id)
	total := g.db.Model(models.User{ID: uint(userID)}).Association(bookmarksStr).Count()

	return bookmarks, total, nil
}

// CreateABookmarkOfAUser this func will create a bookmark and build the relationship between the bookmark and the user
func (g *GormStorage) CreateABookmarkOfAUser(userID string, bookmark models.Bookmark) (models.Bookmark, error) {
	var _bookmark = bookmark

	user, err := g.GetUserByID(userID)
	if err != nil {
		return _bookmark, errors.WithStack(err)
	}

	// get first matched record, or create a new one
	err = g.db.Where("slug = ? AND host = ?", bookmark.Slug, bookmark.Host).FirstOrCreate(&_bookmark).Error

	if err != nil {
		return _bookmark, errors.Wrap(err, fmt.Sprintf("create a bookmark(%#v) occurs error", bookmark))
	}

	err = g.db.Model(&user).Association(bookmarksStr).Append(_bookmark).Error
	if err != nil {
		return _bookmark, errors.Wrap(err, fmt.Sprintf("append the bookmark(%#v) to the user(id: %s) occurs error", bookmark, userID))
	}

	return _bookmark, nil
}

// DeleteABookmarkOfAUser this func will delete the relationship between the user and the bookmark
func (g *GormStorage) DeleteABookmarkOfAUser(userID, bookmarkID string) error {
	user, err := g.GetUserByID(userID)
	if err != nil {
		return errors.WithStack(err)
	}

	bookmark, err := g.GetABookmarkByID(bookmarkID)
	if err != nil {
		return errors.WithStack(err)
	}

	// The reason why here find before delete is to make sure it will return error if record is not found
	err = g.db.Model(&user).Association(bookmarksStr).Find(&bookmark).Delete(bookmark).Error
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("delete bookmark(id: %s) from user(id: %s) occurs error", bookmarkID, userID))
	}

	return nil
}
