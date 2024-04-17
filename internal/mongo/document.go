package mongo

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// BuildDocument is a wrapper for constructing a mongo document
// e.g. {"$sort": {"name": 1}}
//       <-key->  <- value ->
func BuildDocument(key string, value interface{}) bson.D {
	return bson.D{{Key: key, Value: value}}
}

// Build Element is a wrapper for constructing part of an object
// e.g. {"$sort": {"name":      1    }}
//                 <-key->  <-value->
func BuildElement(key string, value interface{}) bson.E {
	return bson.E{Key: key, Value: value}
}

func BuildArray(items interface{}) (arr bson.A, exist bool) {
	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return
	}

	arrLength := val.Len()
	if arrLength > 0 {
		exist = true
	}

	for i := 0; i < arrLength; i++ {
		arr = append(arr, val.Index(i).Interface())
	}
	return
}

func BuildLookupByIDStage(field, fromCol string) bson.D {
	return bson.D{
		{Key: StageLookup, Value: bson.D{
			{Key: MetaFrom, Value: fromCol},
			{Key: MetaLocalField, Value: field},
			{Key: MetaForeignField, Value: "_id"},
			{Key: MetaAs, Value: field},
		},
		},
	}
}

func BuildUnwindStage(field string) bson.D {
	return bson.D{{Key: StageUnwind, Value: bson.D{
		{Key: MetaPath, Value: "$" + field},
		{Key: MetaPreserveNullAndEmptyArrays, Value: true},
	}}}
}

func BuildCategorySetStage() []bson.D {
	var result []bson.D

	// unwind category_set
	result = append(result, bson.D{{Key: StageUnwind, Value: bson.D{
		{Key: "path", Value: "$category_set"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}})

	// lookup postcategories
	result = append(result, bson.D{{Key: StageLookup, Value: bson.D{
		{Key: MetaFrom, Value: "postcategories"},
		{Key: MetaLocalField, Value: "category_set.category"},
		{Key: MetaForeignField, Value: "_id"},
		{Key: MetaAs, Value: "category_set.category"},
	}}})

	// lookup tags
	result = append(result, bson.D{{Key: StageLookup, Value: bson.D{
		{Key: MetaFrom, Value: "tags"},
		{Key: MetaLocalField, Value: "category_set.subcategory"},
		{Key: MetaForeignField, Value: "_id"},
		{Key: MetaAs, Value: "category_set.subcategory"},
	}}})

	// addFields
	result = append(result, bson.D{{Key: StageAddFields, Value: bson.D{
		{Key: "category_set.category", Value: bson.D{
			{Key: "$arrayElemAt", Value: bson.A{"$category_set.category", 0}},
		}},
		{Key: "category_set.subcategory", Value: bson.D{
			{Key: "$arrayElemAt", Value: bson.A{"$category_set.subcategory", 0}},
		}},
	}}})

	// group
	result = append(result, bson.D{{Key: StageGroup, Value: bson.D{
		{Key: "_id", Value: "$_id"},
		{Key: "category_set", Value: bson.D{{Key: "$push", Value: "$category_set"}}},
		{Key: "data", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
	}}})

	// addFields
	result = append(result, bson.D{{Key: StageAddFields, Value: bson.D{
		{Key: "data.category_set", Value: "$category_set"},
	}}})

	// replaceRoot
	result = append(result, bson.D{{Key: StageReplaceRoot, Value: bson.D{
		{Key: "newRoot", Value: "$data"},
	}}})

	return result
}

func BuildReviewLookupStatements() []bson.D {
	var stages []bson.D

	// lookup posts
	stages = append(stages, bson.D{{
		Key: StageLookup, Value: bson.D{
			{Key: MetaFrom, Value: "posts"},
			{Key: MetaLocalField, Value: "post_id"},
			{Key: MetaForeignField, Value: "_id"},
			{Key: MetaAs, Value: "post"},
		},
	}})

	// lookup images
	stages = append(stages, bson.D{{
		Key: StageLookup, Value: bson.D{
			{Key: MetaFrom, Value: "images"},
			{Key: MetaLocalField, Value: "post.og_image"},
			{Key: MetaForeignField, Value: "_id"},
			{Key: MetaAs, Value: "og_image"},
		},
	}})

	// unwind images
	stages = append(stages, BuildUnwindStage("og_image"))
	stages = append(stages, BuildUnwindStage("post"))

	// project fields
	stages = append(stages, bson.D{{
		Key: StageProject, Value: bson.D{
			{Key: "order", Value: 1},
			{Key: "og_image", Value: 1},
			{Key: "post_id", Value: "$post._id"},
			{Key: "slug", Value: "$post.slug"},
			{Key: "title", Value: "$post.title"},
			{Key: "og_description", Value: "$post.og_description"},
			{Key: "reviewWord", Value: "$post.reviewWord"},
		},
	}})

	return stages
}
