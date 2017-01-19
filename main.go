package main

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/middlewares"
)

func main() {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	router.Run() // listen and server on 0.0.0.0:8080
}
