package main

import (
	"fmt"
	"net/http"
	"time"

	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/services"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	var err error
	var cf *controllers.ControllerFactory

	globals.Conf, err = configs.LoadConf("")
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

	// mailSender := services.NewSMTPMailService() // use office365 to send mails
	mailSvc := services.NewAmazonMailService() // use Amazon SES to send mails

	cf = controllers.NewControllerFactory(db, session, mailSvc)

	// set up the router
	router := routers.SetupRouter(cf)

	readTimeout := 5 * time.Second

	// Set writeTimeout bigger than 30 secs.
	// 30 secs is to ensure donation request is handled correctly.
	writeTimeout := 40 * time.Second
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", globals.LocalhostPort),
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	if err = s.ListenAndServe(); err != nil {
		log.Error("Fail to start HTTP server", err.Error())
	}

}
