package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"twreporter.org/go-api/models"
	//log "github.com/Sirupsen/logrus"
)

// SubscribeWebpush - which handles the HTTP POST request,
// and try to create a webpush subscription record into the persistent database
func (mc *MembershipController) SubscribeWebpush(c *gin.Context) (int, gin.H, error) {
	// SubscriptionBody is to store POST body
	type SubscriptionBody struct {
		Endpoint       string `json:"endpoint" form:"endpoint" binding:"required"`
		Keys           string `json:"keys" form:"keys" binding:"required"`
		ExpirationTime string `json:"expirationTime" form:"expirationTime"`
		UserID         string `json:"user_id" form:"user_id"`
	}

	const errorWhere = "MembershipController.SubscribeWebpush"
	var err error
	var sBody SubscriptionBody
	var expirationTime int64
	var userID uint64
	var wpSub models.WebpushSubscription

	if err = c.Bind(&sBody); err != nil {
		return http.StatusBadRequest, gin.H{"status": "fail", "data": SubscriptionBody{
			Endpoint:       "endpoint is required",
			Keys:           "keys is required",
			ExpirationTime: "expirationTime is optional",
			UserID:         "user_id is optional",
		}}, nil
	}

	// HashEndpoint is created by md5 hash.
	// It is a unique key in the persistent database
	// to avoid from creating the duplicate record
	wpSub = models.WebpushSubscription{
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

	if err = mc.Storage.CreateAWebpushSubscription(wpSub); err != nil {
		appErr, _ := err.(models.AppError)
		return 0, gin.H{}, models.NewAppError(errorWhere, "Creating a webpush subscription fails", err.Error(), appErr.StatusCode)
	}

	return http.StatusCreated, gin.H{"status": "success", "data": sBody}, nil
}
