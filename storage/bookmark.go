package storage

import (
	"twreporter.org/go-api/models"
	// log "github.com/Sirupsen/logrus"
)

var (
	bookmarksStr = "Bookmarks"
)

// GetABookmarkBySlug ...
func (g *GormStorage) GetABookmarkBySlug(slug string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "slug = ?", slug).Error
	if err != nil {
		return bookmark, g.NewStorageError(err, "GetABookmarkBySlug", "storage.bookmark.error_to_get")
	}

	return bookmark, err
}

// GetABookmarkByID ...
func (g *GormStorage) GetABookmarkByID(id string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "id = ?", id).Error
	if err != nil {
		return bookmark, g.NewStorageError(err, "GetABookmarkByID", "storage.bookmark.error_to_get")
	}

	return bookmark, err
}

// GetABookmarkOfAUser get a bookmark of a user
func (g *GormStorage) GetABookmarkOfAUser(userID string, bookmarkSlug string, bookmarkHost string) (models.Bookmark, error) {
	var bookmarks []models.Bookmark
	var bookmark models.Bookmark
	var user models.User
	var err error

	if user, err = g.GetUserByID(userID); err != nil {
		return bookmark, err
	}

	err = g.db.Model(&user).Association(bookmarksStr).Find(&bookmarks).Error

	if err != nil {
		return bookmark, g.NewStorageError(err, "GetABookmarkOfAUser", "storage.bookmark.error_to_get_a_bookmark_of_a_user")
	}

	for _, ele := range bookmarks {
		if ele.Slug == bookmarkSlug && ele.Host == bookmarkHost {
			return ele, err
		}
	}

	return bookmark, models.NewAppError("GetABookmarkOfAUser", "Record not found", "storage.bookmark.bookmark_is_not_registered_by_user", 404)
}

// GetBookmarksOfAUser lists bookmarks of the user
func (g *GormStorage) GetBookmarksOfAUser(id string, limit, offset int) ([]models.Bookmark, int, error) {
	var bookmarks []models.Bookmark
	var user models.User
	var err error

	if user, err = g.GetUserByID(id); err != nil {
		return bookmarks, 0, err
	}

	// The reason I write the raw sql statement, not use gorm association(see the following commented code),
	// err = g.db.Model(&user).Limit(limit).Offset(offset).Order("created_at desc").Related(&bookmarks, bookmarksStr).Error
	// is because I need to sort/limit/offset the records occording to `users_bookmarks`.`created_at`.
	err = g.db.Raw("SELECT `users_bookmarks`.created_at AS users_bookmarks_created_at, `bookmarks`.* FROM `bookmarks` INNER JOIN `users_bookmarks` ON `users_bookmarks`.`bookmark_id` = `bookmarks`.`id` WHERE `bookmarks`.deleted_at IS NULL AND ((`users_bookmarks`.`user_id` IN (?))) ORDER BY users_bookmarks_created_at desc LIMIT ? OFFSET ?", id, limit, offset).Scan(&bookmarks).Error

	if err != nil {
		return bookmarks, 0, g.NewStorageError(err, "GetBookmarksOfAUser", "storage.bookmark.error_to_get_bookmarks")
	}

	total := g.db.Model(&user).Association(bookmarksStr).Count()

	return bookmarks, total, err
}

// CreateABookmarkOfAUser this func will create a bookmark and build the relationship between the bookmark and the user
func (g *GormStorage) CreateABookmarkOfAUser(userID string, bookmark models.Bookmark) (models.Bookmark, error) {
	var _bookmark = bookmark
	var err error
	var user models.User

	if user, err = g.GetUserByID(userID); err != nil {
		return _bookmark, err
	}

	// get first matched record, or create a new one
	err = g.db.Where("slug = ? AND host = ?", bookmark.Slug, bookmark.Host).FirstOrCreate(&_bookmark).Error

	if err != nil {
		return _bookmark, g.NewStorageError(err, "CreateABookmarkOfAUser", "storage.bookmark.error_to_create_bookmark")
	}

	// update the bookmark fields
	err = g.db.Model(&_bookmark).Updates(bookmark).Error

	if err != nil {
		return _bookmark, g.NewStorageError(err, "CreateABookmarkOfAUser", "storage.bookmark.error_to_update_bookmark")
	}

	err = g.db.Model(&user).Association(bookmarksStr).Append(_bookmark).Error
	if err != nil {
		return _bookmark, g.NewStorageError(err, "CreateABookmarkOfAUser", "storage.bookmark.error_to_create_user_bookmark_relationship")
	}

	return _bookmark, err
}

// DeleteABookmarkOfAUser this func will delete the relationship between the user and the bookmark
func (g *GormStorage) DeleteABookmarkOfAUser(userID, bookmarkID string) error {
	var err error
	var bookmark models.Bookmark
	var user models.User

	if user, err = g.GetUserByID(userID); err != nil {
		return err
	}

	bookmark, err = g.GetABookmarkByID(bookmarkID)
	if err != nil {
		return g.NewStorageError(err, "DeleteABookmarkOfAUser", "storage.bookmark.error_to_get_bookmark")
	}

	// The reason why here find before delete is to make sure it will return error if record is not found
	err = g.db.Model(&user).Association(bookmarksStr).Find(&bookmark).Delete(bookmark).Error
	if err != nil {
		return g.NewStorageError(err, "DeleteABookmarkOfAUser", "storage.bookmark.error_to_delete_user_bookmark_relationship")
	}

	return err
}
