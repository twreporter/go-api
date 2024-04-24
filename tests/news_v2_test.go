package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
	"io/ioutil"
	"encoding/json"

	"github.com/twreporter/go-api/internal/news"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/twreporter/go-api/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testPost struct {
	ID            primitive.ObjectID
	Editor        primitive.ObjectID
	CreatedAt     time.Time
	Slug          string
	State         string
	Topic         *primitive.ObjectID
	Image         primitive.ObjectID
	Video         primitive.ObjectID
	Relateds      []primitive.ObjectID
	Tags          []primitive.ObjectID
	Categories    []primitive.ObjectID
	Engineers     []primitive.ObjectID
	Designers     []primitive.ObjectID
	Photographers []primitive.ObjectID
	Writers       []primitive.ObjectID
	Category      string
	SubCategory   string
	BookmarkID    string
	ReviewWord    string
}

type testReview struct {
	ID     primitive.ObjectID
	PostID primitive.ObjectID
	Order  int
}

type responseBodyForReview struct {
	Status string        `json:"status"`
	Data   []news.Review `json:"data"`
}

func TestGetPostsByAuthors_AuthorIsAnEngineer(t *testing.T) {
	db, cleanup := setupMongoGoDriverTestDB()
	defer cleanup()
	defer func() { db.Drop(context.Background()) }()
	// setup records
	// authors
	authors := map[string]testAuthor{
		"王小明": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "王小明",
			createdAt: time.Unix(1611817200, 0),
		},
		"劉大華": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "劉大華",
			createdAt: time.Unix(1611817800, 0),
		},
	}
	for _, author := range authors {
		migrateAuthorRecord(db, author)
	}
	// posts
	posts := map[string]testPost{
		"王小明的文章": {
			ID:     primitive.NewObjectID(),
			Editor: primitive.NewObjectID(),
			// Without loss of generosity,
			// use dedicate timestamp to ensure the JSON equality
			// with no decimal points precision
			CreatedAt:  time.Unix(1612337400, 0),
			Slug:       "test-slug-1",
			State:      "published",
			Image:      primitive.NewObjectID(),
			Video:      primitive.NewObjectID(),
			Engineers:  []primitive.ObjectID{authors["王小明"].id},
			Categories: []primitive.ObjectID{primitive.NewObjectID()},
			Tags:       []primitive.ObjectID{primitive.NewObjectID()},
		},
		"劉大華的文章": {
			ID:         primitive.NewObjectID(),
			Editor:     primitive.NewObjectID(),
			CreatedAt:  time.Unix(1612337400, 0),
			Slug:       "test-slug-2",
			State:      "published",
			Image:      primitive.NewObjectID(),
			Video:      primitive.NewObjectID(),
			Engineers:  []primitive.ObjectID{authors["劉大華"].id},
			Categories: []primitive.ObjectID{primitive.NewObjectID()},
			Tags:       []primitive.ObjectID{primitive.NewObjectID()},
		},
	}
	for _, post := range posts {
		migratePostRecord(db, post)
	}
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/authors/%s/posts", authors["王小明"].id.Hex()), "", "", "")
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, postListResponse(metaOfPostResponse(posts["王小明的文章"])), response.Body.String())
}

