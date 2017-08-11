package routers

import (
	// log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/utils"
)

// SetupRouter ...
func SetupRouter(cf *controllers.ControllerFactory) *gin.Engine {
	engine := gin.Default()

	config := cors.DefaultConfig()

	if utils.Cfg.Environment != "development" {
		if len(utils.Cfg.CorsSettings.AllowOrigins) > 0 {
			config.AllowOrigins = utils.Cfg.CorsSettings.AllowOrigins
		} else {
			config.AllowOrigins = []string{"https://www.twreporter.org"}
		}
		engine.Use(cors.New(config))
	}

	engine.Use(cors.New(config))

	routerGroup := engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)

	return engine
}
