package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"
)

type DefaultVariables struct {
	Account string
	Service string
	Token   string

	// objectID
	ImgID1  bson.ObjectId
	ImgID2  bson.ObjectId
	VideoID bson.ObjectId
	PostID1 bson.ObjectId
	PostID2 bson.ObjectId
	TopicID bson.ObjectId
	TagID   bson.ObjectId
	CatID   bson.ObjectId
	ThemeID bson.ObjectId

	// collection
	ImgCol1  models.MongoImage
	ImgCol2  models.MongoImage
	VideoCol models.MongoVideo
	PostCol1 models.Post
	PostCol2 models.Post
	TagCol   models.Tag
	CatCol   models.Category
	TopicCol models.Topic
	ThemeCol models.Theme

	MockPostSlug1 string
	MockTopicSlug string

	// web push
	WebPushEndpoint string
}

type GlobalVariables struct {
	Defaults  DefaultVariables
	GinEngine *gin.Engine
	GormDB    *gorm.DB
	MgoDB     *mgo.Session
}

type WebPushSubscriptionPostBody struct {
	Endpoint       string `json:"endpoint"`
	Keys           string `json:"keys"`
	ExpirationTime string `json:"expiration_time"`
	UserID         string `json:"user_id"`
}