func TestGetPostsByAuthors_AuthorIsADesigner(t *testing.T) {
	db, cleanup := setupMongoGoDriverTestDB()
	defer cleanup()
	defer func() { db.Drop(context.Background()) }()
	// setup records
	// authors
	authors := map[string]testAuthor{
		"王小明": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "王小明",
			createdAt: time.Unix(1611817200, 0),
		},
		"劉大華": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "劉大華",
			createdAt: time.Unix(1611817800, 0),
		},
	}
	for _, author := range authors {
		migrateAuthorRecord(db, author)
	}
	// posts
	posts := map[string]testPost{
		"王小明的文章": {
			ID:     primitive.NewObjectID(),
			Editor: primitive.NewObjectID(),
			// Without loss of generosity,
			// use dedicate timestamp to ensure the JSON equality
			// with no decimal points precision
			CreatedAt:  time.Unix(1612337400, 0),
			Slug:       "test-slug-1",
			State:      "published",
			Image:      primitive.NewObjectID(),
			Video:      primitive.NewObjectID(),
			Designers:  []primitive.ObjectID{authors["王小明"].id},
			Categories: []primitive.ObjectID{primitive.NewObjectID()},
			Tags:       []primitive.ObjectID{primitive.NewObjectID()},
		},
		"劉大華的文章": {
			ID:         primitive.NewObjectID(),
			Editor:     primitive.NewObjectID(),
			CreatedAt:  time.Unix(1612337400, 0),
			Slug:       "test-slug-2",
			State:      "published",
			Image:      primitive.NewObjectID(),
			Video:      primitive.NewObjectID(),
			Designers:  []primitive.ObjectID{authors["劉大華"].id},
			Categories: []primitive.ObjectID{primitive.NewObjectID()},
			Tags:       []primitive.ObjectID{primitive.NewObjectID()},
		},
	}
	for _, post := range posts {
		migratePostRecord(db, post)
	}
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/authors/%s/posts", authors["王小明"].id.Hex()), "", "", "")
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, postListResponse(metaOfPostResponse(posts["王小明的文章"])), response.Body.String())
}

func TestGetPostsByAuthors_AuthorIsAPhotographer(t *testing.T) {
	db, cleanup := setupMongoGoDriverTestDB()
	defer cleanup()
	defer func() { db.Drop(context.Background()) }()
	// setup records
	// authors
	authors := map[string]testAuthor{
		"王小明": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "王小明",
			createdAt: time.Unix(1611817200, 0),
		},
		"劉大華": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "劉大華",
			createdAt: time.Unix(1611817800, 0),
		},
	}
	for _, author := range authors {
		migrateAuthorRecord(db, author)
	}
	// posts
	posts := map[string]testPost{
		"王小明的文章": {
			ID:     primitive.NewObjectID(),
			Editor: primitive.NewObjectID(),
			// Without loss of generosity,
			// use dedicate timestamp to ensure the JSON equality
			// with no decimal points precision
			CreatedAt:     time.Unix(1612337400, 0),
			Slug:          "test-slug-1",
			State:         "published",
			Image:         primitive.NewObjectID(),
			Video:         primitive.NewObjectID(),
			Photographers: []primitive.ObjectID{authors["王小明"].id},
			Categories:    []primitive.ObjectID{primitive.NewObjectID()},
			Tags:          []primitive.ObjectID{primitive.NewObjectID()},
		},
		"劉大華的文章": {
			ID:            primitive.NewObjectID(),
			Editor:        primitive.NewObjectID(),
			CreatedAt:     time.Unix(1612337400, 0),
			Slug:          "test-slug-2",
			State:         "published",
			Image:         primitive.NewObjectID(),
			Video:         primitive.NewObjectID(),
			Photographers: []primitive.ObjectID{authors["劉大華"].id},
			Categories:    []primitive.ObjectID{primitive.NewObjectID()},
			Tags:          []primitive.ObjectID{primitive.NewObjectID()},
		},
	}
	for _, post := range posts {
		migratePostRecord(db, post)
	}
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/authors/%s/posts", authors["王小明"].id.Hex()), "", "", "")
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, postListResponse(metaOfPostResponse(posts["王小明的文章"])), response.Body.String())
}

