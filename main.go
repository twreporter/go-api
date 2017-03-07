package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	// Load config file
	utils.LoadConfig("config.json")

	// security: no one can put it in an iframe
	secureMiddleware := secure.New(secure.Options{
		FrameDeny: true,
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	// set up database connection
	db, _ := utils.InitDB()
	defer db.Close()

	// set up data storage
	// userStorage := storage.NewUserStorage(db)
	userStorage := storage.NewUserStorage(db)

	mailSender := utils.NewSMTPEmailSender(utils.Cfg.EmailSettings)

	// set up the router
	router := routers.SetupRouter(userStorage, mailSender)

	router.Use(secureFunc)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.ListenAndServe()
}
