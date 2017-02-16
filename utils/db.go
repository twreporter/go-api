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
	db, err := gorm.Open("mysql", cfg.DB.User+":"+cfg.DB.Password+"@tcp("+cfg.DB.Address+":"+cfg.DB.Port+")/"+cfg.DB.Name)

	db.AutoMigrate(&models.User{}, &models.OAuthAccount{})

	if err != nil {
		log.Error("Please check the MySQL database connection.")
		return nil, err
	}

	return db, nil
}
