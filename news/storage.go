package news

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"twreporter.org/go-api/globals"
)

var (
	fullPostLookupFields         = []string{"heroImage", "leading_image_portrait", "leading_video", "categories", "tags", "topic_full", "og_image", "writters", "photographers", "designers", "engineers", "relateds", "theme"}
	metaOfPostLookupFields       = []string{"heroImage", "leading_image_portrait", "categories", "tags", "topics", "og_image", "theme"}
	fullTopicLookupFields        = []string{"topic_relateds", "leading_image", "leading_image_portrait", "leading_video", "og_image"}
	metaOfTopicLookupFields      = []string{"leading_image", "leading_image_portrait", "og_image"}
	topicRelatedPostLookupFields = []string{"heroImage", "categories", "tags", "og_image"}

	metaOfPostExcludedFields       = []string{"leading_video", "writters", "photographers", "designers", "engineers", "relateds", "content"}
	metaOfTopicExcludedFields      = []string{"relateds", "leading_video"}
	topicRelatedPostExcludedFields = []string{"leading_image_portrait", "leading_video", "topics", "writters", "photographers", "designers", "engineers", "relateds", "theme", "content"}
)

type mongoStorage struct {
	*mongo.Client
}

func NewMongoStorage(client *mongo.Client) *mongoStorage {
	return &mongoStorage{client}
}

func (m *mongoStorage) GetPosts(ctx context.Context, q *Query) ([]Post, error) {
	var posts []Post
	// build aggregate stage from query
	stages := buildFilterStage(q.Filter)

	stages = append(stages, buildSortStage(q.Sort)...)

	stages = append(stages, buildPaginationStage(q.Pagination)...)

	// build expansion stages according to full/meta expansion
	if q.Full {
		stages = append(stages, buildLookupStages(fullPostLookupFields)...)
	} else {
		stages = append(stages, buildExcludedStage(metaOfPostExcludedFields))
		stages = append(stages, buildLookupStages(metaOfPostLookupFields)...)
	}

	cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection("posts").Aggregate(ctx, stages)
	if err != nil {
		return []Post{}, errors.WithStack(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post Post
		err := cursor.Decode(&post)
		if err != nil {
			return []Post{}, errors.WithStack(err)
		}
		posts = append(posts, post)
	}

	// Perform the query
	// error handling
	return posts, nil
}

func (m *mongoStorage) GetTopics(ctx context.Context, q *Query) ([]Topic, error) {
	var topics []Topic

	// build aggregate stage from query
	stages := buildFilterStage(q.Filter)

	stages = append(stages, buildSortStage(q.Sort)...)

	stages = append(stages, buildPaginationStage(q.Pagination)...)

	// build expansion stages according to full/meta expansion
	if q.Full {
		stages = append(stages, buildLookupStages(fullTopicLookupFields)...)
	} else {
		stages = append(stages, buildExcludedStage(metaOfTopicExcludedFields))
		stages = append(stages, buildLookupStages(metaOfTopicLookupFields)...)
	}

	cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection("topics").Aggregate(ctx, stages)
	if err != nil {
		return []Topic{}, errors.WithStack(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var topic Topic
		err := cursor.Decode(&topic)
		if err != nil {
			return []Topic{}, errors.WithStack(err)
		}
		topics = append(topics, topic)
	}
	// Perform the query
	// error handling
	return topics, nil
}

func (m *mongoStorage) GetPostCount(ctx context.Context, f *Filter) (int, error) {
	return 0, nil
}

func (m *mongoStorage) GetTopicCount(ctx context.Context, f *Filter) (int, error) {
	return 0, nil
}

func buildFilterStage(f Filter) []bson.D {
	var match []primitive.E
	if f.Slug != "" {
		match = append(match, bson.E{Key: "slug", Value: f.Slug})
	}

	if f.State != "" {
		match = append(match, bson.E{Key: "state", Value: f.State})
	}

	if f.Style != "" {
		match = append(match, bson.E{Key: "style", Value: f.Style})
	}

	if len(f.Categories) > 0 {
		var ids bson.A
		for _, v := range f.Categories {
			ids = append(ids, v)
		}
		match = append(match, bson.E{Key: "categories",
			Value: bson.D{{Key: "$in", Value: ids}}})
	}

	if !f.IsFeatured.IsZero() {
		if f.IsFeatured.Bool {
			match = append(match, bson.E{Key: "isFeatured",
				Value: f.IsFeatured.Bool})
		}
	}

	if !f.PublishedDate.IsEmpty() {
		var query bson.E
		if !f.PublishedDate.Exact.IsZero() {
			if !f.PublishedDate.Exact.IsZero() {
				t := time.Unix(f.PublishedDate.Start.Int64, 0)
				pt := primitive.NewDateTimeFromTime(t)
				query = bson.E{Key: "publishedDate", Value: pt}
			}
		} else {
			timeRange := []bson.E{}
			if !f.PublishedDate.Start.IsZero() {
				t := time.Unix(f.PublishedDate.Start.Int64, 0)
				pt := primitive.NewDateTimeFromTime(t)
				timeRange = append(timeRange, bson.E{Key: "$gte", Value: pt})
			}
			if !f.PublishedDate.End.IsZero() {
				t := time.Unix(f.PublishedDate.End.Int64, 0)
				pt := primitive.NewDateTimeFromTime(t)
				timeRange = append(timeRange, bson.E{Key: "$lte", Value: pt})
			}
			query = bson.E{Key: "publishedDate", Value: timeRange}
		}
		match = append(match, query)
	}

	return []bson.D{{{Key: "$match", Value: match}}}
}
func buildExcludedStage(excludeds []string) bson.D {
	var fields []bson.E
	for _, e := range excludeds {
		fields = append(fields, bson.E{Key: e, Value: 0})
	}

	return bson.D{{Key: "$project", Value: fields}}
}

func buildLookupStages(fields []string) []bson.D {
	var stages []bson.D
	for _, v := range fields {
		switch v {
		case "writters", "photographers", "designers", "engineers":
			stages = append(stages, buildLookupByIDStage(v, "contacts"))
		case "heroImage", "leading_image", "leading_image_portrait", "og_image":
			stages = append(stages, buildLookupByIDStage(v, "images"))
			stages = append(stages, buildUnwindStage(v))
		case "leading_video":
			stages = append(stages, buildLookupByIDStage(v, "videos"))
			stages = append(stages, buildUnwindStage(v))
		case "theme":
			stages = append(stages, buildLookupByIDStage(v, "themes"))
			stages = append(stages, buildUnwindStage(v))
		case "categories":
			stages = append(stages, buildLookupByIDStage(v, "postcategories"))
		case "tags":
			stages = append(stages, buildLookupByIDStage(v, "tags"))
		case "relateds":
			stages = append(stages, buildLookupNestedStage(v, "posts", metaOfPostExcludedFields, metaOfPostLookupFields, true))
		case "topic_relateds":
			stages = append(stages, buildLookupNestedStage("relateds", "posts", topicRelatedPostExcludedFields, topicRelatedPostLookupFields, true))
		case "topics":
			stages = append(stages, buildLookupNestedStage(v, "topics", metaOfTopicExcludedFields, metaOfTopicLookupFields, false))
			stages = append(stages, buildUnwindStage(v))

		case "topic_full":
			stages = append(stages, buildLookupNestedStage("topics", "topics", nil, fullTopicLookupFields, false))

			stages = append(stages, buildUnwindStage("topics"))
		}

	}
	return stages
}

func buildLookupByIDStage(field, from string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: from},
			{Key: "localField", Value: field},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: field},
		},
		},
	}
}

