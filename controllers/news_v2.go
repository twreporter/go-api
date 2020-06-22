package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"twreporter.org/go-api/internal/news"
)

type newsV2Storage interface {
	GetFullPosts(context.Context, *news.Query) ([]news.Post, error)
	GetMetaOfPosts(context.Context, *news.Query) ([]news.MetaOfPost, error)
	GetFullTopics(context.Context, *news.Query) ([]news.Topic, error)
	GetMetaOfTopics(context.Context, *news.Query) ([]news.MetaOfTopic, error)

	GetPostCount(context.Context, *news.Query) (int, error)
	GetTopicCount(context.Context, *news.Filter) (int, error)
}

func NewNewsV2Controller(s newsV2Storage) *newsV2Controller {
	return &newsV2Controller{s}
}

type newsV2Controller struct {
	Storage newsV2Storage
}

func (nc *newsV2Controller) GetPosts(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				c.JSON(http.StatusGatewayTimeout, gin.H{"status": "error", "message": "Query upstream server timeout."})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unexpected error."})
			}
			log.Errorf("%+v", err)
		}
	}()

	q := news.ParsePostListQuery(c)

	// TODO(babygoat): config context with proper timeout
	posts, err := nc.Storage.GetMetaOfPosts(c, q)

	if err != nil {
		return
	}

	// TODO(babygoat): config context with proper timeout
	total, err := nc.Storage.GetPostCount(c, q)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": posts, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}}})
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
	var topic interface{}
	var err error

	defer func() {
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				c.JSON(http.StatusGatewayTimeout, gin.H{"status": "error", "message": "Query upstream server timeout."})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unexpected error."})
			}
			log.Errorf("%+v", err)
		}
	}()

	q := news.ParseSingleTopicQuery(c)

	if q.Full {
		var topics []news.Topic
		// TODO(babygoat): config context with proper timeout
		topics, err = nc.Storage.GetFullTopics(c, q)
		if len(topics) > 0 {
			topic = topics[0]
		}
	} else {
		var topics []news.MetaOfTopic
		// TODO(babygoat): config context with proper timeout
		topics, err = nc.Storage.GetMetaOfTopics(c, q)
		if len(topics) > 0 {
			topic = topics[0]
		}
	}

	// server side error
	if err != nil {
		return
	}

	if topic == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"slug": "Cannot find the topic from the slug"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": topic})
}

func (nc *newsV2Controller) GetIndexPage(c *gin.Context) {
}
