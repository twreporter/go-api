package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"
)

// GetPosts receive HTTP GET method request, and return the posts.
// `query`, `limit`, `offset`, `sort` and `full` are the url query params,
// which define the rule we retrieve posts from storage.
func (nc *NewsController) GetPosts(c *gin.Context) {
	var total int
	var posts []models.Post
	var err error

	qs, limit, offset, sort, full := nc.GetQueryParam(c)

	if qs == "" {
		qs = "{}"
	}

	if limit == 0 {
		limit = 10
	}

	if sort == "" {
		sort = "-publishedDate"
	}

	if full {
		posts, total, err = nc.Storage.GetFullPosts(qs, limit, offset, sort, nil)
	} else {
		posts, total, err = nc.Storage.GetMetaOfPosts(qs, limit, offset, sort, nil)
	}

	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": posts, "meta": models.MetaOfResponse{
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}})
}

// GetAPost receive HTTP GET method request, and return the certain post.
func (nc *NewsController) GetAPost(c *gin.Context) {
	var posts []models.Post
	var err error

	slug := c.Param("slug")
	full, _ := strconv.ParseBool(c.Query("full"))

	qs := bson.M{
		"slug": slug,
	}

	if full {
		posts, _, err = nc.Storage.GetFullPosts(qs, 1, 0, "-publishedDate", nil)
	} else {
		posts, _, err = nc.Storage.GetMetaOfPosts(qs, 1, 0, "-publishedDate", nil)
	}

	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Record Not Found", "error": "Record Not Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "record": posts[0]})
}
