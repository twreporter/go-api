package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/internal/news"
)

type newsV2Storage interface {
	GetFullPosts(context.Context, *news.Query) ([]news.Post, error)
	GetMetaOfPosts(context.Context, *news.Query) ([]news.MetaOfPost, error)
	GetFullTopics(context.Context, *news.Query) ([]news.Topic, error)
	GetMetaOfTopics(context.Context, *news.Query) ([]news.MetaOfTopic, error)

	GetPostCount(context.Context, *news.Filter) (int, error)
	GetTopicCount(context.Context, *news.Filter) (int, error)
}

func NewNewsV2Controller(s newsV2Storage) *newsV2Controller {
	return &newsV2Controller{s}
}

type newsV2Controller struct {
	Storage newsV2Storage
}

func (nc *newsV2Controller) GetPosts(c *gin.Context) {
}

func (nc *newsV2Controller) GetAPost(c *gin.Context) {
	var post interface{}
	var err error

	q := news.ParseSinglePostQuery(c)

	if q.Full {
		var posts []news.Post
		// TODO(babygoat): config context with proper timeout
		posts, err = nc.Storage.GetFullPosts(c, q)
		if len(posts) > 0 {
			post = posts[0]
		}
	} else {
		var posts []news.MetaOfPost
		// TODO(babygoat): config context with proper timeout
		posts, err = nc.Storage.GetMetaOfPosts(c, q)
		if len(posts) > 0 {
			post = posts[0]
		}
	}

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		c.JSON(http.StatusGatewayTimeout, gin.H{"status": "error", "message": "Query upstream server timeout."})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unexpected error."})
		return
	case post == nil:
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"slug": "Cannot find the post from the slug"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

func (nc *newsV2Controller) GetTopics(c *gin.Context) {
}

func (nc *newsV2Controller) GetATopic(c *gin.Context) {
}

func (nc *newsV2Controller) GetIndexPage(c *gin.Context) {
}
