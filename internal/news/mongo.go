package news

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/internal/query"
)

type mongoQuery struct {
	mongoPagination
	mongoFilter
	mongoSort
}

func NewMongoQuery(q *Query) *mongoQuery {
	return &mongoQuery{
		fromPagination(q.Pagination),
		fromFilter(q.Filter),
		fromSort(q.Sort),
	}
}

type mongoPagination struct {
	Skip  int
	Limit int
}

func fromPagination(p query.Pagination) mongoPagination {
	return mongoPagination{
		Skip:  p.Offset,
		Limit: p.Limit,
	}
}

type mongoFilter struct {
	Slug       string               `mongo:"slug"`
	State      string               `mongo:"state"`
	Style      string               `mongo:"style"`
	IsFeatured null.Bool            `mongo:"isFeatured"`
	Categories []primitive.ObjectID `mongo:"categories"`
	Tags       []primitive.ObjectID `mongo:"tags"`
	IDs        []primitive.ObjectID `mongo:"_id"`
}

func fromFilter(f Filter) mongoFilter {
	return mongoFilter{
		Slug:       f.Slug,
		State:      f.State,
		Style:      f.Style,
		IsFeatured: f.IsFeatured,
		Categories: hexToObjectIDs(f.Categories),
		Tags:       hexToObjectIDs(f.Tags),
		IDs:        hexToObjectIDs(f.IDs),
	}
}

func hexToObjectIDs(hs []string) []primitive.ObjectID {
	var ids []primitive.ObjectID

	for _, v := range hs {
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			// ignore invalid objectID
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

type mongoSort struct {
	PublishedDate query.Order `mongo:"publishedDate"`
	UpdatedAt     query.Order `mongo:"updatedAt"`
}

func fromSort(s SortBy) mongoSort {
	return mongoSort{
		PublishedDate: s.PublishedDate,
		UpdatedAt:     s.UpdatedAt,
	}
}
