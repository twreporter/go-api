package main

import (
	"net/http"
	"time"

	"twreporter.org/go-api/routers"
)

func main() {
	router := routers.SetupRouter()

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	s.ListenAndServe()
}
