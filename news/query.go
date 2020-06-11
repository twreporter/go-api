package news

import (
	"net/url"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/guregu/null.v3"
)

type (
	Pagination struct {
		Offset int
		Limit  int
	}

	Filter struct {
		Slug       string
		State      string
		Style      string
		IsFeatured null.Bool
		Categories []primitive.ObjectID
	}

	Sort struct {
		PublishedDate SortOrder
		UpdatedAt     SortOrder
	}

	SortOrder struct {
		IsAsc null.Bool
	}

	Query struct {
		Pagination
		Filter Filter
		Sort   Sort
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
		if err == nil {
			q.Offset = offset
		}
		limit, err := strconv.Atoi(u.Get("limitt"))
		if err == nil {
			q.Limit = limit
		}

		full, err := strconv.ParseBool(u.Get("full"))
		if err == nil {
			q.Full = full
		}
	}
}

func NewQuery(options ...Options) *Query {
	q := &Query{
		Pagination: Pagination{Offset: 0, Limit: 10},
		Sort:       Sort{PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)}},
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

func (q *Query) SetSort(sort Sort) {
	q.Sort = sort
}
