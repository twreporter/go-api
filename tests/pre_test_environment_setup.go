package tests

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

const (
	mgoDBName = "mgo"

	// collections name
	mgoPostCol       = "posts"
	mgoTopicCol      = "topics"
	mgoImgCol        = "images"
	mgoVideoCol      = "videos"
	mgoTagCol        = "tags"
	mgoCategoriesCol = "postcategories"
	mgoThemeCol      = "themes"
)

func runMySQLMigration(gormDB *gorm.DB) {
	// Make sure there is no existing schemas
	dropStmt := []byte(`
		SET FOREIGN_KEY_CHECKS = 0;
		SET GROUP_CONCAT_MAX_LEN=32768;
		SET @tables = NULL;
		SELECT GROUP_CONCAT(table_schema, '.', table_name) INTO @tables
		FROM information_schema.tables
		WHERE table_schema = (SELECT DATABASE());
		SELECT IFNULL(@tables,"dummy") INTO @tables;

		SET @tables = CONCAT("DROP TABLE IF EXISTS ", @tables);
		PREPARE stmt FROM @tables;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
		SET FOREIGN_KEY_CHECKS = 1;
	`)
	gormDB.Exec(string(dropStmt))

	m, err := utils.GetMigrateInstance(gormDB.DB())

	if nil != err {
		panic("unable to get migrate instance: " + err.Error())
	}

	// Upgrade to the latest migrations
	if err = m.Up(); err != nil {
		panic("unable to migrate mysql: " + err.Error())
	}

}

func runMgoMigration(mgoDB *mgo.Session) {
	err := mgoDB.DB(mgoDBName).DropDatabase()
	if err != nil {
		panic(fmt.Sprint("Can not drop mongo gorm database"))
	}

	mgoDB.DB(mgoDBName).Run(bson.D{{Name: "create", Value: mgoPostCol}}, nil)
	mgoDB.DB(mgoDBName).Run(bson.D{{Name: "create", Value: mgoTopicCol}}, nil)
	mgoDB.DB(mgoDBName).Run(bson.D{{Name: "create", Value: mgoImgCol}}, nil)
	mgoDB.DB(mgoDBName).Run(bson.D{{Name: "create", Value: mgoVideoCol}}, nil)
	mgoDB.DB(mgoDBName).Run(bson.D{{Name: "create", Value: mgoTagCol}}, nil)
	mgoDB.DB(mgoDBName).Run(bson.D{{Name: "create", Value: mgoCategoriesCol}}, nil)
}

func setGormDefaultRecords(gormDB *gorm.DB) {
	// Set an active reporter account
	ms := storage.NewGormStorage(gormDB)

	ra := models.ReporterAccount{
		Email:         Globs.Defaults.Account,
		ActivateToken: Globs.Defaults.Token,
		ActExpTime:    time.Now().Add(time.Duration(15) * time.Minute),
	}
	_, _ = ms.InsertUserByReporterAccount(ra)

	ms.CreateService(models.ServiceJSON{Name: Globs.Defaults.Service})

	ms.CreateRegistration(Globs.Defaults.Service, models.RegistrationJSON{Email: Globs.Defaults.Account, ActivateToken: Globs.Defaults.Token})

}

func setMgoDefaultRecords(mgoDB *mgo.Session) {
	// insert img1 and  img2
	mgoDB.DB(mgoDBName).C(mgoImgCol).Insert(Globs.Defaults.ImgCol1)

	mgoDB.DB(mgoDBName).C(mgoImgCol).Insert(Globs.Defaults.ImgCol2)
	// insert video
	mgoDB.DB(mgoDBName).C(mgoVideoCol).Insert(Globs.Defaults.VideoCol)

	// insert tag and postcategory
	mgoDB.DB(mgoDBName).C(mgoTagCol).Insert(Globs.Defaults.TagCol)
	mgoDB.DB(mgoDBName).C(mgoCategoriesCol).Insert(Globs.Defaults.CatCol)

	// insert post1 and post2
	mgoDB.DB(mgoDBName).C(mgoPostCol).Insert(Globs.Defaults.PostCol1)

	// insert theme
	mgoDB.DB(mgoDBName).C(mgoThemeCol).Insert(Globs.Defaults.ThemeCol)

	mgoDB.DB(mgoDBName).C(mgoPostCol).Insert(Globs.Defaults.PostCol2)

	// insert topic
	mgoDB.DB(mgoDBName).C(mgoTopicCol).Insert(Globs.Defaults.TopicCol)
}

func openGormConnection() (db *gorm.DB, err error) {
	// CREATE USER 'gorm'@'localhost' IDENTIFIED BY 'gorm';
	// CREATE DATABASE gorm;
	// GRANT ALL ON gorm.* TO 'gorm'@'localhost';
	dbhost := os.Getenv("GORM_DBADDRESS")
	if dbhost != "" {
		dbhost = fmt.Sprintf("tcp(%v)", dbhost)
	} else {
		dbhost = "tcp(127.0.0.1:3306)"
	}
	db, err = gorm.Open("mysql", fmt.Sprintf("gorm:gorm@%v/gorm?charset=utf8mb4,utf8&parseTime=True&multiStatements=true", dbhost))

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(10)

	return
}

func openMgoConnection() (session *mgo.Session, err error) {
	dbhost := os.Getenv("MGO_DBADDRESS")
	if dbhost == "" {
		dbhost = "localhost"
	}
	session, err = mgo.Dial(dbhost)

	// set settings
	globals.Conf.DB.Mongo.DBname = mgoDBName

	return
}

func setUpDBEnvironment() (*gorm.DB, *mgo.Session) {
	var err error
	var gormDB *gorm.DB
	var mgoDB *mgo.Session

	// Create DB connections
	if gormDB, err = openGormConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to test database, but got err=%+v", err))
	}

	gormDB.SetJoinTableHandler(&models.User{}, globals.TableBookmarks, &models.UsersBookmarks{})

	// Create Mongo DB connections
	if mgoDB, err = openMgoConnection(); err != nil {
		panic(fmt.Sprintf("No error should happen when connecting to mongo database, but got err=%+v", err))
	}

	// set up tables in gorm DB
	runMySQLMigration(gormDB)

	// set up default records in gorm DB
	setGormDefaultRecords(gormDB)

	// set up collections in mongoDB
	runMgoMigration(mgoDB)

	// set up default collections in mongoDB
	setMgoDefaultRecords(mgoDB)

	return gormDB, mgoDB
}
