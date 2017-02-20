package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/models"
)

// InitDB initiates the database connection
func InitDB() (*gorm.DB, error) {
	cfg := configs.GetConfig()

	// connect to MySQL database
	db, err := gorm.Open("mysql", cfg.DB.User+":"+cfg.DB.Password+"@tcp("+cfg.DB.Address+":"+cfg.DB.Port+")/"+cfg.DB.Name+"?parseTime=true")

	if err != nil {
		log.Error("Please check the MySQL database connection.")
		return nil, err
	}

	// automatically migrate the schema, to keep them update to date.
	db.AutoMigrate(&models.User{}, &models.OAuthAccount{}, &models.ReporterAccount{})

	return db, nil
}
