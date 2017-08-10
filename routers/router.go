package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/utils"
)

// SetupRouter ...
func SetupRouter(cf *controllers.ControllerFactory) *gin.Engine {
	engine := gin.Default()

	if utils.Cfg.Environment == "production" {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"https://twreporter.org", "https://dev.twreporter.org",
			"https://www.twreporter.org", "http://twreporter.org", "http://www.twreporter.org",
			"http://dev.twreporter.org", "http://staging.twreporter.org", "https://staging.twreporter.org"}
		engine.Use(cors.New(config))
		engine.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "https://www.twreporter.org")
		})
	} else {
		// TODO: use cors.Default() after new version of github.com/gin-contrib/cors
		engine.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		})
	}

	routerGroup := engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)

	return engine
}
