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
	sortByDescending    = "-"

	queryFull       = "full"
	querySlug       = "slug"
	queryCategoryID = "category_id"
	queryTagID      = "tag_id"
	queryPostID     = "id"
	querySort       = "sort"
	queryOffset     = "offset"
	queryLimit      = "limit"
)

func ParseSinglePostQuery(c *gin.Context) *Query {
	var q Query

	if slug := c.Param(querySlug); slug != "" {
		q.Filter.Slug = slug
	}

	if full, err := strconv.ParseBool(c.Query(queryFull)); err == nil {
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
	if len(c.QueryArray(queryCategoryID)) > 0 {
		q.Filter.Categories = c.QueryArray(queryCategoryID)
	}
	if len(c.QueryArray(queryTagID)) > 0 {
		q.Filter.Tags = c.QueryArray(queryTagID)
	}
	if len(c.QueryArray(queryPostID)) > 0 {
		q.Filter.Tags = c.QueryArray(queryPostID)
	}

	// Parse pagination
	if offset, err := strconv.Atoi(c.Query(queryOffset)); err == nil {
		q.Offset = offset
	}
	if limit, err := strconv.Atoi(c.Query(queryLimit)); err == nil {
		q.Limit = limit
	}

	// Parse sorting
	if sort := c.Query(querySort); sort != "" {
		switch sort {
		case sortByPublishedDate:
			q.Sort = SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(true)}}
		case sortByDescending + sortByPublishedDate:
			q.Sort = SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}}
		case sortByUpdatedAt:
			q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(true)}}
		case sortByDescending + sortByUpdatedAt:
			q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(false)}}
		}
	}
	return &q
}
