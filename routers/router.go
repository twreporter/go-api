package routers

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"

	"github.com/twreporter/go-api/controllers"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/middlewares"
)

const (
	maxAge = 3600
)

type wrappedFn func(c *gin.Context) (int, gin.H, error)

func ginResponseWrapper(fn wrappedFn) func(c *gin.Context) {
	return func(c *gin.Context) {
		statusCode, obj, err := fn(c)
		if err != nil {
			if globals.Conf.Environment == "development" {
				log.Errorf("%+v", err)
			} else {
				log.WithField("detail", err).Errorf("%s", f.FormatStack(err))
			}
		}
		c.JSON(statusCode, obj)
	}
}

// SetupRouter ...
func SetupRouter(cf *controllers.ControllerFactory) (engine *gin.Engine) {
	switch globals.Conf.Environment {
	case "production", "staging":
		// Disable default logger(stdout/stderr)
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(middlewares.Recovery())
		engine.Use(gin.LoggerWithFormatter(f.NewGinLogFormatter()))
	default:
		engine = gin.Default()
	}

	config := cors.DefaultConfig()

	var allowOrigins = globals.Conf.Cors.AllowOrigins
	if len(allowOrigins) > 0 {
		config.AllowOrigins = allowOrigins
	} else {
		switch globals.Conf.Environment {
		case globals.DevelopmentEnvironment:
			config.AllowAllOrigins = true
		case globals.StagingEnvironment:
			config.AllowOrigins = []string{globals.MainSiteStagingOrigin, globals.SupportSiteStagingOrigin, globals.AccountsSiteStagingOrigin}
		case globals.ProductionEnvironment:
			config.AllowOrigins = []string{globals.MainSiteOrigin, globals.SupportSiteOrigin, globals.AccountsSiteOrigin}
		default:
			// omit intentionally
		}
	}

	config.AddAllowHeaders("Authorization")
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
	// endpoints for bookmarks of users
	v1Group.GET("/users/:userID/bookmarks", middlewares.ValidateAuthorization(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	v1Group.GET("/users/:userID/bookmarks/:bookmarkSlug", middlewares.ValidateAuthorization(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	v1Group.POST("/users/:userID/bookmarks", middlewares.ValidateAuthorization(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.CreateABookmarkOfAUser))
	v1Group.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.ValidateAuthorization(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.DeleteABookmarkOfAUser))

	// endpoints for donation
	v1Group.POST("/periodic-donations", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.ValidateUserIDInReqBody(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.CreateAPeriodicDonationOfAUser))
	v1Group.PATCH("/periodic-donations/orders/:order", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.ValidateUserIDInReqBody(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.PatchADonationOfAUser(c, globals.PeriodicDonationType)
	}))
	v1Group.GET("/periodic-donations/orders/:order", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.GetADonationOfAUser(c, globals.PeriodicDonationType)
	}))
	v1Group.POST("/donations/prime", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.ValidateUserIDInReqBody(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.CreateADonationOfAUser))
	v1Group.PATCH("/donations/prime/orders/:order", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.ValidateUserIDInReqBody(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.PatchADonationOfAUser(c, globals.PrimeDonationType)
	}))
	// v1Group.GET("/users/:userID/donations", middlewares.ValidateAuthorization(), middlewares.ValidateUserID(), ginResponseWrapper(mc.GetDonationsOfAUser))
	// one-time donation including credit_card, line pay, apple pay, google pay and samsung pay
	v1Group.GET("/donations/prime/orders/:order", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.GetADonationOfAUser(c, globals.PrimeDonationType)
	}))
	v1Group.GET("/donations/prime/orders/:order/transaction_verification", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetVerificationInfoOfADonation))

	v1Group.POST("/donations/prime/line-notify", ginResponseWrapper(mc.PatchLinePayOfAUser))
	v1Group.POST("/tappay_query", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.QueryTappayServer))
	// TODO
	// donations derived from the periodic donation
	// v1Group.GET("/users/:userID/donations/token/:id", middlewares.ValidateAuthorization(), middlewares.ValidateUserID(), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
	//  return mc.GetADonationOfAUser(c, globals.TokenDonationType)
	//}))

	// other donations not included in the above endpoints
	v1Group.GET("/donations/others/orders/:order", middlewares.ValidateAuthentication(), middlewares.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.GetADonationOfAUser(c, globals.OthersDonationType)
	}))

	// endpoints for web push subscriptions
	v1Group.POST("/web-push/subscriptions" /*middlewares.ValidateAuthorization()*/, middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.SubscribeWebPush))
	v1Group.GET("/web-push/subscriptions", middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.IsWebPushSubscribed))

	// =============================
	// news service endpoints
	// =============================
	nc := cf.GetNewsController()
	// endpoints for authors
	v1Group.GET("/authors", middlewares.SetCacheControl("public,max-age=600"), ginResponseWrapper(nc.GetAuthors))
	// endpoints for posts
	v1Group.GET("/posts", middlewares.SetCacheControl("public,max-age=900"), ginResponseWrapper(nc.GetPosts))
	v1Group.GET("/posts/:slug", middlewares.SetCacheControl("public,max-age=900"), ginResponseWrapper(nc.GetAPost))
	// endpoints for topics
	v1Group.GET("/topics", middlewares.SetCacheControl("public,max-age=900"), ginResponseWrapper(nc.GetTopics))
	v1Group.GET("/topics/:slug", middlewares.SetCacheControl("public,max-age=900"), ginResponseWrapper(nc.GetATopic))
	v1Group.GET("/index_page", middlewares.SetCacheControl("public,max-age=1800"), nc.GetIndexPageContents)
	v1Group.GET("/index_page_categories", middlewares.SetCacheControl("public,max-age=1800"), nc.GetCategoriesPosts)
	// endpoints for search
	v1Group.GET("/search/authors", middlewares.SetCacheControl("public,max-age=3600"), nc.SearchAuthors)
	v1Group.GET("/search/posts", middlewares.SetCacheControl("public,max-age=3600"), nc.SearchPosts)

	// =============================
	// mail service endpoints
	// =============================

	mailContrl := cf.GetMailController()
	mailMiddleware := middlewares.GetMailServiceMiddleware()
	v1Group.POST(fmt.Sprintf("/%s", globals.SendActivationRoutePath), mailMiddleware.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mailContrl.SendActivation))
	v1Group.POST(fmt.Sprintf("/%s", globals.SendSuccessDonationRoutePath), mailMiddleware.ValidateAuthorization(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mailContrl.SendDonationSuccessMail))

	v2Group := engine.Group("/v2")
	ncV2 := cf.GetNewsV2Controller()
	v2Group.GET("/posts", middlewares.SetCacheControl("public,max-age=900"), ncV2.GetPosts)
	v2Group.GET("/posts/:slug", middlewares.SetCacheControl("public,max-age=900"), ncV2.GetAPost)
	// endpoints for topics
	v2Group.GET("/topics", middlewares.SetCacheControl("public,max-age=900"), ncV2.GetTopics)
	v2Group.GET("/topics/:slug", middlewares.SetCacheControl("public,max-age=900"), ncV2.GetATopic)
	v2Group.GET("/index_page", middlewares.SetCacheControl("public,max-age=1800"), ncV2.GetIndexPage)

	v2Group.GET("/authors", middlewares.SetCacheControl("public,max-age=600"), ncV2.GetAuthors)
	v2Group.GET("/authors/:author_id", middlewares.SetCacheControl("public,max-age=600"), ncV2.GetAuthorByID)
	// =============================
	// v2 oauth endpoints
	// =============================
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
	v2AuthGroup.POST("/token", middlewares.ValidateAuthentication(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.TokenDispatch))
	v2AuthGroup.GET("/logout", mc.TokenInvalidate)
	return
}
