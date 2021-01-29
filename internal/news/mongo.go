package news

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/twreporter/go-api/internal/mongo"
	"github.com/twreporter/go-api/internal/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/guregu/null.v3"
)

// tagMongo is used to map the query field to the corresponded field in real mongo document
const tagMongo = "mongo"

type mongoQuery struct {
	mongoPagination
	mongoFilter
	mongoSort
}

func (mq *mongoQuery) GetFilter() mongoFilter {
	return mq.mongoFilter
}

// BuildQueryStatments build query statements from pagination/filter/sort objects
// TODO(babygoat): this can be further refactor to accept generic storage query interface
func BuildQueryStatements(mq *mongoQuery) []bson.D {
	var stages []bson.D

	stages = append(stages, mq.mongoFilter.BuildStage()...)
	stages = append(stages, mq.mongoSort.BuildStage()...)
	stages = append(stages, mq.mongoPagination.BuildStage()...)

	return stages
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

func (mp mongoPagination) BuildStage() []bson.D {
	var stages []bson.D
	if mp.Skip > 0 {
		stages = append(stages, mongo.BuildDocument(mongo.StageSkip, mp.Skip))
	}
	if mp.Limit > 0 {
		stages = append(stages, mongo.BuildDocument(mongo.StageLimit, mp.Limit))
	}

	return stages
}

type mongoFilter struct {
	Slug       string               `mongo:"slug"`
	State      string               `mongo:"state"`
	Style      string               `mongo:"style"`
	IsFeatured null.Bool            `mongo:"isFeatured"`
	Categories []primitive.ObjectID `mongo:"categories"`
	Tags       []primitive.ObjectID `mongo:"tags"`
	IDs        []primitive.ObjectID `mongo:"_id"`
	Name       primitive.Regex      `mongo:"name"`
}

func (mf mongoFilter) BuildStage() []bson.D {
	var match []bson.D
	if elements := mf.BuildElements(); len(elements) > 0 {
		match = append(match, mongo.BuildDocument(mongo.StageMatch, elements))
	}
	return match
}

func (mf mongoFilter) BuildElements() []bson.E {
	typ := reflect.TypeOf(mf)
	val := reflect.ValueOf(mf)

	var elements []bson.E
	for i := 0; i < typ.NumField(); i++ {
		fieldT := typ.Field(i)
		fieldV := val.Field(i)

		tag := fieldT.Tag.Get(tagMongo)

		switch fieldV.Interface().(type) {
		case string:
			v := fieldV.Interface().(string)
			if v != "" {
				elements = append(elements, mongo.BuildElement(tag, v))
			}
		case null.Bool:
			v := fieldV.Interface().(null.Bool)
			if !v.IsZero() {
				elements = append(elements, mongo.BuildElement(tag, v.Bool))
			}

		case []primitive.ObjectID:
			if v, ok := mongo.BuildArray(fieldV.Interface().([]primitive.ObjectID)); ok {
				elements = append(elements, mongo.BuildElement(tag, mongo.BuildDocument(mongo.OpIn, v)))
			}
		case primitive.Regex:
			v := fieldV.Interface().(primitive.Regex)
			if v.Pattern != "" {
				elements = append(elements, mongo.BuildElement(tag, v))
			}
		default:
			log.Errorf("Unimplemented type %+v", fieldT.Type)
		}
	}
	return elements
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
		Name:       primitive.Regex{Pattern: f.Name},
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

func (ms mongoSort) BuildStage() []bson.D {
	typ := reflect.TypeOf(ms)
	val := reflect.ValueOf(ms)

	var sortBy []bson.E
	for i := 0; i < typ.NumField(); i++ {
		fieldT := typ.Field(i)
		fieldV := val.Field(i)

		tag := fieldT.Tag.Get(tagMongo)

		switch fieldV.Interface().(type) {
		case query.Order:
			v := fieldV.Interface().(query.Order)
			if !v.IsAsc.IsZero() {
				if v.IsAsc.Bool {
					sortBy = append(sortBy, mongo.BuildElement(tag, mongo.OrderAsc))
				} else {
					sortBy = append(sortBy, mongo.BuildElement(tag, mongo.OrderDesc))
				}
			}
		default:
		}
	}

	if len(sortBy) == 0 {
		return nil
	}

	return []bson.D{mongo.BuildDocument(mongo.StageSort, sortBy)}
}

func fromSort(s SortBy) mongoSort {
	return mongoSort{
		PublishedDate: s.PublishedDate,
		UpdatedAt:     s.UpdatedAt,
	}
}

const (
	ColContacts       = "contacts"
	ColImages         = "images"
	ColVideos         = "videos"
	ColThemes         = "themes"
	ColPostCategories = "postcategories"
	ColTags           = "tags"
	ColPosts          = "posts"
	ColTopics         = "topics"

	// TODO: rename fields to writer
	fieldWriters              = "writters"
	fieldPhotographers        = "photographers"
	fieldDesigners            = "designers"
	fieldEngineers            = "engineers"
	fieldHeroImage            = "heroImage"
	fieldLeadingImage         = "leading_image"
	fieldLeadingImagePortrait = "leading_image_portrait"
	fieldOgImage              = "og_image"
	fieldLeadingVideo         = "leading_video"
	fieldTheme                = "theme"
	fieldCategories           = "categories"
	fieldTags                 = "tags"
	// TODO: rename the field to topic
	fieldTopics           = "topics"
	fieldRelatedDocuments = "relateds"
	fieldID               = "_id"
	fieldState            = "state"
	fieldThumbnail        = "thumbnail"
	fieldBio              = "bio"
)

type lookupInfo struct {
	Collection string
	ToUnwind   bool
}

var (
	LookupFullPost = map[string]lookupInfo{
		fieldCategories:           {Collection: ColPostCategories},
		fieldDesigners:            {Collection: ColContacts},
		fieldEngineers:            {Collection: ColContacts},
		fieldHeroImage:            {Collection: ColImages, ToUnwind: true},
		fieldLeadingImagePortrait: {Collection: ColImages, ToUnwind: true},
		fieldOgImage:              {Collection: ColImages, ToUnwind: true},
		fieldPhotographers:        {Collection: ColContacts},
		fieldTags:                 {Collection: ColTags},
		fieldTopics:               {Collection: ColTopics, ToUnwind: true},
		fieldWriters:              {Collection: ColContacts},
	}

	LookupMetaOfPost = map[string]lookupInfo{
		fieldCategories:           {Collection: ColPostCategories},
		fieldHeroImage:            {Collection: ColImages, ToUnwind: true},
		fieldLeadingImagePortrait: {Collection: ColImages, ToUnwind: true},
		fieldTags:                 {Collection: ColTags},
		fieldOgImage:              {Collection: ColImages, ToUnwind: true},
	}

	LookupFullTopic = map[string]lookupInfo{
		fieldLeadingImage:         {Collection: ColImages, ToUnwind: true},
		fieldLeadingImagePortrait: {Collection: ColImages, ToUnwind: true},
		fieldLeadingVideo:         {Collection: ColVideos, ToUnwind: true},
		fieldOgImage:              {Collection: ColImages, ToUnwind: true},
	}

	LookupMetaOfTopic = map[string]lookupInfo{
		fieldLeadingImage:         {Collection: ColImages, ToUnwind: true},
		fieldLeadingImagePortrait: {Collection: ColImages, ToUnwind: true},
		fieldOgImage:              {Collection: ColImages, ToUnwind: true},
	}

	LookupAuthor = map[string]lookupInfo{
		fieldThumbnail: {Collection: ColImages, ToUnwind: true},
	}
)

func BuildLookupStatements(m map[string]lookupInfo) []bson.D {
	var stages []bson.D
	for field, info := range m {
		if shouldPreserveOrder(field) {
			stages = append(stages, buildPreserveLookupOrderStatement(field, info)...)
		} else {
			stages = append(stages, mongo.BuildLookupByIDStage(field, info.Collection))
		}
		if info.ToUnwind {
			stages = append(stages, mongo.BuildUnwindStage(field))
		}
	}
	return stages
}

// Filter related posts that is published already
func BuildFilterRelatedPost() []bson.D {
	var stages []bson.D

	// TODO(babygoat): add test item ensure the order of related document
	// First, retrieve full post data by joing the documents.
	// Per this mongodb issue (https://jira.mongodb.org/browse/SERVER-32947),
	// the order of the output documents after $lookup operation is in natural order (i.e. internal disk order)
	// and thus the order does not persist.
	// As the order of related documents is crucial, we need to keep a copy of the array for order reference.
	stages = append(stages, bson.D{
		{Key: mongo.StageAddFields, Value: bson.D{{Key: "relatedsCopy", Value: "$" + fieldRelatedDocuments}}},
	})
	stages = append(stages, mongo.BuildLookupByIDStage(fieldRelatedDocuments, ColPosts))

	// Then, match the posts with published state
	stages = append(stages, bson.D{
		{Key: mongo.StageAddFields, Value: bson.D{
			{Key: fieldRelatedDocuments, Value: bson.D{
				{Key: mongo.StageFilter, Value: bson.D{
					{Key: mongo.MetaInput, Value: "$" + fieldRelatedDocuments},
					{Key: mongo.MetaAs, Value: fieldRelatedDocuments},
					{Key: mongo.MetaCond, Value: bson.D{
						{Key: mongo.OpEq, Value: bson.A{
							"$$" + fieldRelatedDocuments + "." + fieldState,
							"published",
						}},
					}},
				}},
			},
			},
		}},
	})

	// Next, promote the _id field to array of ObjectIDs
	stages = append(stages, bson.D{
		{Key: mongo.StageAddFields, Value: mongo.BuildDocument(fieldRelatedDocuments, "$"+fieldRelatedDocuments+"."+fieldID)},
	})

	// Finally, constructing the ordered relateds field by
	// filtering the copied array with the output array of the previous stage
	// as filter operator returns elements in original order.
	stages = append(stages, bson.D{
		{Key: mongo.StageAddFields, Value: bson.D{
			{Key: fieldRelatedDocuments, Value: bson.D{
				{Key: mongo.StageFilter, Value: bson.D{
					{Key: mongo.MetaInput, Value: "$relatedsCopy"},
					{Key: mongo.MetaAs, Value: "relatedsCopy"},
					{Key: mongo.MetaCond, Value: bson.D{
						{Key: mongo.OpIn, Value: bson.A{
							"$$relatedsCopy",
							"$" + fieldRelatedDocuments,
						}},
					}},
				}},
			},
			},
		}},
	})
	return stages
}

func buildPreserveLookupOrderStatement(orderedField string, info lookupInfo) []bson.D {
	var stages []bson.D

	// Copy the original orderedField for order reference
	copyField := orderedField + "Copy"
	stages = append(stages, bson.D{
		{Key: mongo.StageAddFields, Value: bson.D{{Key: copyField, Value: "$" + orderedField}}},
	})

	// Perform lookup for joined documents
	stages = append(stages, mongo.BuildLookupByIDStage(orderedField, info.Collection))

	// Construct the ordered documents from reference of the original copy
	// An example is given below for illustration of the query performed
	// Field writers
	// {
	//   $addFields: {
	//     writers: {
	// Use the reduce operator so we can push(concat) the element in correct order(i.e. input array)
	// from left to right
	// https://docs.mongodb.com/manual/reference/operator/aggregation/reduce/
	//       $reduce: {
	//         input: "writersCopy",
	//         initialValue: [],
	// Declare the expression used to apply on the input array element from left to right
	// Two pre-defined variables are available.
	// $$this variable refers to the element in input array field.
	// $$value variable refers to the cumulative value of the expression starting from initialValue
	//         in: {
	// Declare accessible varialbes in below "in" expression
	//           $let: {
	//             vars: {writersCopy: "$$this"},
	//             in: {
	//               $concatArrays:[
	//                 "$$value",
	// Filter the element from the (joined) ordered field
	//                 {
	//                   $filter: {
	//                     input: "$writers",
	//                     as: "writers",
	//                     cond:{ $eq:["$$writersCopy", "$$writers._id"] }
	//                   }
	//                 }
	//               ]
	//             }
	//           }
	//         }
	//       }
	//     }
	//   }
	// }
	stages = append(stages, bson.D{
		{Key: mongo.StageAddFields, Value: bson.D{
			{Key: orderedField, Value: bson.D{
				{Key: mongo.OpReduce, Value: bson.D{
					{Key: mongo.MetaInput, Value: "$" + copyField},
					{Key: mongo.MetaInitialValue, Value: bson.A{}},
					{Key: mongo.MetaIn, Value: bson.D{
						{Key: mongo.OpLet, Value: bson.D{
							{Key: mongo.MetaVars, Value: bson.D{{Key: copyField, Value: "$$this"}}},
							{Key: mongo.MetaIn, Value: bson.D{
								{Key: mongo.OpConcatArrays, Value: bson.A{
									"$$value",
									bson.D{
										{Key: mongo.StageFilter, Value: bson.D{
											{Key: mongo.MetaInput, Value: "$" + orderedField},
											{Key: mongo.MetaAs, Value: orderedField},
											{Key: mongo.MetaCond, Value: bson.D{
												{Key: mongo.OpEq, Value: bson.A{
													"$$" + copyField,
													"$$" + orderedField + "." + fieldID,
												}},
											}},
										}},
									},
								}},
							}},
						}},
					}},
				}},
			}},
		}},
	})

	return stages
}

var (
	// Fields that should be preserved order after lookup stage according to the requirements
	orderedFields = []string{fieldDesigners, fieldEngineers, fieldPhotographers, fieldWriters}
)

func shouldPreserveOrder(field string) bool {
	for _, v := range orderedFields {
		if field == v {
			return true
		}
	}
	return false
}

// BuildBioMarkdownOnlyStatement returns statement for rewriting `bio` field with markdown format
func BuildBioMarkdownOnlyStatement() bson.D {
	return bson.D{{Key: mongo.StageAddFields, Value: bson.D{{Key: fieldBio, Value: "$" + fieldBio + ".md"}}}}
}
