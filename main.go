package main

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/middlewares"
)

func main() {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	r.Run() // listen and server on 0.0.0.0:8080
}
