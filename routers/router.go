package routers

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/middlewares"
)

// SetupRouter ...
func SetupRouter(cf *controllers.ControllerFactory) *gin.Engine {
	engine := gin.Default()
	engine.Use(middlewares.CORSMiddleware())

	// Simple group: v1
	routerGroup := engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)

	return engine
}
