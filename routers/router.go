package routers

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/controllers/oauth/facebook"
	"twreporter.org/go-api/controllers/oauth/google"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

// SetupRouter ...
func SetupRouter(userStorage storage.UserStorage, bookmarkStorage storage.BookmarkStorage, mailSender utils.EmailSender) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		v1.GET("/ping", menuitems.Retrieve)

		v1.GET("/secured/ping", middlewares.CheckJWT(), func(g *gin.Context) {
			g.JSON(200, gin.H{"text": "Hello from private"})
		})

		// handle oauth login
		fbAuth := facebook.Facebook{userStorage}
		v1.GET("/auth/facebook", fbAuth.BeginAuth)
		v1.GET("/auth/facebook/callback", fbAuth.Authenticate)
		gooAuth := google.Google{userStorage}
		v1.GET("/auth/google", gooAuth.BeginAuth)
		v1.GET("/auth/google/callback", gooAuth.Authenticate)

		// handle login
		account := controllers.AccountController{userStorage}
		v1.POST("/login", account.Authenticate)
		v1.POST("/signup", func(c *gin.Context) {
			account.Signup(c, mailSender)
		})
		v1.GET("/activate", account.Activate)

		// handle bookmarks of users
		bc := controllers.BookmarkController{bookmarkStorage, userStorage}
		v1.GET("/users/:userID/bookmarks", middlewares.ValidateUserID(), bc.ListBookmarkByUser)
		v1.POST("/users/:userID/bookmarks/", middlewares.ValidateUserID(), bc.CreateBookmarkByUser)
		v1.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.ValidateUserID(), bc.DeleteBookmarkByUser)
	}

	return router
}
