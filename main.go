package main

import (
	"go/build"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	p, _ := build.Default.Import("twreporter.org/go-api", "", build.FindOnly)

	fname := filepath.Join(p.Dir, "configs/config.json")

	// Load config file
	err := utils.LoadConfig(fname)
	if err != nil {
		log.Fatal("main.load_config.fatal_error: ", err.Error())
	}

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

	cf, err := controllers.NewControllerFactory()

	if err != nil {
		panic(err)
	}

	defer cf.Close()

	// set up the router
	router := routers.SetupRouter(cf)
	router.Use(secureFunc)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	s.ListenAndServe()

}
