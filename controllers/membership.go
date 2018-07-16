package controllers

import (
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/middlewares"
	// "twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
	//log "github.com/Sirupsen/logrus"
)

// NewMembershipController ...
func NewMembershipController(s storage.MembershipStorage) Controller {
	return &MembershipController{s}
}

// MembershipController ...
type MembershipController struct {
	Storage storage.MembershipStorage
}

// Close is the method of Controller interface
func (mc *MembershipController) Close() error {
	err := mc.Storage.Close()
	if err != nil {
		return err
	}
	return nil
}

// SetRoute is the method of Controller interface
func (mc *MembershipController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	// mailSender := utils.NewSMTPEmailSender()                          // use office365 to send mails
	mailSender := utils.NewAmazonEmailSender() // use Amazon SES to send mails

	// endpoints for account
	group.POST("/signin", middlewares.SetCacheControl("no-store"), ginResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.SignIn(c, mailSender)
	}))
	group.GET("/activate", middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.Activate))
	group.GET("/token/:userID", middlewares.CheckJWT(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.RenewJWT))

	// endpoints for bookmarks of users
	group.GET("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	group.GET("/users/:userID/bookmarks/:bookmarkSlug", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.GetBookmarksOfAUser))
	group.POST("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.CreateABookmarkOfAUser))
	group.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), ginResponseWrapper(mc.DeleteABookmarkOfAUser))

	// endpoints for donation
	group.POST("/users/:userID/periodic_donations" /*middlewares.CheckJWT(),*/, ginResponseWrapper(mc.CreateAPeriodicDonationOfAUser))
	group.POST("/users/:userID/donations/:pay_method" /*middlewares.CheckJWT(),*/, ginResponseWrapper(mc.CreateADonationOfAUser))
	group.GET("/users/:userID/donations" /*middlewares.CheckJWT(),*/, ginResponseWrapper(mc.GetDonationsOfAUser))

	// endpoints for web push subscriptions
	group.POST("/web-push/subscriptions" /*middlewares.CheckJWT()*/, ginResponseWrapper(mc.SubscribeWebPush))
	group.GET("/web-push/subscriptions", ginResponseWrapper(mc.IsWebPushSubscribed))

	/* Comment out these unused endpoints
	// endpoint for registration
	// TODO add middleware to check the request from twreporter.org domain
	group.POST("/registrations/:service", mc.Register)
	// TODO add middleware to check the email to delete is the email of the user sending the request
	group.DELETE("/registrations/:service/:userEmail", mc.Deregister)
	// TODO add middleware to check the request from twreporter.org domain
	group.GET("/registrations/:service/:userEmail", mc.GetRegisteredUser)
	group.GET("/registrations/:service", mc.GetRegisteredUsers)
	group.GET("/activation/:service/:userEmail", mc.ActivateRegistration)

	// endpoints for service
	group.POST("/services", middlewares.CheckJWT(), middlewares.ValidateAdminUsers(), mc.Create)
	group.DELETE("/services/:name", middlewares.CheckJWT(), middlewares.ValidateAdminUsers(), mc.Delete)
	group.PUT("/services/:name", middlewares.CheckJWT(), middlewares.ValidateAdminUsers(), mc.Update)
	group.GET("/services/:name", middlewares.CheckJWT(), mc.Read)
	*/

	return group
}
