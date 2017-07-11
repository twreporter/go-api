package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	// "twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/storage"
)

// NewsController has methods to handle requests which wants posts, topics ... etc news resource.
type NewsController struct {
	Storage storage.NewsStorage
}

// NewNewsController ...
func NewNewsController(s storage.NewsStorage) Controller {
	return &NewsController{s}
}

// Close is the method of Controller interface
func (nc *NewsController) Close() error {
	err := nc.Storage.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetQueryParam pares url param
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

// SetRoute is the method of Controller interface
func (nc *NewsController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	// endpoints for posts
	group.GET("/posts", nc.GetPosts)
	group.GET("/posts/:slug", nc.GetAPost)

	// endpoints for topics
	group.GET("/topics", nc.GetTopics)
	group.GET("/topics/:slug", nc.GetATopic)

	group.GET("/index_page", nc.GetIndexPageContents)
	group.GET("/index_page_categories", nc.GetCategoriesPosts)

	// endpoints for search
	group.GET("/search/authors", nc.SearchAuthors)
	group.GET("/search/posts", nc.SearchPosts)
	return group
}
