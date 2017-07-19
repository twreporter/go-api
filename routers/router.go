package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/utils"
)

// SetupRouter ...
func SetupRouter(cf *controllers.ControllerFactory) *gin.Engine {
	engine := gin.Default()

	if utils.Cfg.Environment == "production" {
		engine.Use(cors.New(cors.Config{
			AllowOrigins: []string{"https://twreporter.org", "https://dev.twreporter.org",
				"https://www.twreporter.org", "http://twreporter.org", "http://www.twreporter.org",
				"http://dev.twreporter.org"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "DELETE", "UPDATE"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "Accept",
				"Accept-Encoding", "Client-Security-Token", "X-Requested-With", "x-access-token"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           30 * time.Minute,
		}))
	}

	// Simple group: v1
	routerGroup := engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)

	return engine
}
