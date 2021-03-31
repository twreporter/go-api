package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/twreporter/go-api/configs"
	"github.com/twreporter/go-api/controllers"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/internal/mongo"
	"github.com/twreporter/go-api/routers"
	"github.com/twreporter/go-api/services"
	"github.com/twreporter/go-api/utils"
)

func main() {
	var err error
	var cf *controllers.ControllerFactory

	defer func() {
		if err != nil {
			if globals.Conf.Environment == "development" {
				log.Errorf("%+v", err)
			} else {
				log.WithField("detail", err).Errorf("%s", f.FormatStack(err))
			}
		}
	}()

	globals.Conf, err = configs.LoadConf("")
	if err != nil {
		err = errors.Wrap(err, "Fatal error config file")
		return
	}

	configLogger()

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

	log.Info("Connection to MongoDB with mongo-go-driver")
	ctx := context.Background()
	opts := options.Client()
	client, err := mongo.NewClient(ctx, opts.ApplyURI(globals.Conf.DB.Mongo.URL).SetReadPreference(readpref.Nearest()))

	if err != nil {
		return
	}
	defer func() {
		client.Disconnect(ctx)
	}()
	// mailSender := services.NewSMTPMailService() // use office365 to send mails
	mailSvc := services.NewAmazonMailService() // use Amazon SES to send mails

	sClient := search.NewClient(globals.Conf.Algolia.ApplicationID, globals.Conf.Algolia.APIKey)
	cf = controllers.NewControllerFactory(db, session, mailSvc, client, sClient.InitIndex("contacts-index-v3"))

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

func configLogger() {
	env := globals.Conf.Environment
	switch env {
	// production/staging environments writes the log into standard output
	// and delegates log collector (fluentd) in k8s cluster to export to
	// stackdriver sink.
	case "production", "staging":
		log.SetOutput(os.Stdout)
		log.SetFormatter(f.NewStackdriverFormatter("go-api", env))
	// development environment reports the log location
	default:
		log.SetReportCaller(true)
	}
}
