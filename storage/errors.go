package storage

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
var ErrRecordNotFound = gorm.ErrRecordNotFound

// ErrMgoNotFound ...
var ErrMgoNotFound = mgo.ErrNotFound

// ErrDuplicateEntry record is already existed
var ErrDuplicateEntry uint16 = 1062

// ErrMgoDuplicateEntry ...
var ErrMgoDuplicateEntry = 11000
