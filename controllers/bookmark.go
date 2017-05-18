package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

// BookmarkController this struct contains two stroages which have those methods to inteact with DB
type BookmarkController struct {
	Storage storage.MembershipStorage
}

// SetRoute is the method of Controller interface
func (bc BookmarkController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	// handle bookmarks of users
	group.GET("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), bc.GetBookmarksOfAUser)
	group.POST("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), bc.CreateABookmarkOfAUser)
	group.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.CheckJWT(), middlewares.ValidateUserID(), bc.DeleteABookmarkOfAUser)
	return group
}

// GetBookmarksOfAUser given userID this func will list all the bookmarks belongs to the user
func (bc BookmarkController) GetBookmarksOfAUser(c *gin.Context) {
	var err error
	var bookmarks []models.Bookmark

	// get userID according to the url param
	userID := c.Param("userID")
	bookmarks, err = bc.Storage.GetBookmarksOfAUser(userID)

	if err != nil && err.Error() == utils.ErrRecordNotFound.Error() {
		c.JSON(http.StatusNotFound, gin.H{"status": "User not found", "error": err.Error()})
		return
	} else if err != nil {
		log.Error("controllers.bookmark.get_bookmarks_of_a_user.error_to_get_bookmarks: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": bookmarks})
}

// DeleteABookmarkOfAUser given userID and bookmarkHref, this func will remove the relationship between user and bookmark
func (bc BookmarkController) DeleteABookmarkOfAUser(c *gin.Context) {
	bookmarkID := c.Param("bookmarkID")
	userID := c.Param("userID")

	err := bc.Storage.DeleteABookmarkOfAUser(userID, bookmarkID)

	if err != nil && err.Error() == utils.ErrRecordNotFound.Error() {
		c.JSON(http.StatusNotFound, gin.H{"status": "User not found", "error": err.Error()})
		return
	} else if err != nil {
		log.Error("controllers.bookmark.delete_a_bookmark_of_a_user.error_to_delete_bookmark: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// CreateABookmarkOfAUser given userID and bookmark POST body, this func will try to create bookmark record in the bookmarks table,
// and build the relationship between bookmark and user
func (bc BookmarkController) CreateABookmarkOfAUser(c *gin.Context) {
	var bookmark models.Bookmark
	var err error

	userID := c.Param("userID")
	bookmark, err = bc.parseBody(c)
	if err != nil {
		log.Error("controllers.bookmark.create_bookmark.error_to_parse_post_body: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad request", "error": err.Error()})
		return
	}

	err = bc.Storage.CreateABookmarkOfAUser(userID, bookmark)

	if err != nil && err.Error() == utils.ErrRecordNotFound.Error() {
		c.JSON(http.StatusNotFound, gin.H{"status": "User not found", "error": err.Error()})
		return
	} else if err != nil {
		log.Error("controllers.bookmark.create_bookmark_of_a_user.error_to_create_bookmark: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (bc BookmarkController) parseBody(c *gin.Context) (models.Bookmark, error) {
	var err error
	var form models.BookmarkForm
	var json models.BookmarkJSON

	contentType := c.ContentType()

	if contentType == "application/json" {
		err = c.Bind(&json)
		if err != nil {
			return models.Bookmark{}, err
		}
		return models.Bookmark{Href: json.Href, Title: json.Title, Desc: utils.ToNullString(json.Desc), Thumbnail: utils.ToNullString(json.Thumbnail)}, nil
	} else if contentType == "x-www-form-urlencoded" {
		err = c.Bind(&form)
		if err != nil {
			return models.Bookmark{}, err
		}
		return models.Bookmark{Href: form.Href, Title: form.Title, Desc: utils.ToNullString(form.Desc), Thumbnail: utils.ToNullString(form.Thumbnail)}, nil
	}

	return models.Bookmark{}, models.NewAppError("parseBody", "controllers.account.parse_post_body", "POST body is neither JSON nor x-www-form-urlencoded", http.StatusBadRequest)
}
