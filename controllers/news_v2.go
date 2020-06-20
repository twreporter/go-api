package controllers

import (
	"context"

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
}

func (nc *newsV2Controller) GetTopics(c *gin.Context) {
}

func (nc *newsV2Controller) GetATopic(c *gin.Context) {
}

func (nc *newsV2Controller) GetIndexPage(c *gin.Context) {
}
