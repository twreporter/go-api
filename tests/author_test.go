package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/twreporter/go-api/globals"

	"github.com/twreporter/go-api/internal/news"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testAuthor struct {
	id, tid   primitive.ObjectID
	name      string
	createdAt time.Time
}

func TestGetAuthors_ByKeywords(t *testing.T) {
	defer cleanupAuthorRecords()
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
		"李明": {
			id:        primitive.NewObjectID(),
			tid:       primitive.NewObjectID(),
			name:      "李明",
			createdAt: time.Unix(1611818400, 0),
		},
	}
	// setup records
	for _, v := range authors {
		migrateAuthorRecord(v)
	}
	// TODO(babygoat): remove pre variable to overwrite Mongo.DBname after v1 endpoints & tests are removed
	pre := globals.Conf.DB.Mongo.DBname
	globals.Conf.DB.Mongo.DBname = testMongoDB

	response := serveHTTP(http.MethodGet, "/v2/authors?keywords=小明", "", "", "")
	// TODO(babygoat): remove pre variable to overwrite Mongo.DBname after v1 endpoints & tests are removed
	globals.Conf.DB.Mongo.DBname = pre

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, authorListResponse(authorResponse(authors["王小明"])), response.Body.String())
}

func TestGetAuthors_NoContent(t *testing.T) {
	response := serveHTTP(http.MethodGet, "/v2/authors", "", "", "")
	assert.Equal(t, http.StatusNoContent, response.Code)
}

func TestGetAuthorByID_ByValidID(t *testing.T) {
	defer cleanupAuthorRecords()

	author := testAuthor{
		id:        primitive.NewObjectID(),
		tid:       primitive.NewObjectID(),
		name:      "王小明",
		createdAt: time.Unix(1611817200, 0),
	}
	// setup records
	migrateAuthorRecord(author)
	// TODO(babygoat): remove pre variable to overwrite Mongo.DBname after v1 endpoints & tests are removed
	pre := globals.Conf.DB.Mongo.DBname
	globals.Conf.DB.Mongo.DBname = testMongoDB

	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/authors/%s", author.id.Hex()), "", "", "")
	// TODO(babygoat): remove pre variable to overwrite Mongo.DBname after v1 endpoints & tests are removed
	globals.Conf.DB.Mongo.DBname = pre

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, singleRecordResponse(authorResponse(author)), response.Body.String())
}

func cleanupAuthorRecords() {
	testMongoClient.Database(testMongoDB).Collection(news.ColContacts).Drop(context.Background())
	testMongoClient.Database(testMongoDB).Collection(news.ColImages).Drop(context.Background())
}

func migrateAuthorRecord(author testAuthor) {
	contact := createContactDocument(author.id, author.tid, author.name, author.createdAt)
	image := createImageDocument(author.tid)
	testMongoClient.Database(testMongoDB).Collection(news.ColContacts).InsertOne(context.Background(), contact)
	testMongoClient.Database(testMongoDB).Collection(news.ColImages).InsertOne(context.Background(), image)
}

func authorListResponse(authors ...string) string {
	return listResponse(len(authors), authors)
}

func authorResponse(author testAuthor) string { // use time as the id generation seed
	return fmt.Sprintf(`{
		"id":        "%s",
		"email":      "test@twreporter.org",
		"bio":        "test bio",
		"name":       "%s",
		"job_title":  "test job title",
		"thumbnail":  {
			"id":         "%s",
			"description": "test description",
			"filetype":  "image/jpeg",
			"resized_targets": {
				"tiny": {
					"height": 150,
					"width":  150,
					"url":    "https://www.twreporter.org/images/test-tiny.jpg"
				},
				"w400": {
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-w400.jpg"
				},
				"mobile": {
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-mobile.jpg"
				},
				"tablet": {
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-tablet.jpg"
				},
				"desktop": {
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-desktop.jpg"
				}
			}
		},
		"updated_at": "%s"
	}`, author.id.Hex(), author.name, author.tid.Hex(), author.createdAt.UTC().Format(time.RFC3339))

}

func listResponse(total int, records []string) string {
	return fmt.Sprintf(`{
	  "status": "success",
	  "data": {
		"meta": {
		  "offset": 0,
		  "limit": 10,
		  "total": %d
		},
		"records":[%s]
	  }
}`, total, strings.Join(records, ","))
}

func singleRecordResponse(record string) string {
	return fmt.Sprintf(`{
	"status": "success",
	"data": %s
}`, record)
}

func createContactDocument(id, thumbnailID primitive.ObjectID, name string, t time.Time) bson.M {
	return bson.M{
		"_id":        id,
		"email":      "test@twreporter.org",
		"bio":        bson.M{"html": "<p>test bio</p>", "md": "test bio"},
		"name":       name,
		"job_title":  "test job title",
		"thumbnail":  thumbnailID,
		"updated_at": t,
	}
}

func createImageDocument(id primitive.ObjectID) bson.M {
	return bson.M{
		"_id":         id,
		"description": "test description",
		"copyright":   "copyrighted",
		"keywords":    "keyword1, keyword2",
		"sale":        false,
		"image": bson.M{
			"filename":  "test name",
			"filetype":  "image/jpeg",
			"gcsBucket": "",
			"gcsDir":    "",
			"height":    400,
			"size":      160000,
			"width":     400,
			"resizedTargets": bson.M{
				"tiny": bson.M{
					"height": 150,
					"width":  150,
					"url":    "https://www.twreporter.org/images/test-tiny.jpg",
				},
				"w400": bson.M{
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-w400.jpg",
				},
				"mobile": bson.M{
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-mobile.jpg",
				},
				"tablet": bson.M{
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-tablet.jpg",
				},
				"desktop": bson.M{
					"height": 400,
					"width":  400,
					"url":    "https://www.twreporter.org/images/test-desktop.jpg",
				},
			},
		},
		"iptc": bson.M{
			"caption":      "test caption",
			"country":      "台灣",
			"country_code": "TWN",
			"byline":       "攝影師",
			"created_time": "",
			"created_date": "20210129",
			"keywords":     bson.A{"keyword1", "keyword2"},
			"city":         "taipei",
		},
	}
}
