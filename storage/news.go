package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
)

// NewsStorage defines the methods we need to implement,
// in order to provide the news resource to twreporter main site.
type NewsStorage interface {
	/** Close DB Connection **/
	Close() error

	/** Posts methods **/
	GetMetaOfPosts(models.MongoQuery, int, int, string, []string) ([]models.Post, int, error)
	GetFullPosts(models.MongoQuery, int, int, string, []string) ([]models.Post, int, error)
	GetMetaOfTopics(models.MongoQuery, int, int, string, []string) ([]models.Topic, int, error)
	GetFullTopics(models.MongoQuery, int, int, string, []string) ([]models.Topic, int, error)

	/** Authors methods **/
	GetFullAuthors(int, int, string) ([]models.FullAuthor, int, error)
}

// NewMongoStorage initializes the storage connected to Mongo database
func NewMongoStorage(db *mgo.Session) *MongoStorage {
	return &MongoStorage{db}
}

// MongoStorage implements `NewsStorage`
type MongoStorage struct {
	db *mgo.Session
}

// Close quits the DB connection gracefully
func (m *MongoStorage) Close() error {
	m.db.Close()
	return nil
}

// GetDocuments ...
func (m *MongoStorage) GetDocuments(qs models.MongoQuery, limit int, offset int, sort string, collection string, documents interface{}) (count int, err error) {
	var dbname = globals.Conf.DB.Mongo.DBname

	session := m.db.Copy()
	defer session.Close()

	err = session.DB(dbname).C(collection).Find(qs).Limit(limit).Skip(offset).Sort(sort).All(documents)

	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("get documents by conditions(where: %#v, limit: %d, offset: %d, sort: %s, collection:%s) occurs error", qs, limit, offset, sort, collection))
	}

	count, err = session.DB(dbname).C(collection).Find(qs).Count()

	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("count documents by condition(where: %#v, collection: %s) occurs error", qs, collection))
	}

	return count, nil
}

// GetDocument ...
func (m *MongoStorage) GetDocument(id bson.ObjectId, collection string, doc interface{}) error {
	if id == "" {
		return errors.Wrap(ErrMgoNotFound, "can not get document by zeroed string")
	}
	session := m.db.Copy()
	defer session.Close()

	err := session.DB(globals.Conf.DB.Mongo.DBname).C(collection).FindId(id).One(doc)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("get document(id: %v, collection: %s) occurs error", id, collection))
	}
	return nil
}
