package storage

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
	//log "github.com/Sirupsen/logrus"
)

// MembershipStorage defines the methods we need to implement,
// in order to fulfill the functionalities a membership system needs.
// Such as, let user signup, login w/o oauth, CRUD bookmarks, CRUD registrations.
type MembershipStorage interface {
	/** Close DB Connection **/
	Close() error

	/** User methods **/
	GetUserByID(string) (models.User, error)
	GetOAuthData(sql.NullString, string) (models.OAuthAccount, error)
	GetUserDataByOAuth(models.OAuthAccount) (models.User, error)
	GetReporterAccountData(string) (*models.ReporterAccount, error)
	GetUserDataByReporterAccount(*models.ReporterAccount) (*models.User, error)
	InsertUserByOAuth(models.OAuthAccount) models.User
	InsertUserByReporterAccount(models.ReporterAccount) (models.User, error)
	UpdateOAuthData(models.OAuthAccount) (models.OAuthAccount, error)
	UpdateReporterAccountPassword(*models.ReporterAccount, string) (*models.ReporterAccount, error)
	UpdateReporterAccountActive(*models.ReporterAccount, bool) (*models.ReporterAccount, error)

	/** Bookmark methods **/
	GetABookmarkByHref(string) (models.Bookmark, error)
	GetABookmarkByID(string) (models.Bookmark, error)
	GetBookmarksOfAUser(string) ([]models.Bookmark, error)
	CreateABookmarkOfAUser(string, models.Bookmark) error
	DeleteABookmarkOfAUser(string, string) error

	/** Service methods **/
	GetService(string) (models.Service, error)
	CreateService(models.ServiceJSON) (models.Service, error)
	UpdateService(string, models.ServiceJSON) (models.Service, error)
	DeleteService(string) error

	/** Registration methods **/
	GetRegistration(string, string) (models.Registration, error)
	GetRegistrationsByService(string, int, int, string, int) ([]models.Registration, error)
	GetRegistrationsAmountByService(string, int) (uint, error)
	CreateRegistration(string, models.RegistrationJSON) (models.Registration, error)
	UpdateRegistration(string, models.RegistrationJSON) (models.Registration, error)
	DeleteRegistration(string, string) error
}

// NewGormStorage initializes the storage connected to MySQL database by gorm library
func NewGormStorage(db *gorm.DB) *GormStorage {
	return &GormStorage{db}
}

// GormStorage implements MembershipStorage interface
type GormStorage struct {
	db *gorm.DB
}

// Close quits the DB connection gracefully
func (gs *GormStorage) Close() error {
	err := gs.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// NewsStorage defines the methods we need to implement,
// in order to provide the news resource to twreporter main site.
type NewsStorage interface {
	/** Close DB Connection **/
	Close() error

	/** Posts methods **/
	GetMetaOfPosts(string, int, int, string) ([]models.PostMeta, error)
	// GetTopics(string) (models.PostMeta, error)
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
func (ms *MongoStorage) Close() error {
	ms.db.Close()
	return nil
}

func (m *MongoStorage) GetDocuments(qs interface{}, limit int, offset int, sort string, collection string, documents interface{}) error {
	var err error
	var q models.MongoQuery

	_qs, ok := qs.(string)

	if ok {
		err = models.GetQuery(_qs, &q)

		if err != nil {
			return m.NewStorageError(err, "GetDocuments", "storage.mongo_storage.get_documents.parse_query_error")
		}

		qs = q
	}

	err = m.db.DB(utils.Cfg.MongoDBSettings.DBName).C(collection).Find(qs).Limit(limit).Skip(offset).Sort(sort).All(documents)

	if err != nil {
		return m.NewStorageError(err, "GetDocuments", "storage.mongo_storage.get_documents_error")
	}

	return nil
}

func (m *MongoStorage) GetDocument(id bson.ObjectId, collection string, doc interface{}) error {
	if id == "" {
		return m.NewStorageError(ErrMgoNotFound, "GetTopicMeta", "storage.mongo_storage.get_document.id_not_provided")
	}

	err := m.db.DB(utils.Cfg.MongoDBSettings.DBName).C(collection).FindId(id).One(doc)

	if err != nil {
		return m.NewStorageError(err, "GetDocument", "storage.mongo_storage.get_document.error")
	}
	return nil
}
