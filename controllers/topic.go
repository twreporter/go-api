package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// GetTopics receive HTTP GET method request, and return the topics.
// `query`, `limit`, `offset` and `sort` are the url query params,
// which define the rule we retrieve topics from storage.
func (nc *NewsController) GetTopics(c *gin.Context) {
	var topics []models.Topic
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
		topics, err = nc.Storage.GetTopics(qs, limit, offset, sort, nil)
	} else {
		topics, err = nc.Storage.GetTopics(qs, limit, offset, sort, []string{"leading_image", "og_image"})
	}

	if err != nil {
		appErr := err.(models.AppError)
		c.JSON(appErr.StatusCode, gin.H{"status": appErr.Message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": topics})
	return

}
