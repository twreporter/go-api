package routers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

type wrappedFn func(c *gin.Context) (int, gin.H, error)

func ginResponseWrapper(fn wrappedFn) func(c *gin.Context) {
	return func(c *gin.Context) {
		statusCode, obj, err := fn(c)
		if err != nil {
			appErr := err.(*models.AppError)
			log.Error(appErr.Error())
			c.JSON(appErr.StatusCode, gin.H{"status": "error", "message": appErr.Message})
			return
		}
		c.JSON(statusCode, obj)
	}
}

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
	} else {
		config.AllowAllOrigins = true
	}

	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("DELETE")

	engine.Use(cors.New(config))

	v1Group := engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		v1Group.GET("/ping", menuitems.Retrieve)
	}

	// =============================
	// membership service endpoints
	// =============================
	mc := cf.GetMembershipController()
	// endpoints for account
	v1Group.POST("/signin", middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.SignIn(c, cf.GetMailSender())
	}))
	v1Group.GET("/activate", middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.Activate))
	v1Group.GET("/token/:userID", middlewares.CheckJWT(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.RenewJWT))
	// endpoints for bookmarks of users
	v1Group.GET("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	v1Group.GET("/users/:userID/bookmarks/:bookmarkSlug", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	v1Group.POST("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.CreateABookmarkOfAUser))
	v1Group.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.DeleteABookmarkOfAUser))
	// endpoints for web push subscriptions
	v1Group.POST("/web-push/subscriptions" /*middlewares.CheckJWT()*/, ginResponseWrapper(mc.SubscribeWebPush))
	v1Group.GET("/web-push/subscriptions", ginResponseWrapper(mc.IsWebPushSubscribed))

	// =============================
	// news service endpoints
	// =============================
	nc := cf.GetNewsController()
	// endpoints for authors
	v1Group.GET("/authors", middlewares.SetCacheControl("public,max-age=600"), ginResponseWrapper(nc.GetAuthors))
	// endpoints for posts
	v1Group.GET("/posts", middlewares.SetCacheControl("public,max-age=900"), nc.GetPosts)
	v1Group.GET("/posts/:slug", middlewares.SetCacheControl("public,max-age=900"), nc.GetAPost)
	// endpoints for topics
	v1Group.GET("/topics", middlewares.SetCacheControl("public,max-age=900"), nc.GetTopics)
	v1Group.GET("/topics/:slug", middlewares.SetCacheControl("public,max-age=900"), nc.GetATopic)
	v1Group.GET("/index_page", middlewares.SetCacheControl("public,max-age=1800"), nc.GetIndexPageContents)
	v1Group.GET("/index_page_categories", middlewares.SetCacheControl("public,max-age=1800"), nc.GetCategoriesPosts)
	// endpoints for search
	v1Group.GET("/search/authors", middlewares.SetCacheControl("public,max-age=3600"), nc.SearchAuthors)
	v1Group.GET("/search/posts", middlewares.SetCacheControl("public,max-age=3600"), nc.SearchPosts)

	// =============================
	// oauth endpoints
	// =============================
	authGroup := v1Group.Group("/auth")

	gc := cf.GetGoogleController()
	authGroup.GET("/google", middlewares.SetCacheControl("no-store"), gc.BeginAuth)
	authGroup.GET("/google/callback", middlewares.SetCacheControl("no-store"), gc.Authenticate)
	fc := cf.GetFacebookController()
	authGroup.GET("/facebook", middlewares.SetCacheControl("no-store"), fc.BeginAuth)
	authGroup.GET("/facebook/callback", middlewares.SetCacheControl("no-store"), fc.Authenticate)

	return engine
}
