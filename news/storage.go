package news

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"twreporter.org/go-api/globals"
)

var (
	fullPostLookupFields         = []string{"hero_image", "leading_image_portrait", "leading_video", "categories", "tags", "topic_full", "og_image", "writters", "photographers", "designers", "engineers", "relateds", "theme"}
	metaOfPostLookupFields       = []string{"hero_image", "leading_image_portrait", "categories", "tags", "topic", "og_image", "theme"}
	fullTopicLookupFields        = []string{"topic_relateds", "leading_image", "leading_image_portrait", "leading_video", "og_image"}
	metaOfTopicLookupFields      = []string{"leading_image", "leading_image_portrait", "og_image"}
	topicRelatedPostLookupFields = []string{"hero_image", "categories", "tags", "og_image"}
)

type mongoStorage struct {
	*mongo.Client
}

func NewMongoStorage(client *mongo.Client) *mongoStorage {
	return &mongoStorage{client}
}

func (m *mongoStorage) GetPosts(ctx context.Context, q *Query) ([]Post, error) {
	posts := make([]Post, q.Limit)
	// build aggregate stage from query
	stages := buildFilterStage(q.Filter)

	stages = append(stages, buildSortStage(q.Sort)...)

	stages = append(stages, buildPaginationStage(q.Pagination)...)

	// build expansion stages according to full/meta expansion
	var fields []string
	if q.Full {
		fields = fullPostLookupFields
	} else {
		fields = metaOfPostLookupFields
	}

	stages = append(stages, buildLookupStages(fields)...)
	cursor, err := m.Database(globals.Conf.DB.Mongo.DBname).Collection("posts").Aggregate(ctx, stages)
	if err != nil {
		return []Post{}, errors.WithStack(err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &posts); err != nil {
		return []Post{}, errors.WithStack(err)
	}
	// Perform the query
	// error handling
	return posts, nil
}

func (m *mongoStorage) GetTopics(ctx context.Context, q *Query) ([]Topic, error) {
	return nil, nil
}

func (m *mongoStorage) GetPostCount(ctx context.Context, q *Query) (int, error) {
	return 0, nil
}

func (m *mongoStorage) GetTopicCount(ctx context.Context, q *Query) (int, error) {
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

	return []bson.D{{{Key: "$match", Value: match}}}
}

func buildLookupStages(fields []string) []bson.D {
	var stages []bson.D
	for _, v := range fields {
		switch v {
		case "writters", "photographers", "designers", "engineers":
			stages = append(stages, buildLookupByIDStage(v, "contacts", false)...)
		case "hero_image", "leading_image", "leading_image_portrait", "og_image":
			stages = append(stages, buildLookupByIDStage(v, "images", true)...)
		case "leading_video":
			stages = append(stages, buildLookupByIDStage(v, "videos", true)...)
		case "theme":
			stages = append(stages, buildLookupByIDStage(v, "themes", true)...)
		case "categories":
			stages = append(stages, buildLookupByIDStage(v, "postcategories", true)...)
		case "tags":
			stages = append(stages, buildLookupByIDStage(v, "tags", false)...)

		case "relateds":
			nestedPipeline := []bson.D{
				// match stage
				// {$match: {$expr: {$in: ["$_id", "$$relateds"]}}}
				{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$in", Value: bson.A{"$_id", "$$" + v}}}}}}},
			}
			nestedPipeline = append(nestedPipeline, buildLookupStages(metaOfPostLookupFields)...)
			stages = append(stages, bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "posts"},
					{Key: "let", Value: bson.D{{Key: v, Value: "$" + v}}},
					{Key: "pipeline", Value: nestedPipeline},
					{Key: "as", Value: v},
				},
				},
			})
		case "topic_relateds":
			nestedPipeline := []bson.D{
				{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$in", Value: bson.A{"$_id", "$$relateds"}}}}}}},
			}
			nestedPipeline = append(nestedPipeline, buildLookupStages(topicRelatedPostLookupFields)...)
			stages = append(stages, bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "posts"},
					{Key: "let", Value: bson.D{{Key: "relateds", Value: "$relateds"}}},
					{Key: "pipeline", Value: nestedPipeline},
					{Key: "as", Value: "relateds"},
				},
				},
			})
		case "topic":
			nestedPipeline := []bson.D{
				{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$in", Value: bson.A{"$_id", "$$" + v}}}}}}},
			}
			nestedPipeline = append(nestedPipeline, buildLookupStages(metaOfTopicLookupFields)...)
			stages = append(stages, bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "topics"},
					{Key: "let", Value: bson.D{{Key: v, Value: "$" + v}}},
					{Key: "pipeline", Value: nestedPipeline},
					{Key: "as", Value: v},
				},
				},
			})

		case "topic_full":
			nestedPipeline := []bson.D{
				{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$in", Value: bson.A{"$_id", "$$topic"}}}}}}},
			}
			nestedPipeline = append(nestedPipeline, buildLookupStages(metaOfTopicLookupFields)...)
			stages = append(stages, bson.D{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "topics"},
					{Key: "let", Value: bson.D{{Key: "topic", Value: "$topic"}}},
					{Key: "pipeline", Value: nestedPipeline},
					{Key: "as", Value: "topic"},
				},
				},
			})
		}
	}
	return stages
}

func buildLookupByIDStage(field, from string, expand bool) []bson.D {
	var stages []bson.D

	stages = append(stages, bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: from},
			{Key: "localField", Value: field},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: field},
		},
		},
	})
	if expand {
		stages = append(stages, bson.D{{Key: "$unwind", Value: "$" + field}})

	}
	return stages
}

func buildSortStage(sort Sort) []bson.D {
	var sortBy bson.D
	if !sort.PublishedDate.IsAsc.IsZero() {
		if sort.PublishedDate.IsAsc.Bool {
			sortBy = bson.D{{Key: "publishedDate", Value: 1}}
		} else {
			sortBy = bson.D{{Key: "publishedDate", Value: -1}}
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
