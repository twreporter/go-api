package storage

import (
	"twreporter.org/go-api/models"
)

var (
	bookmarksStr = "Bookmarks"
)

// GetABookmarkByHref ...
func (g *GormMembershipStorage) GetABookmarkByHref(href string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "href = ?", href).Error
	return bookmark, err
}

// GetABookmarkByID ...
func (g *GormMembershipStorage) GetABookmarkByID(id string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "id = ?", id).Error
	return bookmark, err
}

// GetBookmarksOfAUser lists bookmarks of the user
func (g *GormMembershipStorage) GetBookmarksOfAUser(id string) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	var user models.User
	var err error

	user, err = g.GetUserByID(id)

	if err != nil {
		return []models.Bookmark{}, err
	}

	err = g.db.Model(&user).Association(bookmarksStr).Find(&bookmarks).Error
	return bookmarks, err
}

// CreateABookmarkOfAUser this func will create a bookmark and build the relationship between the bookmark and the user
func (g *GormMembershipStorage) CreateABookmarkOfAUser(userID string, bookmark models.Bookmark) error {
	var _bookmark models.Bookmark

	user, err := g.GetUserByID(userID)

	if err != nil {
		return err
	}

	// get first matched record, or create a new one
	err = g.db.Where(bookmark).FirstOrCreate(&_bookmark).Error

	if err != nil {
		return err
	}

	err = g.db.Model(&user).Association(bookmarksStr).Append(_bookmark).Error
	return err
}

// DeleteABookmarkOfAUser this func will delete the relationship between the user and the bookmark
func (g *GormMembershipStorage) DeleteABookmarkOfAUser(userID, bookmarkID string) error {
	var err error
	var bookmark models.Bookmark
	var user models.User

	user, err = g.GetUserByID(userID)
	if err != nil {
		return err
	}

	bookmark, err = g.GetABookmarkByID(bookmarkID)
	if err != nil {
		return err
	}

	err = g.db.Model(&user).Association(bookmarksStr).Delete(bookmark).Error
	return err
}
