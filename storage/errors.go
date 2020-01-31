package storage

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
var ErrRecordNotFound = gorm.ErrRecordNotFound

// ErrMgoNotFound record not found error when accessing MongoDB
var ErrMgoNotFound = mgo.ErrNotFound

func IsNotFound(err error) bool {
	cause := errors.Cause(err)

	switch cause {
	case ErrRecordNotFound:
		return true
	case ErrMgoNotFound:
		return true
	default:
		// omit intentionally
	}
	return false
}

func IsConflict(err error) bool {
	// ErrDuplicateEntry record is already existed in MySQL database
	var ErrDuplicateEntry uint16 = 1062
	// ErrMgoDuplicateEntry record is already existed in MongoDB
	var ErrMgoDuplicateEntry = 11000

	cause := errors.Cause(err)

	switch e := cause.(type) {
	case *mysql.MySQLError:
		return e.Number == ErrDuplicateEntry
	case *mgo.LastError:
		return e.Code == ErrMgoDuplicateEntry
	default:
		// omit intentionally
	}
	return false
}
