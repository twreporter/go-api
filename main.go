package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/routers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cfg := configs.GetConfig()

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

	// connect to MySQL database
	db, err := gorm.Open("mysql", cfg.DB.User+":"+cfg.DB.Password+"@tcp("+cfg.DB.Address+":"+cfg.DB.Port+")/"+cfg.DB.Name)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	router := routers.SetupRouter()

	router.Use(secureFunc)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.ListenAndServe()
}
