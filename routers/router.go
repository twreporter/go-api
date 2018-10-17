package routers

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo"
	"github.com/gin-gonic/gin"

	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
)

const (
	maxAge = 3600
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

	if globals.Conf.Environment != "development" {
		var allowOrigins = globals.Conf.Cors.AllowOrigins
		if len(allowOrigins) > 0 {
			config.AllowOrigins = allowOrigins
		} else {
			config.AllowOrigins = []string{"https://www.twreporter.org"}
		}
	} else {
		config.AllowAllOrigins = true
	}

	config.AddAllowHeaders("Authorization")
	config.AddAllowHeaders("ContentType")
	config.AddAllowMethods("DELETE")
	config.AddAllowMethods("PATCH")

	// Enable Access-Control-Allow-Credentials header for axios pre-flight(OPTION) request
	// so that the subsequent request could carry cookie
	config.AllowCredentials = true

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
	v1Group.POST("/signin", middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.SignIn))
	v1Group.GET("/activate", middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.Activate))
	v1Group.GET("/token/:userID", middlewares.CheckJWT(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.RenewJWT))
	// endpoints for bookmarks of users
	v1Group.GET("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	v1Group.GET("/users/:userID/bookmarks/:bookmarkSlug", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	v1Group.POST("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.CreateABookmarkOfAUser))
	v1Group.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.DeleteABookmarkOfAUser))

	// endpoints for donation
	v1Group.POST("/periodic-donations", middlewares.CheckJWT(), middlewares.ValidateUserIDInReqBody(), ginResponseWrapper(mc.CreateAPeriodicDonationOfAUser))
	v1Group.PATCH("/periodic-donations/:id", middlewares.CheckJWT(), middlewares.ValidateUserIDInReqBody(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.PatchADonationOfAUser(c, globals.PeriodicDonationType)
	}))
	v1Group.GET("/periodic-donations/:id", middlewares.ValidateAuthentication(), middlewares.CheckJWT(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.GetADonationOfAUser(c, globals.PeriodicDonationType)
	}))
	v1Group.POST("/donations/:pay_method", middlewares.CheckJWT(), middlewares.ValidateUserIDInReqBody(), ginResponseWrapper(mc.CreateADonationOfAUser))
	v1Group.PATCH("/donations/prime/:id", middlewares.CheckJWT(), middlewares.ValidateUserIDInReqBody(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.PatchADonationOfAUser(c, globals.PrimeDonaitionType)
	}))
	// v1Group.GET("/users/:userID/donations", middlewares.CheckJWT(), middlewares.ValidateUserID(), ginResponseWrapper(mc.GetDonationsOfAUser))
	// one-time donation including credit_card, line pay, apple pay, google pay and samsung pay
	v1Group.GET("/donations/prime/:id", middlewares.ValidateAuthentication(), middlewares.CheckJWT(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.GetADonationOfAUser(c, globals.PrimeDonaitionType)
	}))

	// TODO
	// donations derived from the periodic donation
	// v1Group.GET("/users/:userID/donations/token/:id", middlewares.CheckJWT(), middlewares.ValidateUserID(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
	//  return mc.GetADonationOfAUser(c, globals.TokenDonationType)
	//}))

	// other donations not included in the above endpoints
	v1Group.GET("/donations/others/:id", middlewares.ValidateAuthentication(), middlewares.CheckJWT(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.GetADonationOfAUser(c, globals.OthersDonationType)
	}))

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
	// v1 oauth endpoints
	// =============================
	authGroup := v1Group.Group("/auth")

	gc := cf.GetGoogleController()
	authGroup.GET("/google", middlewares.SetCacheControl("no-store"), gc.BeginAuth)
	authGroup.GET("/google/callback", middlewares.SetCacheControl("no-store"), gc.Authenticate)
	fc := cf.GetFacebookController()
	authGroup.GET("/facebook", middlewares.SetCacheControl("no-store"), fc.BeginAuth)
	authGroup.GET("/facebook/callback", middlewares.SetCacheControl("no-store"), fc.Authenticate)

	// =============================
	// mail service endpoints
	// =============================

	mailContrl := cf.GetMailController()
	v1Group.POST(fmt.Sprintf("/%s", globals.SendActivationRoutePath), ginResponseWrapper(mailContrl.SendActivation))
	v1Group.POST(fmt.Sprintf("/%s", globals.SendSuccessDonationRoutePath), ginResponseWrapper(mailContrl.SendDonationSuccessMail))

	// =============================
	// v2 oauth endpoints
	// =============================
	v2Group := engine.Group("/v2")
	v2AuthGroup := v2Group.Group("/auth")

	session := cf.GetMgoSession()
	c := session.DB("go-api").C("sessions")
	store := mongo.NewStore(c, maxAge, true, []byte("secret"))
	v2AuthGroup.Use(sessions.Sessions("go-api-session", store))
	store.Options(sessions.Options{
		Domain:   globals.Conf.App.Domain,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   globals.Conf.Environment != "development",
	})

	ogc := cf.GetOAuthController(globals.GoogleOAuth)
	v2AuthGroup.GET("/google", middlewares.SetCacheControl("no-store"), ogc.BeginOAuth)
	v2AuthGroup.GET("/google/callback", middlewares.SetCacheControl("no-store"), ogc.Authenticate)
	ofc := cf.GetOAuthController(globals.FacebookOAuth)
	v2AuthGroup.GET("/facebook", middlewares.SetCacheControl("no-store"), ofc.BeginOAuth)
	v2AuthGroup.GET("/facebook/callback", middlewares.SetCacheControl("no-store"), ofc.Authenticate)

	// =============================
	// v2 membership service endpoints
	// =============================
	v2AuthGroup.POST("/signin", middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.SignInV2))
	v2AuthGroup.GET("/activate", middlewares.SetCacheControl("no-store"), mc.ActivateV2)
	v2AuthGroup.POST("/token", middlewares.CheckJWT(), middlewares.SetCacheControl("no-store"), mc.TokenDispatch)
	v2AuthGroup.GET("/logout", mc.TokenInvalidate)
	return engine
}
