package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/middlewares"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		v1.GET("/ping", menuitems.Retrieve)
	}

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.ListenAndServe()
}