func TestGetPostsByAuthors_AuthorIsAWriter(t *testing.T) {
	db, cleanup := setupMongoGoDriverTestDB()
	defer cleanup()
	defer func() { db.Drop(context.Background()) }()
	// setup records
	// authors
	authors := map[string]testAuthor{
		"王小明": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "王小明",
			createdAt: time.Unix(1611817200, 0),
		},
		"劉大華": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "劉大華",
			createdAt: time.Unix(1611817800, 0),
		},
	}
	for _, author := range authors {
		migrateAuthorRecord(db, author)
	}
	// posts
	posts := map[string]testPost{
		"王小明的文章": {
			ID:     primitive.NewObjectID(),
			Editor: primitive.NewObjectID(),
			// Without loss of generosity,
			// use dedicate timestamp to ensure the JSON equality
			// with no decimal points precision
			CreatedAt:   time.Unix(1612337400, 0),
			Slug:        "test-slug-1",
			State:       "published",
			Image:       primitive.NewObjectID(),
			Video:       primitive.NewObjectID(),
			Writers:     []primitive.ObjectID{authors["王小明"].id},
			Categories:  []primitive.ObjectID{primitive.NewObjectID()},
			Tags:        []primitive.ObjectID{primitive.NewObjectID()},
			Category:    "",
			SubCategory: "",
		},
		"劉大華的文章": {
			ID:         primitive.NewObjectID(),
			Editor:     primitive.NewObjectID(),
			CreatedAt:  time.Unix(1612337400, 0),
			Slug:       "test-slug-2",
			State:      "published",
			Image:      primitive.NewObjectID(),
			Video:      primitive.NewObjectID(),
			Writers:    []primitive.ObjectID{authors["劉大華"].id},
			Categories: []primitive.ObjectID{primitive.NewObjectID()},
			Tags:       []primitive.ObjectID{primitive.NewObjectID()},
		},
	}
	for _, post := range posts {
		migratePostRecord(db, post)
	}
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/authors/%s/posts", authors["王小明"].id.Hex()), "", "", "")
	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, postListResponse(metaOfPostResponse(posts["王小明的文章"])), response.Body.String())
}

func TestGetPostReviews_Success(t *testing.T) {
	var resBody responseBodyForReview

	db, cleanup := setupMongoGoDriverTestDB()
	defer cleanup()
	defer func() { db.Drop(context.Background()) }()

	// setup post records
	posts := map[string]testPost{
		"mock1": {
			ID:          primitive.NewObjectID(),
			Slug:        "test-slug-1",
			Image:       primitive.NewObjectID(),
			ReviewWord:  "test review word 1",
		},
		"mock2": {
			ID:         primitive.NewObjectID(),
			Slug:       "test-slug-2",
			Image:      primitive.NewObjectID(),
			ReviewWord:  "test review word 2",
		},
	}
	for _, post := range posts {
		migratePostRecord(db, post)
	}

	// setup post review record
	reviews := map[string]testReview{
		"mock1": {
			ID:     primitive.NewObjectID(),
			Order:  1,
			PostID: posts["mock1"].ID,
		},
		"mock2": {
			ID:     primitive.NewObjectID(),
			Order:  2,
			PostID: posts["mock2"].ID,
		},
	}
	for _, review := range reviews {
		migratePostReviewRecord(db, review)
	}

	// Mocking user
	mockEmail := "get-post-review@twreporter.org"
	user := createUser(mockEmail)
	defer func() { deleteUser(user) }()
	authorization, cookie := helperSetupAuth(user)

	// Send request to test GetPostReviews function
	response := serveHTTPWithCookies(http.MethodGet, "/v2/post_reviews", "", "", authorization, cookie)
	resBodyInBytes, _ := ioutil.ReadAll(response.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 2, len(resBody.Data))
	assert.Equal(t, reviews["mock1"].PostID, resBody.Data[0].PostID)
	assert.Equal(t, reviews["mock1"].Order, resBody.Data[0].Order)
	assert.Equal(t, posts["mock1"].Slug, resBody.Data[0].Slug)
	assert.Equal(t, posts["mock1"].ReviewWord, resBody.Data[0].ReviewWord)
	assert.Equal(t, reviews["mock2"].PostID, resBody.Data[1].PostID)
	assert.Equal(t, reviews["mock2"].Order, resBody.Data[1].Order)
	assert.Equal(t, posts["mock2"].Slug, resBody.Data[1].Slug)
	assert.Equal(t, posts["mock2"].ReviewWord, resBody.Data[1].ReviewWord)
}

