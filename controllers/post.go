package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// GetQueryParam ...
func (nc *NewsController) GetQueryParam(c *gin.Context) (qs string, limit int, offset int, sort string, full bool) {
	qs = c.Query("where")
	_limit := c.Query("limit")
	_offset := c.Query("offset")
	_full := c.Query("full")
	sort = c.Query("sort")

	limit, _ = strconv.Atoi(_limit)
	offset, _ = strconv.Atoi(_offset)
	full, _ = strconv.ParseBool(_full)
	return
}

// GetPosts receive HTTP GET method request, and return the posts.
// `query`, `limit`, `offset` and `sort` are the url query params,
// which define the rule we retrieve posts from storage.
func (nc *NewsController) GetPosts(c *gin.Context) {
	var metaOfPosts []models.PostMeta
	var err error

	qs, limit, offset, sort, full := nc.GetQueryParam(c)

	if limit == 0 {
		limit = 10
	}

	if sort == "" {
		sort = "-publishedDate"
	}

	if full {
	} else {
		metaOfPosts, err = nc.Storage.GetMetaOfPosts(qs, limit, offset, sort, []string{"hero_image", "categories", "tags", "topic", "og_image"})
		if err != nil {
			appErr := err.(models.AppError)
			c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "records": metaOfPosts})
	}
}
