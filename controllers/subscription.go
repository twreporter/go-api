package controllers

import (
	"hash/crc32"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/models"
)

// IsWebPushSubscribed - which handles the HTTP Get request,
// and try to check if the web push subscription is existed or not
func (mc *MembershipController) IsWebPushSubscribed(c *gin.Context) (int, gin.H, error) {
	endpoint := c.Query("endpoint")
	if endpoint == "" {
		return http.StatusNotFound, gin.H{"status": "error", "message": "Fail to get a web push subscription since you do not provide endpoint in URL query param"}, nil
	}

	crc32Endpoint := crc32.Checksum([]byte(endpoint), crc32.IEEETable)

	wpSub, err := mc.Storage.GetAWebPushSubscription(crc32Endpoint, endpoint)
	if err != nil {
		return toResponse(err)
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

	var endpoint string
	var err error
	var expirationTime int64
	var crc32Endpoint uint32
	var sBody subscriptionBody
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

	endpoint = sBody.Endpoint
	crc32Endpoint = crc32.Checksum([]byte(endpoint), crc32.IEEETable)

	wpSub = models.WebPushSubscription{
		Endpoint:      endpoint,
		Crc32Endpoint: crc32Endpoint,
		Keys:          sBody.Keys,
	}

	if userID, err = strconv.ParseUint(sBody.UserID, 10, 0); err == nil {
		wpSub.SetUserID(uint(userID))
	}

	if expirationTime, err = strconv.ParseInt(sBody.ExpirationTime, 10, 64); err == nil {
		wpSub.SetExpirationTime(expirationTime)
	}

	if err = mc.Storage.CreateAWebPushSubscription(wpSub); err != nil {
		return toResponse(err)
	}

	return http.StatusCreated, gin.H{"status": "success", "data": sBody}, nil
}
