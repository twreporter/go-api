package news

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type newsStorage interface {
	GetPosts(*Query) ([]Post, error)
	GetTopics(*Query) ([]Topic, error)

	GetPostCount(*Filter) (int, error)
	GetTopicCount(*Filter) (int, error)
}

type newsController struct {
	Storage newsStorage
}

func NewController(s newsStorage) *newsController {
	return &newsController{s}
}

func (nc *newsController) GetPosts(c *gin.Context) {
	q := NewQuery(FromUrlQueryMap(c.Request.URL.Query()))

	posts, err := nc.Storage.GetPosts(q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	total, err := nc.Storage.GetPostCount(&q.Filter)
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
	q := NewQuery(FromSlug(c.Query("slug")))
	q.Limit = 1

	posts, err := nc.Storage.GetPosts(q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": posts[0]})
}

func (nc *newsController) GetTopics(c *gin.Context) {
	q := NewQuery(FromUrlQueryMap(c.Request.URL.Query()))

	topics, err := nc.Storage.GetTopics(q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	total, err := nc.Storage.GetTopicCount(&q.Filter)
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
	q := NewQuery(FromSlug(c.Query("slug")))
	q.Limit = 1

	topics, err := nc.Storage.GetTopics(q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": topics[0]})
}
