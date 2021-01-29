package news

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/internal/query"
	"gopkg.in/guregu/null.v3"
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
	Name       string
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
	queryKeywords   = "keywords"
)

type Option func(*Query)

var defaultQuery = Query{
	Pagination: query.Pagination{Offset: 0, Limit: 10},
	Filter:     Filter{State: "published"},
	Sort:       SortBy{PublishedDate: query.Order{IsAsc: null.BoolFrom(false)}},
}

var defaultAuthorQuery = Query{
	Pagination: query.Pagination{Offset: 0, Limit: 10},
	Sort:       SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(false)}},
}

// NewQuery returns a default query along with the options(pagination/sort/filter).
// Note that if the same options are specified multiple times, the last one will be applied.
func NewQuery(options ...Option) *Query {
	q := defaultQuery
	for _, o := range options {
		o(&q)
	}
	return &q
}

func WithOffset(offset int) Option {
	return func(q *Query) {
		q.Offset = offset
	}
}

func WithLimit(limit int) Option {
	return func(q *Query) {
		q.Limit = limit
	}
}

// FilterCategoryIDs adds category ids into filter on the query
func WithFilterCategoryIDs(ids ...string) Option {
	return func(q *Query) {
		if len(ids) > 0 {
			q.Filter.Categories = ids
		}
	}
}

// FilterState adds the post publish state filter on the query
func WithFilterState(state string) Option {
	return func(q *Query) {
		q.Filter.State = state
	}
}

// FilterStyle adds the post style filter on the query
func WithFilterStyle(style string) Option {
	return func(q *Query) {
		q.Filter.Style = style
	}
}

// FilterIsFeatured adds the isFeatured filter on the query
func WithFilterIsFeatured(isFeatured bool) Option {
	return func(q *Query) {
		q.Filter.IsFeatured = null.BoolFrom(isFeatured)
	}
}

// SortUpdatedAt updates the query to sort by updatedAt field
func WithSortUpdatedAt(isAsc bool) Option {
	return func(q *Query) {
		q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(isAsc)}}
	}
}

func ParseSinglePostQuery(c *gin.Context) *Query {
	return parseSingleQuery(c)
}

func ParsePostListQuery(c *gin.Context) *Query {
	var q Query

	q = defaultQuery
	// Parse filter
	if len(c.QueryArray(queryCategoryID)) > 0 {
		q.Filter.Categories = c.QueryArray(queryCategoryID)
	}
	if len(c.QueryArray(queryTagID)) > 0 {
		q.Filter.Tags = c.QueryArray(queryTagID)
	}
	if len(c.QueryArray(queryPostID)) > 0 {
		q.Filter.IDs = c.QueryArray(queryPostID)
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

func ParseSingleTopicQuery(c *gin.Context) *Query {
	return parseSingleQuery(c)
}

func ParseTopicListQuery(c *gin.Context) *Query {
	var q Query

	q = defaultQuery

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
		}
	}
	return &q
}

func parseSingleQuery(c *gin.Context) *Query {
	var q Query

	if slug := c.Param(querySlug); slug != "" {
		q.Filter.Slug = slug
	}

	if full, err := strconv.ParseBool(c.Query(queryFull)); err == nil {
		q.Full = full
	}
	return &q
}

func ParseAuthorListQuery(c *gin.Context) *Query {
	var q Query

	q = defaultAuthorQuery
	if keywords := c.Query(queryKeywords); keywords != "" {
		q.Filter.Name = keywords
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
		case sortByUpdatedAt:
			q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(true)}}
		case sortByDescending + sortByUpdatedAt:
			q.Sort = SortBy{UpdatedAt: query.Order{IsAsc: null.BoolFrom(false)}}
		}
	}
	return &q
}
