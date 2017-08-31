package controllers

import (
	"net/http"

	// log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// GetBookmarksOfAUser given userID this func will list all the bookmarks belongs to the user
func (mc *MembershipController) GetBookmarksOfAUser(c *gin.Context) {
	var err error
	var appErr models.AppError
	var bookmarks []models.Bookmark

	// get userID according to the url param
	userID := c.Param("userID")
	bookmarks, err = mc.Storage.GetBookmarksOfAUser(userID)

	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": bookmarks})
}

// DeleteABookmarkOfAUser given userID and bookmarkHref, this func will remove the relationship between user and bookmark
func (mc *MembershipController) DeleteABookmarkOfAUser(c *gin.Context) {
	var appErr models.AppError

	bookmarkID := c.Param("bookmarkID")
	userID := c.Param("userID")

	err := mc.Storage.DeleteABookmarkOfAUser(userID, bookmarkID)

	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.Data(http.StatusNoContent, gin.MIMEHTML, nil)
}

// CreateABookmarkOfAUser given userID and bookmark POST body, this func will try to create bookmark record in the bookmarks table,
// and build the relationship between bookmark and user
func (mc *MembershipController) CreateABookmarkOfAUser(c *gin.Context) {
	var appErr models.AppError
	var bookmark models.Bookmark
	var err error

	userID := c.Param("userID")
	bookmark, err = mc.parseBookmarkPOSTBody(c)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	err = mc.Storage.CreateABookmarkOfAUser(userID, bookmark)
	if err != nil {
		appErr = err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (mc *MembershipController) parseBookmarkPOSTBody(c *gin.Context) (models.Bookmark, error) {
	var err error
	var bm models.Bookmark

	if err = c.Bind(&bm); err != nil {
		return models.Bookmark{}, models.NewAppError("parseBookmarkPOSTBody", "Bad request", "POST body is neither JSON nor x-www-form-urlencoded", http.StatusBadRequest)
	}

	return bm, nil
}
