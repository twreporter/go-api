package routers

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/controllers/oauth"
	"twreporter.org/go-api/middlewares"
)

// SetupRouter ...
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		v1.GET("/ping", menuitems.Retrieve)
		// handle login
		oauth := new(oauth.Facebook)
		v1.GET("/auth/facebook", oauth.BeginAuth)
		v1.GET("/auth/facebook/callback", oauth.Authenticate)
	}

	return router
}
