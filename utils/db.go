package utils

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/globalsign/mgo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/matryer/try.v1"

	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
)

// InitDB initiates the MySQL database connection
func InitDB(attempts, retryMaxDelay int) (*gorm.DB, error) {
	var db *gorm.DB
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		var config = globals.Conf.DB.MySQL

		// connect to MySQL database
		var endpoint = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=true", config.User, config.Password, config.Address, config.Port, config.Name)
		log.Debug("connect to mysql ", endpoint)
		db, err = gorm.Open("mysql", endpoint)

		if err != nil {
			time.Sleep(time.Duration(retryMaxDelay) * time.Second)
		}

		return attempt < attempts, errors.WithStack(err)
	})

	if err != nil {
		return nil, errors.Wrap(err, "Please check the MySQL database connection: ")
	}

	db.SetJoinTableHandler(&models.User{}, globals.TableBookmarks, &models.UsersBookmarks{})

	//db.LogMode(true)

	return db, nil
}

// InitMongoDB initiates the Mongo DB connection
func InitMongoDB() (*mgo.Session, error) {
	var timeout = globals.Conf.DB.Mongo.Timeout
	// Set connection timeout
	session, err := mgo.DialWithTimeout(globals.Conf.DB.Mongo.URL, time.Duration(timeout)*time.Second)
	log.Debug("connect to mongodb ", globals.Conf.DB.Mongo.URL)

	if err != nil {
		return nil, errors.Wrap(err, "Establishing a new session to the mongo occurs error: ")
	}

	// Set operation timeout
	session.SetSyncTimeout(time.Duration(timeout) * time.Second)

	// Set socket timeout to 3 mins
	session.SetSocketTimeout(3 * time.Minute)

	// As our mongo cluster comprises cost-effective solution(Replica set arbiter),
	// use Nearest read concern(https://docs.mongodb.com/manual/core/read-preference/#nearest)
	// to distribute the read load acorss primary and secondary node evenly.
	session.SetMode(mgo.Nearest, true)

	return session, nil
}

func InitMongoDBV2() (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(globals.Conf.DB.Mongo.URL)
	clientOpts = clientOpts.SetReadPreference(readpref.Nearest())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, errors.Wrap(err, "Establishing a new connection to cluster occurs error:")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "Connection to cluster does not response:")
	}
	return client, nil
}

// Get the migrate instance for operating migration
func GetMigrateInstance(dbInstance *sql.DB) (*migrate.Migrate, error) {
	const migrateMysqlDriver = "mysql"
	const migrateSourceDriver = "file"
	var migrateSourceDir string = filepath.Join(GetProjectRoot(), "migrations")

	driver, _ := mysql.WithInstance(dbInstance, &mysql.Config{})

	sourceUrl := fmt.Sprintf("%s://%s", migrateSourceDriver, migrateSourceDir)
	m, err := migrate.NewWithDatabaseInstance(sourceUrl, migrateMysqlDriver, driver)

	return m, errors.WithStack(err)
}
