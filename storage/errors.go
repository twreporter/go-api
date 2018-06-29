package storage

import (
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
var ErrRecordNotFound = gorm.ErrRecordNotFound

// ErrMgoNotFound record not found error when accessing MongoDB
var ErrMgoNotFound = mgo.ErrNotFound

// ErrDuplicateEntry record is already existed in MySQL database
var ErrDuplicateEntry uint16 = 1062

// ErrMgoDuplicateEntry record is already existed in MongoDB
var ErrMgoDuplicateEntry = 11000

// IsRecordNotFoundError check if err is equal to gorm.ErrRecordNotFound
func IsRecordNotFoundError(err error) bool {
	return err == ErrRecordNotFound
}

// IsDuplicateEntryError check if err belongs to mysql.MySQLError and its error number is equal to ErrDuplicateEntry
func IsDuplicateEntryError(err error) bool {
	errStruct, ok := err.(*mysql.MySQLError)
	if ok && errStruct.Number == ErrDuplicateEntry {
		return true
	}
	return false
}

// NewStorageError return AppError with detailed information.
// This method is mainly used to deal with MySQL error.
func (g *GormStorage) NewStorageError(err error, where string, message string) error {
	switch {
	case IsRecordNotFoundError(err):
		return models.NewAppError(where, "Record not found", fmt.Sprintf("%v : %v", message, err.Error()), http.StatusNotFound)
	case IsDuplicateEntryError(err):
		return models.NewAppError(where, "Record is already existed", fmt.Sprintf("%v : %v", message, err.Error()), http.StatusConflict)
	case err != nil:
		log.Error(err.Error())
		return models.NewAppError(where, "Internal server error", fmt.Sprintf("%v : %v", message, err.Error()), http.StatusInternalServerError)
	default:
		return nil
	}
}

// NewStorageError return AppError with detailed information.
// This method is mainly used to deal with MongoDB error.
func (m *MongoStorage) NewStorageError(err error, where string, message string) error {
	errStruct, ok := err.(*mgo.LastError)

	if err != nil && err.Error() == ErrMgoNotFound.Error() {
		return models.NewAppError(where, "Record not found", fmt.Sprintf("%v : %v", message, err.Error()), http.StatusNotFound)
	} else if ok && errStruct.Code == ErrMgoDuplicateEntry {
		return models.NewAppError(where, "Record is already existed", fmt.Sprintf("%v : %v", message, err.Error()), http.StatusConflict)
	} else if err != nil {
		log.Error(err.Error())
		return models.NewAppError(where, "Internal server error", fmt.Sprintf("%v : %v", message, err.Error()), http.StatusInternalServerError)
	}

	return nil
}