func buildLookupNestedStage(field, from string, excluded, nestedLookup []string, isArrayField bool) bson.D {
	var nestedPipeline []bson.D
	if isArrayField {
		nestedPipeline = append(nestedPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$in", Value: bson.A{"$_id", "$$" + field}}}}}}})
	} else {
		nestedPipeline = append(nestedPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id", "$$" + field}}}}}}})
	}
	if len(excluded) > 0 {
		nestedPipeline = append(nestedPipeline, buildExcludedStage(excluded))
	}
	nestedPipeline = append(nestedPipeline, buildLookupStages(nestedLookup)...)

	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: from},
			{Key: "let", Value: bson.D{{Key: field, Value: "$" + field}}},
			{Key: "pipeline", Value: nestedPipeline},
			{Key: "as", Value: field},
		},
		},
	}
}

// Build sort stage with respect to the specified field
// Currently, only single sorting criteria is allowed.
func buildSortStage(sort Sort) []bson.D {
	var sortBy bson.D
	if !sort.PublishedDate.IsAsc.IsZero() {
		if sort.PublishedDate.IsAsc.Bool {
			sortBy = bson.D{{Key: "publishedDate", Value: 1}}
		} else {
			sortBy = bson.D{{Key: "publishedDate", Value: -1}}
		}
	}

	if !sort.UpdatedAt.IsAsc.IsZero() {
		if sort.UpdatedAt.IsAsc.Bool {
			sortBy = bson.D{{Key: "UpdatedAt", Value: 1}}
		} else {
			sortBy = bson.D{{Key: "UpdatedAt", Value: -1}}
		}
	}

	return []bson.D{{{Key: "$sort", Value: sortBy}}}
}

func buildPaginationStage(p Pagination) []bson.D {
	var stages []bson.D
	if p.Offset > 0 {
		stages = append(stages, bson.D{{Key: "$skip", Value: p.Offset}})
	}
	if p.Limit > 0 {
		stages = append(stages, bson.D{{Key: "$limit", Value: p.Limit}})
	}

	return stages
}

func buildUnwindStage(field string) bson.D {
	return bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$" + field},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}
}
