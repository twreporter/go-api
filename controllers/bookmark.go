package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"twreporter.org/go-api/models"
)

// GetBookmarksOfAUser given userID this func will list all the bookmarks belongs to the user
func (mc *MembershipController) GetBookmarksOfAUser(c *gin.Context) (int, gin.H, error) {
	var err error
	var bookmarks []models.Bookmark
	var bookmark models.Bookmark
	var total int

	// get userID according to the url param
	userID := c.Param("userID")

	// get bookmarkSlug in url param
	bookmarkSlug := c.Param("bookmarkSlug")

	// Get a specific bookmark from a user
	if bookmarkSlug != "" {
		host := c.Query("host")

		if bookmark, err = mc.Storage.GetABookmarkOfAUser(userID, bookmarkSlug, host); err != nil {
			return toResponse(err)
		}

		return http.StatusOK, gin.H{"status": "ok", "record": bookmark}, nil
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if limit == 0 {
		limit = 10
	}

	if bookmarks, total, err = mc.Storage.GetBookmarksOfAUser(userID, limit, offset); err != nil {
		return toResponse(err)
	}

	// TODO The response JSON should be like
	//	{
	//		"status": "success",
	//		"data":  {
	//			"meta": meta,
	//			"records": bookmarks
	//		}
	//	}
	return http.StatusOK, gin.H{"status": "ok", "records": bookmarks, "meta": models.MetaOfResponse{
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}}, nil
}

// DeleteABookmarkOfAUser given userID and bookmarkHref, this func will remove the relationship between user and bookmark
func (mc *MembershipController) DeleteABookmarkOfAUser(c *gin.Context) (int, gin.H, error) {
	bookmarkID := c.Param("bookmarkID")
	userID := c.Param("userID")

	if err := mc.Storage.DeleteABookmarkOfAUser(userID, bookmarkID); err != nil {
		return toResponse(err)
	}

	return http.StatusNoContent, gin.H{}, nil
}

// CreateABookmarkOfAUser given userID and bookmark POST body, this func will try to create bookmark record in the bookmarks table,
// and build the relationship between bookmark and user
func (mc *MembershipController) CreateABookmarkOfAUser(c *gin.Context) (int, gin.H, error) {
	var bookmark models.Bookmark
	var err error

	userID := c.Param("userID")
	if bookmark, err = mc.parseBookmarkPOSTBody(c); err != nil {
		// For legacy code, the response returns with status "error"
		// TODO rewrite with status "fail"
		return http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()}, nil
	}

	if bookmark, err = mc.Storage.CreateABookmarkOfAUser(userID, bookmark); err != nil {
		return toResponse(err)
	}

	// TODO The response JSON should be like
	//	{
	//		"status": "success",
	//		"data": bookmark
	//	}
	return http.StatusCreated, gin.H{"status": "ok", "record": bookmark}, nil
}

func (mc *MembershipController) parseBookmarkPOSTBody(c *gin.Context) (models.Bookmark, error) {
	var bm models.Bookmark

	if err := c.Bind(&bm); err != nil {
		return models.Bookmark{}, errors.Wrap(err, "POST body is neither JSON nor x-www-form-urlencoded")
	}
	return bm, nil
}
