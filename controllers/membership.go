package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/middlewares"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
	//	log "github.com/Sirupsen/logrus"
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

type wrappedFn func(c *gin.Context) (int, gin.H, error)

// GinResponseWrapper ...
func GinResponseWrapper(fn wrappedFn) func(c *gin.Context) {
	return func(c *gin.Context) {
		statusCode, obj, err := fn(c)
		if err != nil {
			switch appErr := err.(type) {
			case *models.AppError:
				c.JSON(appErr.StatusCode, gin.H{"status": "error", "message": appErr.Message})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Unknown Error. %s", err.Error())})
				return
			}
		}
		c.JSON(statusCode, obj)
	}
}

// SetRoute is the method of Controller interface
func (mc *MembershipController) SetRoute(group *gin.RouterGroup) *gin.RouterGroup {
	// mailSender := utils.NewSMTPEmailSender()                          // use office365 to send mails
	mailSender := utils.NewAmazonEmailSender() // use Amazon SES to send mails

	// endpoints for account
	group.POST("/signin", middlewares.SetCacheControl("no-store"), GinResponseWrapper(func(c *gin.Context) (int, gin.H, error) {
		return mc.SignIn(c, mailSender)
	}))
	group.GET("/activate", middlewares.SetCacheControl("no-store"), GinResponseWrapper(mc.Activate))
	group.GET("/token/:userID", middlewares.CheckJWT(), middlewares.SetCacheControl("no-store"), GinResponseWrapper(mc.RenewJWT))

	// endpoints for bookmarks of users
	group.GET("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), GinResponseWrapper(mc.GetBookmarksOfAUser))
	group.GET("/users/:userID/bookmarks/:bookmarkSlug", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), GinResponseWrapper(mc.GetBookmarksOfAUser))
	group.POST("/users/:userID/bookmarks", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), GinResponseWrapper(mc.CreateABookmarkOfAUser))
	group.DELETE("/users/:userID/bookmarks/:bookmarkID", middlewares.CheckJWT(), middlewares.ValidateUserID(), middlewares.SetCacheControl("no-store"), GinResponseWrapper(mc.DeleteABookmarkOfAUser))

	// endpoints for web push subscriptions
	group.POST("/web-push/subscriptions" /*middlewares.CheckJWT()*/, GinResponseWrapper(mc.SubscribeWebPush))
	group.GET("/web-push/subscriptions", GinResponseWrapper(mc.IsWebPushSubscribed))

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
