package storage

import (
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
)

var (
	bookmarksStr = "Bookmarks"
)

// BookmarkStorage this is an interface defines methods for users and bookmarks tables
type BookmarkStorage interface {
	// get
	GetBookmarkByHref(string) (models.Bookmark, error)
	GetBookmarkByID(string) (models.Bookmark, error)
	GetBookmarkByUser(models.User) ([]models.Bookmark, error)

	// create
	CreateBookmarkByUser(models.User, models.Bookmark) error

	// delete
	DeleteBookmarkByUser(models.User, models.Bookmark) error
}

// NewGormBookmarkStorage this initializes the user storage
func NewGormBookmarkStorage(db *gorm.DB) BookmarkStorage {
	return &gormBookmarkStorage{db}
}

// gormBookmarkStorage this implements UserStorage interface
type gormBookmarkStorage struct {
	db *gorm.DB
}

func (g *gormBookmarkStorage) GetBookmarkByHref(href string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "href = ?", href).Error
	return bookmark, err
}

func (g *gormBookmarkStorage) GetBookmarkByID(id string) (models.Bookmark, error) {
	var bookmark models.Bookmark
	err := g.db.First(&bookmark, "id = ?", id).Error
	return bookmark, err
}

// GetBookmarkByUser this func will list bookmarks of the user
func (g *gormBookmarkStorage) GetBookmarkByUser(user models.User) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	err := g.db.Model(&user).Association(bookmarksStr).Find(&bookmarks).Error
	return bookmarks, err
}

// CreateBookmarkByUser this func will create a bookmark and build the relationship between the bookmark and the user
func (g *gormBookmarkStorage) CreateBookmarkByUser(user models.User, bookmark models.Bookmark) error {
	var err error

	// try to create a bookmark record
	err = g.db.Create(&bookmark).Error

	// bookmark record is already existed (href is unique key in table)
	if err != nil {
		bookmark, err = g.GetBookmarkByHref(bookmark.Href)

		if err != nil {
			return err
		}
	}

	err = g.db.Model(&user).Association(bookmarksStr).Append(bookmark).Error
	return err
}

// DeleteBookmarkByUser this func will delete the relationship between the user and the bookmark
func (g *gormBookmarkStorage) DeleteBookmarkByUser(user models.User, bookmark models.Bookmark) error {
	err := g.db.Model(&user).Association(bookmarksStr).Delete(bookmark).Error
	return err
}
