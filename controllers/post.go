package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// GetPosts receive HTTP GET method request, and return the posts.
// `query`, `limit`, `offset` and `sort` are the url query params,
// which define the rule we retrieve posts from storage.
func (nc *NewsController) GetPosts(c *gin.Context) {
	var metaOfPosts []models.Post
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
	} else {
		metaOfPosts, err = nc.Storage.GetMetaOfPosts(qs, limit, offset, sort, nil)
		if err != nil {
			appErr := err.(models.AppError)
			c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "records": metaOfPosts})
	}
}
