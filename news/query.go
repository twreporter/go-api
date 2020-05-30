package news

import (
	"net/url"
	"strconv"
)

type (
	Pagination struct {
		Offset int
		Limit  int
	}

	Filter struct {
		Slug  string
		State string
	}

	Query struct {
		Pagination
		Filter Filter
		Sort   string
		Full   bool
	}

	Options func(*Query)
)

func FromSlug(slug string) Options {
	return func(q *Query) {
		q.Filter.Slug = slug
	}
}

func FromUrlQueryMap(u url.Values) Options {
	return func(q *Query) {
		offset, err := strconv.Atoi(u.Get("offset"))
		if err != nil {
			q.Offset = offset
		}
		limit, err := strconv.Atoi(u.Get("limitt"))
		if err != nil {
			q.Limit = limit
		}

		full, err := strconv.ParseBool(u.Get("full"))
		if err != nil {
			q.Full = full
		}
	}
}

func NewQuery(options ...Options) *Query {
	q := &Query{
		Pagination: Pagination{Offset: 0, Limit: 10},
		Sort:       "-publishedDate",
	}

	for _, o := range options {
		o(q)
	}
	return q
}

func (q *Query) SetPagination(offset, limit int) {
	q.Offset = offset
	q.Limit = limit
}

func (q *Query) SetSort(sort string) {
	q.Sort = sort
}
