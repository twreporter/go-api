package news

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type newsStorage interface {
	GetPosts(context.Context, *Query) ([]Post, error)
	GetTopics(context.Context, *Query) ([]Topic, error)

	GetPostCount(context.Context, *Filter) (int, error)
	GetTopicCount(context.Context, *Filter) (int, error)
}

type newsController struct {
	Storage newsStorage
}

func NewController(s newsStorage) *newsController {
	return &newsController{s}
}

func (nc *newsController) GetPosts(c *gin.Context) {
	q := NewQuery(FromUrlQueryMap(c.Request.URL.Query()))

	posts, err := nc.Storage.GetPosts(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	total, err := nc.Storage.GetPostCount(c, &q.Filter)
	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": posts, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}})
}

func (nc *newsController) GetAPost(c *gin.Context) {
	q := NewQuery(FromSlug(c.Param("slug")), FromUrlQueryMap(c.Request.URL.Query()))
	q.Limit = 1

	posts, err := nc.Storage.GetPosts(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": posts[0]})
}

func (nc *newsController) GetTopics(c *gin.Context) {
	q := NewQuery(FromUrlQueryMap(c.Request.URL.Query()))

	topics, err := nc.Storage.GetTopics(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	total, err := nc.Storage.GetTopicCount(c, &q.Filter)
	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": topics, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}})
}

func (nc *newsController) GetATopic(c *gin.Context) {
	q := NewQuery(FromSlug(c.Param("slug")), FromUrlQueryMap(c.Request.URL.Query()))
	q.Limit = 1

	topics, err := nc.Storage.GetTopics(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": topics[0]})
}
