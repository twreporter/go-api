package news

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/internal/query"
)

type Query struct {
	query.Pagination
	Filter Filter
	Sort   SortBy
	Full   bool
}

type Filter struct {
	Slug       string
	State      string
	Style      string
	IsFeatured null.Bool
	Categories []string
	Tags       []string
	IDs        []string
}

type SortBy struct {
	PublishedDate query.Order
	UpdatedAt     query.Order
}

func ParseSinglePostQuery(c *gin.Context) *Query {
	var q Query

	if slug := c.Param("slug"); slug != "" {
		q.Filter.Slug = slug
	}

	if full, err := strconv.ParseBool(c.Query("full")); err == nil {
		q.Full = full
	}
	return &q
}