func TestGetPostReviews_AuthFail(t *testing.T) {
	// Mocking user
	mockEmail := "get-post-review@twreporter.org"
	user := createUser(mockEmail)
	defer func() { deleteUser(user) }()
	authorization, cookie := helperSetupAuth(user)

	// Test no authorization header
	response := serveHTTPWithCookies(http.MethodGet, "/v2/post_reviews", "", "", "", cookie)
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	// Test no cookie
	response = serveHTTP(http.MethodGet, "/v2/post_reviews", "", "", authorization)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

// setupMongoGoDriverTestDB overwrites the global mongo DB temporarily
// for testing news-v2 endpoints with mongo-go-driver
// it can be removed once the news v1 endpoints deprecated.
func setupMongoGoDriverTestDB() (*mongo.Database, func()) {
	// TODO(babygoat): remove pre variable to overwrite Mongo.DBname after v1 endpoints & tests are removed
	pre := globals.Conf.DB.Mongo.DBname
	globals.Conf.DB.Mongo.DBname = testMongoDB
	db := testMongoClient.Database(globals.Conf.DB.Mongo.DBname)
	cleanup := func() {
		// TODO(babygoat): remove pre variable to overwrite Mongo.DBname after v1 endpoints & tests are removed
		globals.Conf.DB.Mongo.DBname = pre
	}
	return db, cleanup
}

func postListResponse(posts ...string) string {
	return listResponse(len(posts), posts)
}

func migratePostRecord(db *mongo.Database, post testPost) {
	image := createImageDocument(post.Image)
	video := createVideoDocument(post.Video)
	var categories, tags []interface{}
	for _, tag := range post.Tags {
		tags = append(tags, createTagDocument(tag))
	}
	for _, category := range post.Categories {
		categories = append(categories, createCategoryDocument(category))
	}

	db.Collection(news.ColImages).InsertOne(context.Background(), image)
	db.Collection(news.ColVideos).InsertOne(context.Background(), video)
	if tags != nil {
		db.Collection(news.ColTags).InsertMany(context.Background(), tags)
	}
	if categories != nil {
		db.Collection(news.ColPostCategories).InsertMany(context.Background(), categories)
	}
	db.Collection(news.ColPosts).InsertOne(context.Background(), createPostDocument(post))

}

func migratePostReviewRecord(db *mongo.Database, review testReview) {
	db.Collection(news.ColReviews).InsertOne(context.Background(), createReviewDocument(review))
}

func metaOfPostResponse(p testPost) string {
	return fmt.Sprintf(`
	{
	"id": "%s",
	"style": "article:v2:default",
	"slug": "%s",
	"category_set": [{"category": null, "subcategory": null}],
	"hero_image": %s,
	"leading_image_portrait": %s,
	"og_image": %s,
	"og_description": "測試分享描述",
	"title": "測試標題",
	"subtitle": "測試副標",
	"published_date": "%s",
	"is_external": false,
	"tags": %s,
	"full": false,
	"bookmarkId": ""
	}
`,
		p.ID.Hex(),
		p.Slug,
		imageResponse(p.Image),
		imageResponse(p.Image),
		imageResponse(p.Image),
		p.CreatedAt.UTC().Format(time.RFC3339),
		tagsResponse(p.Tags...))
}

func tagsResponse(ids ...primitive.ObjectID) string {
	var s strings.Builder
	s.WriteString("[")
	for _, id := range ids {
		s.WriteString(tagResponse(id))
	}
	s.WriteString("]")
	return s.String()
}

func tagResponse(id primitive.ObjectID) string {
	return fmt.Sprintf(`
	{
	"id": "%s",
	"name": "測試標籤",
	"category": [],
	"key": "%s",
	"latest_order": 0
	}
`, id.Hex(), id.Hex())
}

func categoriesResponse(ids ...primitive.ObjectID) string {
	var s strings.Builder
	s.WriteString("[")
	for _, id := range ids {
		s.WriteString(categoryResponse(id))
	}
	s.WriteString("]")
	return s.String()
}

func categoryResponse(id primitive.ObjectID) string {
	return fmt.Sprintf(`
	{
	"id": "%s",
	"sort_order": 1,
	"name": "測試分類"
	}
`, id.Hex())
}

// createPostDocument builds the post document for testing
// to facilitate mocking document some fields are filled with same data
// updatedBy, createdBy: testPost.Editor
// updatedAt, createdAt, publishedDate: testPost.createdAt
// slug, mame,
func createPostDocument(p testPost) bson.M {
	return bson.M{
		"_id":          p.ID,
		"updatedBy":    p.Editor,
		"updatedAt":    p.CreatedAt,
		"createdBy":    p.Editor,
		"createdAt":    p.CreatedAt,
		"slug":         p.Slug,
		"name":         p.Slug,
		"toAutoNotify": true,
		"relateds":     p.Relateds,
		"tags":         p.Tags,
		"style":        "article:v2:default",
		"copyright":    "Copyrighted",
		"category_set": bson.M{
			"category":    nil,
			"subcategory": nil,
		},
		"heroImageSize": "normal",
		"engineers":     p.Engineers,
		"designers":     p.Designers,
		"photographers": p.Photographers,
		"writters":      p.Writers,
		"publishedDate": p.CreatedAt,
		"state":         p.State,
		"title":         "測試標題",
		"content": bson.M{
			"apiData": bson.A{
				bson.M{
					"styles":    bson.M{},
					"content":   bson.A{"測試本文"},
					"alignment": "center",
					"type":      "unstyled",
					"id":        "abcde",
				},
			},
			"draft": bson.M{
				"blocks": bson.A{
					bson.M{
						"data":              bson.M{},
						"entityRanges":      bson.A{},
						"inlineStyleRanges": bson.A{},
						"depth":             0,
						"type":              "unstyled",
						"text":              "測試本文",
						"key":               "abcde",
					},
				},
				"entityMap": bson.M{},
			},
			"html": "<p>測試本文</p>",
		},
		"extend_byline":             "其他使用者",
		"isFeatured":                false,
		"is_external":               false,
		"leading_image_description": "測試首圖描述",
		"og_description":            "測試分享描述",
		"og_title":                  "測試分享標題",
		"subtitle":                  "測試副標",
		"brief": bson.M{
			"apiData": bson.A{
				bson.M{
					"styles":    bson.M{},
					"content":   bson.A{"測試前言"},
					"alignment": "center",
					"type":      "unstyled",
					"id":        "abcde",
				},
			},
			"draft": bson.M{
				"blocks": bson.A{
					bson.M{
						"data":              bson.M{},
						"entityRanges":      bson.A{},
						"inlineStyleRanges": bson.A{},
						"depth":             0,
						"type":              "unstyled",
						"text":              "測試前言",
						"key":               "abcde",
					},
				},
				"entityMap": bson.M{},
			},
			"html": "<p>測試前言</p>",
		},
		"topics":                 p.Topic,
		"topics_ref":             p.Topic,
		"og_image":               p.Image,
		"heroImage":              p.Image,
		"leading_image_portrait": p.Image,
		"leading_video":          p.Video,
		"reviewWord":             p.ReviewWord,
	}
}

func createVideoDocument(id primitive.ObjectID) bson.M {
	return bson.M{
		"_id":   id,
		"title": "測試影片",
		"tags":  bson.A{},
		"video": bson.M{
			"filename":  "test.mp4",
			"filetype":  "video/mp4",
			"gcsBucket": "",
			"gcsDir":    "",
			"size":      100,
			"url":       "https://www.twreporter.org/test-video.mp4",
		},
	}
}

func createTagDocument(id primitive.ObjectID) bson.M {
	return bson.M{
		"_id":          id,
		"key":          id.Hex(),
		"name":         "測試標籤",
		"category":     bson.A{},
		"latest_order": 0,
	}
}

func createCategoryDocument(id primitive.ObjectID) bson.M {
	return bson.M{
		"_id":       id,
		"key":       id.Hex(),
		"sortOrder": 1,
		"name":      "測試分類",
	}
}

func createReviewDocument(r testReview) bson.M {
	return bson.M{
		"_id":     r.ID,
		"post_id": r.PostID,
		"order":   r.Order,
	}
}
