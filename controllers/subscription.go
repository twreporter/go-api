package controllers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"

	//log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// IsWebPushSubscribed - which handles the HTTP Get request,
// and try to check if the web push subscription is existed or not
func (mc *MembershipController) IsWebPushSubscribed(c *gin.Context) (int, gin.H, error) {
	const errorWhere = "MembershipController.IsWebPushSubscribed"
	var endpoint = c.Query("endpoint")
	var err error
	var hashEndpoint string
	var wpSub models.WebPushSubscription

	if endpoint == "" {
		return http.StatusNotFound, gin.H{"status": "error", "message": "Fail to get a web push subscription since you do not provide endpoint in URL query param"}, nil
	}

	hashEndpoint = fmt.Sprintf("%x", md5.Sum([]byte(endpoint)))

	if wpSub, err = mc.Storage.GetAWebPushSubscriptionByHashEndpoint(hashEndpoint); err != nil {
		switch appErr := err.(type) {
		case models.AppError:
			return 0, gin.H{}, models.NewAppError(errorWhere, "Fail to get a web push subscription", appErr.Error(), appErr.StatusCode)
		default:
			return http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown Error Type. Fail to get a web push subscription"}, nil
		}
	}

	return http.StatusOK, gin.H{"status": "success", "data": wpSub}, nil
}

// SubscribeWebPush - which handles the HTTP POST request,
// and try to create a web push subscription record into the persistent database
func (mc *MembershipController) SubscribeWebPush(c *gin.Context) (int, gin.H, error) {
	// subscriptionBody is to store POST body
	type subscriptionBody struct {
		Endpoint       string `json:"endpoint" form:"endpoint" binding:"required"`
		Keys           string `json:"keys" form:"keys" binding:"required"`
		ExpirationTime string `json:"expirationTime" form:"expirationTime"`
		UserID         string `json:"user_id" form:"user_id"`
	}

	const errorWhere = "MembershipController.SubscribeWebPush"
	var err error
	var sBody subscriptionBody
	var expirationTime int64
	var userID uint64
	var wpSub models.WebPushSubscription

	if err = c.Bind(&sBody); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": subscriptionBody{
			Endpoint:       "endpoint is required, and need to be a string",
			Keys:           "keys is required, and need to be a string",
			ExpirationTime: "expirationTime is optional, if provide, need to be a string of timestamp",
			UserID:         "user_id is optional, if provide, need to be a string",
		}}, nil
	}

	// HashEndpoint is created by md5 hash.
	// It is a unique key in the persistent database
	// to avoid from creating the duplicate record
	wpSub = models.WebPushSubscription{
		Endpoint:     sBody.Endpoint,
		HashEndpoint: fmt.Sprintf("%x", md5.Sum([]byte(sBody.Endpoint))),
		Keys:         sBody.Keys,
	}

	if userID, err = strconv.ParseUint(sBody.UserID, 10, 0); err == nil {
		wpSub.SetUserID(uint(userID))
	}

	if expirationTime, err = strconv.ParseInt(sBody.ExpirationTime, 10, 64); err == nil {
		wpSub.SetExpirationTime(expirationTime)
	}

	if err = mc.Storage.CreateAWebPushSubscription(wpSub); err != nil {
		switch appErr := err.(type) {
		case models.AppError:
			return 0, gin.H{}, models.NewAppError(errorWhere, "Fails to create a web push subscription", appErr.Error(), appErr.StatusCode)
		default:
			return http.StatusInternalServerError, gin.H{"status": "error", "message": "Unknown Error Type. Fails to create a web push subscription"}, nil
		}
	}

	return http.StatusCreated, gin.H{"status": "success", "data": sBody}, nil
}
