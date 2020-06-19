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

const (
	sortByPublishedDate = "published_date"
	sortByUpdatedAt     = "updated_at"
	descending          = "-"
)

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

func ParsePostListQuery(c *gin.Context) *Query {
	var q Query

	defaultQuery := Query{
		Pagination: query.Pagination{Offset: 0, Limit: 10},
		Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
	}

	q = defaultQuery
	// Parse filter
	if len(c.QueryArray("category_id")) > 0 {
		q.Filter.Categories = c.QueryArray("category_id")
	}
	if len(c.QueryArray("tag_id")) > 0 {
		q.Filter.Tags = c.QueryArray("tag_id")
	}
	if len(c.QueryArray("id")) > 0 {
		q.Filter.Tags = c.QueryArray("id")
	}

	// Parse pagination
	if offset, err := strconv.Atoi(c.Query("offset")); err == nil {
		q.Offset = offset
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		q.Limit = limit
	}

	// Parse sorting
	if sort := c.Query("sort"); sort != "" {
		switch sort {
		case sortByPublishedDate:
			q.Sort = SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(true)}}
		case descending + sortByPublishedDate:
			q.Sort = SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}}
		case sortByUpdatedAt:
			q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(true)}}
		case descending + sortByUpdatedAt:
			q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(false)}}
		}
	}
	return &q
}
