package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
)

type defaultVariables struct {
	Account string
	Service string
	Token   string

	ErrorEmailAddress string
}

type globalVariables struct {
	Defaults  defaultVariables
	GinEngine *gin.Engine
	GormDB    *gorm.DB
	MgoDB     *mgo.Session
}

type webPushSubscriptionPostBody struct {
	Endpoint       string `json:"endpoint"`
	Keys           string `json:"keys"`
	ExpirationTime string `json:"expiration_time"`
	UserID         string `json:"user_id"`
}
