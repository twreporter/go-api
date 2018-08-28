package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/utils"

	//log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	var err error
	var cf *controllers.ControllerFactory

	viper.SetConfigType("json")
	viper.SetConfigFile("./configs/config.json")
	err = viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetDefault("consumersettings", map[string]string{
		"protocol": "http",
		"host":     "testtest.twreporter.org",
		"port":     "3000",
	})

	viper.SetDefault("appsettings", map[string]interface{}{
		"host":       "localhost",
		"port":       "8080",
		"protocol":   "http",
		"version":    "v1",
		"token":      "twreporter-token",
		"expiration": 168,
	})

	viper.SetDefault("mongodbsettings", map[string]interface{}{
		"url":     "locahost",
		"dbname":  "plate",
		"timeout": 5,
	})

	viper.SetDefault("dbsettings", map[string]string{
		"name":     "gorm",
		"user":     "gorm",
		"password": "gorm",
		"address":  "127.0.0.1",
		"port":     "3306",
	})

	viper.SetDefault("encryptsettings.salt", "salt")

	viper.SetDefault("mongodbsettings", map[string]interface{}{
		"url":     "locahost",
		"dbname":  "plate",
		"timeout": 5,
	})

	viper.SetDefault("encryptsettings.salt", "salt")

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// set up database connection
	log.Info("Connecting to MySQL cloud")
	db, err := utils.InitDB(10, 5)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	log.Info("Connecting to MongoDB replica")
	session, err := utils.InitMongoDB()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// mailSender := utils.NewSMTPEmailSender()                          // use office365 to send mails
	mailSender := utils.NewAmazonEmailSender() // use Amazon SES to send mails

	cf = controllers.NewControllerFactory(db, session, mailSender)

	// set up the router
	router := routers.SetupRouter(cf)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err = s.ListenAndServe(); err != nil {
		log.Error("Fail to start HTTP server", err.Error())
	}

}
