package utils

import (
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/matryer/try.v1"
	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
)

// InitDB initiates the MySQL database connection
func InitDB(attempts, retryMaxDelay int) (*gorm.DB, error) {
	var db *gorm.DB
	err := try.Do(func(attempt int) (bool, error) {
		var err error

		// connect to MySQL database
		db, err = gorm.Open("mysql", Cfg.DBSettings.User+":"+Cfg.DBSettings.Password+"@tcp("+Cfg.DBSettings.Address+":"+Cfg.DBSettings.Port+")/"+Cfg.DBSettings.Name+"?parseTime=true")

		if err != nil {
			time.Sleep(time.Duration(retryMaxDelay) * time.Second)
		}

		return attempt < attempts, err
	})

	if err != nil {
		log.Error("Please check the MySQL database connection: ", err.Error())
		return nil, err
	}

	// db.LogMode(true)

	return db, nil
}

// InitMongoDB initiates the Mongo DB connection
func InitMongoDB() (*mgo.Session, error) {
	// Set connection timeout
	session, err := mgo.DialWithTimeout(Cfg.MongoDBSettings.URL, time.Duration(Cfg.MongoDBSettings.Timeout)*time.Second)

	if err != nil {
		log.Error("Establishing a new session to the mongo occurs error: ", err.Error())
		return nil, err
	}

	// Set operation timeout
	session.SetSyncTimeout(time.Duration(Cfg.MongoDBSettings.Timeout) * time.Second)

	// Set socket timeout to 3 mins
	session.SetSocketTimeout(3 * time.Minute)

	// Since we don't have writes here and don't care about the consistency between Mongo Master and Slave,
	// we choose Eventual mode here.
	// The Eventual mode is the fastest and most resource-friendly,
	// but is also the one offering the least guarantees about ordering of the data read and written.
	session.SetMode(mgo.Eventual, true)

	return session, nil
}
