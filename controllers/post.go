package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/models"
)

// GetPosts receive HTTP GET method request, and return the posts.
// `query`, `limit`, `offset`, `sort` and `full` are the url query params,
// which define the rule we retrieve posts from storage.
func (nc *NewsController) GetPosts(c *gin.Context) (int, gin.H, error) {
	var total int
	var posts []models.Post = make([]models.Post, 0)

	err, mq, limit, offset, sort, full := nc.GetQueryParam(c)

	// response empty records if parsing url query param occurs error
	if err != nil {
		return http.StatusOK, gin.H{"status": "ok", "records": posts, "meta": models.MetaOfResponse{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		}}, nil
	}

	if limit == 0 {
		limit = 10
	}

	if sort == "" {
		sort = "-publishedDate"
	}

	if full {
		posts, total, err = nc.Storage.GetFullPosts(mq, limit, offset, sort, nil)
	} else {
		posts, total, err = nc.Storage.GetMetaOfPosts(mq, limit, offset, sort, nil)
	}

	if err != nil {
		return toPostResponse(err)
	}

	// make sure `response.records`
	// would be `[]` rather than  `null`
	if posts == nil {
		posts = make([]models.Post, 0)
	}

	return http.StatusOK, gin.H{"status": "ok", "records": posts, "meta": models.MetaOfResponse{
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}}, nil
}

// GetAPost receive HTTP GET method request, and return the certain post.
func (nc *NewsController) GetAPost(c *gin.Context) (int, gin.H, error) {
	var posts []models.Post
	var err error

	slug := c.Param("slug")
	full, _ := strconv.ParseBool(c.Query("full"))

	mq := models.MongoQuery{
		Slug: slug,
	}

	if full {
		posts, _, err = nc.Storage.GetFullPosts(mq, 1, 0, "-publishedDate", nil)
	} else {
		posts, _, err = nc.Storage.GetMetaOfPosts(mq, 1, 0, "-publishedDate", nil)
	}

	if err != nil {
		return toPostResponse(err)
	}

	if len(posts) == 0 {
		return http.StatusNotFound, gin.H{"status": "Record Not Found", "error": "Record Not Found"}, nil
	}

	return http.StatusOK, gin.H{"status": "ok", "record": posts[0]}, nil
}
