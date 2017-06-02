package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// GetPosts receive HTTP GET method request, and return the posts.
// `query`, `limit` and `offset` are the url query params,
// which define the rule we retrieve posts from storage.
func (nc *NewsController) GetPosts(c *gin.Context) {
	var metaOfPosts []models.PostMeta
	var err error

	qs := c.Query("where")
	limit := c.Query("limit")
	offset := c.Query("offset")
	full := c.Query("full")

	_limit, _ := strconv.Atoi(limit)
	_offset, _ := strconv.Atoi(offset)
	_full, _ := strconv.ParseBool(full)

	if _limit == 0 {
		_limit = 10
	}

	if _full {
	} else {
		metaOfPosts, err = nc.Storage.GetMetaOfPosts(qs, _limit, _offset, []string{"hero_image", "categories", "tags", "topic", "og_image"})
		if err != nil {
			appErr := err.(models.AppError)
			c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "records": metaOfPosts})
	}
}
