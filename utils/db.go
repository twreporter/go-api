package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
)

// InitDB initiates the database connection
func InitDB() (*gorm.DB, error) {
	// connect to MySQL database
	db, err := gorm.Open("mysql", Cfg.DBSettings.User+":"+Cfg.DBSettings.Password+"@tcp("+Cfg.DBSettings.Address+":"+Cfg.DBSettings.Port+")/"+Cfg.DBSettings.Name+"?parseTime=true")

	if err != nil {
		log.Error("Please check the MySQL database connection: ", err.Error())
		return nil, err
	}

	// automatically migrate the schema, to keep them update to date.
	db.AutoMigrate(&models.User{}, &models.OAuthAccount{}, &models.ReporterAccount{}, &models.Bookmark{})

	return db, nil
}
