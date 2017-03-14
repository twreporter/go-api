package controllers

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

type bookmarkForm struct {
	Href      string `form:"href" binding:"required"`
	Title     string `form:"title" binding:"required"`
	Desc      string `form:"desc"`
	Thumbnail string `form:"thumbnail"`
}

type bookmarkJSON struct {
	Href      string `json:"href" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Desc      string `json:"desc"`
	Thumbnail string `json:"thumbnail"`
}

func getPropsFromPOSTBody(c *gin.Context) (models.Bookmark, error) {
	var form bookmarkForm
	var json bookmarkJSON

	// Request Header
	// Content-Type: x-www-form-urlencoded
	formErr := c.Bind(&form)

	// Content-Type: application-json
	jsonErr := c.Bind(&json)

	if formErr == nil || jsonErr == nil {
		if formErr == nil {
			return models.Bookmark{Href: form.Href, Title: form.Title, Desc: utils.ToNullString(form.Desc), Thumbnail: utils.ToNullString(form.Thumbnail)}, nil
		}

		return models.Bookmark{Href: json.Href, Title: json.Title, Desc: utils.ToNullString(json.Desc), Thumbnail: utils.ToNullString(json.Thumbnail)}, nil
	}

	return models.Bookmark{}, models.NewAppError("getPropsFromPOSTBody", "controllers.account.parse_post_body", "POST body is neither JSON nor x-www-form-urlencoded", 500)
}

// BookmarkController this struct contains two stroages which have those methods to inteact with DB
type BookmarkController struct {
	BookmarkStorage storage.BookmarkStorage
	UserStorage     storage.UserStorage
}

// ListBookmarkByUser given userID this func will list all the bookmarks belongs to the user
func (bc BookmarkController) ListBookmarkByUser(c *gin.Context) {
	// get userID according to the url param
	userID := c.Param("userID")
	user, errToGetUser := bc.UserStorage.GetUserByID(userID)

	if errToGetUser != nil {
		log.Error("controllers.bookmark.list_bookmark.error_to_get_user: ", errToGetUser.Error())
		c.JSON(404, gin.H{"status": "User " + userID + " is not available.", "error": errToGetUser.Error()})
		return
	}

	log.Info("user:", user)

	bookmarks, errToGetBookmark := bc.BookmarkStorage.GetBookmarkByUser(user)

	if errToGetBookmark != nil {
		log.Error("controllers.bookmark.list_bookmark.error_to_get_bookmarks: ", errToGetBookmark.Error())
		c.JSON(404, gin.H{"status": "Bookmarks belonging to user " + userID + " is not available", "error": errToGetBookmark.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "bookmarks": bookmarks})
}

// DeleteBookmarkByUser given userID and bookmarkHref, this func will remove the relationship between user and bookmark
func (bc BookmarkController) DeleteBookmarkByUser(c *gin.Context) {
	bookmarkID := c.Param("bookmarkID")

	bookmark, errToGetBookmark := bc.BookmarkStorage.GetBookmarkByID(bookmarkID)

	if errToGetBookmark != nil {
		log.Error("controllers.bookmark.delete_bookmark.error_to_get_bookmark: ", errToGetBookmark.Error())
		c.JSON(404, gin.H{"status": "Bookmark with id " + bookmarkID + " is not available", "error": errToGetBookmark.Error()})
		return
	}

	userID := c.Param("userID")
	user, errToGetUser := bc.UserStorage.GetUserByID(userID)

	if errToGetUser != nil {
		log.Error("controllers.bookmark.delete_bookmark.error_to_get_user: ", errToGetUser.Error())
		c.JSON(404, gin.H{"status": "User " + userID + " is not available", "error": errToGetUser.Error()})
		return
	}

	errToDeleteBookmark := bc.BookmarkStorage.DeleteBookmarkByUser(user, bookmark)

	if errToDeleteBookmark != nil {
		log.Error("controllers.bookmark.delete_bookmark.error_to_delete: ", errToDeleteBookmark.Error())
		c.JSON(404, gin.H{"status": "Bookmark belongs to user " + userID + " is not available", "error": errToDeleteBookmark.Error()})
		return
	}

	c.Data(204, gin.MIMEHTML, nil)
}

// CreateBookmarkByUser given userID and bookmark POST body, this func will try to create bookmark record in the bookmarks table,
// and build the relationship between bookmark and user
func (bc BookmarkController) CreateBookmarkByUser(c *gin.Context) {
	var bookmark models.Bookmark

	userID := c.Param("userID")
	user, errToGetUser := bc.UserStorage.GetUserByID(userID)

	if errToGetUser != nil {
		log.Error("controllers.bookmark.create_bookmark.error_to_get_user: ", errToGetUser.Error())
		c.JSON(404, gin.H{"status": "User " + userID + " is not available", "error": errToGetUser.Error()})
		return
	}

	bookmark, errToParseBody := getPropsFromPOSTBody(c)
	if errToParseBody != nil {
		log.Error("controllers.bookmark.create_bookmark.error_to_parse_post_body: ", errToParseBody.Error())
		c.JSON(400, gin.H{"status": "Bad request", "error": errToParseBody.Error()})
		return
	}

	errToCreateBookmark := bc.BookmarkStorage.CreateBookmarkByUser(user, bookmark)

	if errToCreateBookmark != nil {
		log.Error("controllers.bookmark.create_bookmark.error_to_create_bookmark: ", errToCreateBookmark.Error())
		c.JSON(500, gin.H{"status": "Internal server error", "error": errToCreateBookmark.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "ok"})
}
