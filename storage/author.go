package storage

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

// GetFullAuthors finds the authors according to mongo aggregation pipeline stages
func (m *MongoStorage) GetFullAuthors(limit int, offset int, sort string) ([]models.FullAuthor, int, error) {
	type author struct {
		ID       bson.ObjectId `bson:"_id"`
		JobTitle string        `bson:"job_title"`
		Name     string        `bson:"name"`
		Bio      struct {
			Html string `bson:"html"`
			Md   string `bson:"md"`
		}
		Email      string              `bson:"email"`
		Thumbnails []models.MongoImage `bson:"thumbnails"`
		UpdatedAt  time.Time           `bson:"updatedAt"`
	}

	var authors []models.FullAuthor
	var fa models.FullAuthor
	var total int
	var err error
	var results []author
	var result author

	pipeline := []bson.M{
		bson.M{"$sort": bson.M{sort: -1}},
		bson.M{"$skip": offset},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{"from": "images", "localField": "image", "foreignField": "_id", "as": "thumbnails"}},
	}

	collection := m.db.DB(globals.Conf.DB.Mongo.DBname).C("contacts")
	if total, err = collection.Count(); err != nil {
		return authors, 0, errors.Wrap(err, "can not get total count of authors")
	}

	pipe := collection.Pipe(pipeline)

	if err = pipe.All(&results); err != nil {
		return authors, 0, errors.Wrap(err, "can not get authors from storage")
	}

	// Copy fields/values from `author`s to `FullAuthor`s
	for i := range results {
		result = results[i]
		fa = models.FullAuthor{}
		copier.Copy(&fa, &result)
		if len(result.Thumbnails) > 0 {
			fa.Thumbnail = &result.Thumbnails[0]
		}
		authors = append(authors, fa)
	}

	return authors, total, nil
}
