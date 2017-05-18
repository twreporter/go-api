package tests

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
)

var (
	DefaultID        = "1"
	DefaultAccount   = "nickhsine@twreporter.org"
	DefaultPassword  = "0000"
	DefaultID2       = "2"
	DefaultAccount2  = "hsunpei_wang@twreporter.org"
	DefaultPassword2 = "1111"
)

func OpenTestConnection() (db *gorm.DB, err error) {
	// CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm';
	// CREATE DATABASE gorm;
	// GRANT ALL ON gorm.* TO 'gorm'@'localhost';
	fmt.Println("testing mysql...")
	dbhost := os.Getenv("GORM_DBADDRESS")
	if dbhost != "" {
		dbhost = fmt.Sprintf("tcp(%v)", dbhost)
	} else {
		dbhost = "tcp(127.0.0.1:3306)"
	}
	db, err = gorm.Open("mysql", fmt.Sprintf("gorm:gorm@%v/gorm?charset=utf8&parseTime=True", dbhost))

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(10)

	return
}

func RunMigration(db *gorm.DB) {
	for _, table := range []string{"users_bookmarks"} {
		db.Exec(fmt.Sprintf("drop table %v;", table))
	}

	values := []interface{}{&models.User{}, &models.OAuthAccount{}, &models.ReporterAccount{}, &models.Bookmark{}, &models.Registration{}, &models.Service{}}
	for _, value := range values {
		db.DropTable(value)
	}
	if err := db.AutoMigrate(values...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}
}

func SetDefaultRecords(db *gorm.DB) {
	// Set an active reporter account
	as := storage.NewMembershipStorage(db)

	key, _ := scrypt.Key([]byte(DefaultPassword), []byte(""), 16384, 8, 1, 32)

	ra := models.ReporterAccount{
		Account:       DefaultAccount,
		Password:      fmt.Sprintf("%x", key),
		Active:        true,
		ActivateToken: "",
	}
	_, _ = as.InsertUserByReporterAccount(ra)

	key, _ = scrypt.Key([]byte(DefaultPassword2), []byte(""), 16384, 8, 1, 32)

	ra = models.ReporterAccount{
		Account:       DefaultAccount2,
		Password:      fmt.Sprintf("%x", key),
		Active:        true,
		ActivateToken: "",
	}
	_, _ = as.InsertUserByReporterAccount(ra)
}

func RequestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
