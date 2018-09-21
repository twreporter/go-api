package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
)

// NewsController has methods to handle requests which wants posts, topics ... etc news resource.
type NewsController struct {
	Storage storage.NewsStorage
}

// NewNewsController ...
func NewNewsController(s storage.NewsStorage) *NewsController {
	return &NewsController{s}
}

// GetQueryParam pares url param
func (nc *NewsController) GetQueryParam(c *gin.Context) (mq models.MongoQuery, limit int, offset int, sort string, full bool) {
	where := c.Query("where")
	_limit := c.Query("limit")
	_offset := c.Query("offset")
	_full := c.Query("full")
	sort = c.Query("sort")

	// provide default param if error occurs
	limit, _ = strconv.Atoi(_limit)
	offset, _ = strconv.Atoi(_offset)
	full, _ = strconv.ParseBool(_full)

	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	if where == "" {
		where = "{}"
	}

	_ = models.GetQuery(where, &mq)

	return
}
