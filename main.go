package main

import (
	"go/build"
	"net/http"
	"path/filepath"
	"time"

	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	var err error
	var cf *controllers.ControllerFactory

	p, _ := build.Default.Import("twreporter.org/go-api", "", build.FindOnly)

	fname := filepath.Join(p.Dir, "configs/config.json")

	// Load config file
	err = utils.LoadConfig(fname)
	if err != nil {
		log.Fatal("main.load_config.fatal_error: ", err.Error())
	}

	// set up database connection
	log.Info("Connecting to MySQL cloud")
	db, err := utils.InitDB(10, 5)
	if err != nil {
		panic(err)
	}

	log.Info("Connecting to MongoDB replica")
	session, err := utils.InitMongoDB()
	if err != nil {
		panic(err)
	}

	// mailSender := utils.NewSMTPEmailSender()                          // use office365 to send mails
	mailSender := utils.NewAmazonEmailSender() // use Amazon SES to send mails

	cf = controllers.NewControllerFactory(db, session, mailSender)

	defer db.Close()
	defer session.Close()

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
