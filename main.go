package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/services"
	"twreporter.org/go-api/utils"
)

func main() {
	var err error
	var cf *controllers.ControllerFactory

	defer func() {
		if err != nil {
			log.Errorf("%+v", err)
		}
	}()

	globals.Conf, err = configs.LoadConf("")
	if err != nil {
		err = errors.Wrap(err, "Fatal error config file")
		return
	}

	// set up database connection
	log.Info("Connecting to MySQL cloud")
	db, err := utils.InitDB(10, 5)
	defer db.Close()
	if err != nil {
		return
	}

	log.Info("Connecting to MongoDB replica")
	session, err := utils.InitMongoDB()
	defer session.Close()
	if err != nil {
		return
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
		err = errors.Wrap(err, "Fail to start HTTP server")
	}
	return
}
