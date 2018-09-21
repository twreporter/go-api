package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/routers"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

var (
	DefaultAccount = "developer@twreporter.org"
	DefaultService = "default_service"
	DefaultToken   = "default_token"
	Engine         *gin.Engine
	DB             *gorm.DB
	MgoDB          *mgo.Session
	MgoDBName      = "gorm"

	// collections name
	MgoPostCol       = "posts"
	MgoTopicCol      = "topics"
	MgoImgCol        = "images"
	MgoVideoCol      = "videos"
	MgoTagCol        = "tags"
	MgoCategoriesCol = "postcategories"
	MgoThemeCol      = "themes"

	// objectID
	ImgID1  = bson.NewObjectId()
	ImgID2  = bson.NewObjectId()
	VideoID = bson.NewObjectId()
	PostID1 = bson.NewObjectId()
	PostID2 = bson.NewObjectId()
	TopicID = bson.NewObjectId()
	TagID   = bson.NewObjectId()
	CatID   = bson.NewObjectId()
	ThemeID = bson.NewObjectId()

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

	MockPostSlug1 = "mock-post-slug-1"
	MockTopicSlug = "mock-topic-slug"
)

func OpenGormConnection() (db *gorm.DB, err error) {
	db, _ = utils.InitDB(10, 5)

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}

	return
}

func OpenMgoConnection() (session *mgo.Session, err error) {
	session, err = utils.InitMongoDB()

	return
}

func RunMigration() {
	RunGormMigration()
	RunMgoMigration()
}

func RunGormMigration() {
	values := []interface{}{&models.User{}, &models.OAuthAccount{}, &models.ReporterAccount{}, &models.Bookmark{}, &models.Registration{}, &models.Service{}, &models.UsersBookmarks{}, &models.WebPushSubscription{}}
	for _, value := range values {
		DB.DropTable(value)
	}
	if err := DB.AutoMigrate(values...).Error; err != nil {
		panic(fmt.Sprintf("No error should happen when create table, but got %+v", err))
	}
}

func RunMgoMigration() {
	err := MgoDB.DB(MgoDBName).DropDatabase()
	if err != nil {
		panic(fmt.Sprint("Can not drop mongo gorm database"))
	}

	MgoDB.DB(MgoDBName).Run(bson.D{{"create", MgoPostCol}}, nil)
	MgoDB.DB(MgoDBName).Run(bson.D{{"create", MgoTopicCol}}, nil)
	MgoDB.DB(MgoDBName).Run(bson.D{{"create", MgoImgCol}}, nil)
	MgoDB.DB(MgoDBName).Run(bson.D{{"create", MgoVideoCol}}, nil)
	MgoDB.DB(MgoDBName).Run(bson.D{{"create", MgoTagCol}}, nil)
	MgoDB.DB(MgoDBName).Run(bson.D{{"create", MgoCategoriesCol}}, nil)
}

func SetDefaultRecords() {
	SetGormDefaultRecords()
	SetMgoDefaultRecords()
}

func SetMgoDefaultRecords() {
	ImgCol1 = models.MongoImage{
		ID:          ImgID1,
		Description: "mock image desc",
		Copyright:   "",
		Image: models.MongoImageAsset{
			Height:   1200,
			Filetype: "image/jpg",
			Width:    2000,
			URL:      "https://www.twreporter.org/images/mock-image-1.jpg",
			ResizedTargets: models.ResizedTargets{
				Mobile: models.ImageAsset{
					Height: 600,
					Width:  800,
					URL:    "https://www.twreporter.org/images/mock-image-1-mobile.jpg",
				},
				Tablet: models.ImageAsset{
					Height: 1000,
					Width:  1400,
					URL:    "https://www.twreporter.org/images/mock-image-1-tablet.jpg",
				},
				Desktop: models.ImageAsset{
					Height: 1200,
					Width:  2000,
					URL:    "https://www.twreporter.org/images/mock-image-1-desktop.jpg",
				},
				Tiny: models.ImageAsset{
					Height: 60,
					Width:  80,
					URL:    "https://www.twreporter.org/images/mock-image-1-tiny.jpg",
				},
				W400: models.ImageAsset{
					Height: 300,
					Width:  400,
					URL:    "https://www.twreporter.org/images/mock-image-1-w400.jpg",
				},
			},
		},
	}

	VideoCol = models.MongoVideo{
		ID:    VideoID,
		Title: "mock video title",
		Video: models.MongoVideoAsset{
			Filetype: "video/mp4",
			Size:     1000,
			URL:      "https://www.twreporter.org/videos/mock-video.mp4",
		},
	}

	TagCol = models.Tag{
		ID:   TagID,
		Name: "mock tag",
	}

	CatCol = models.Category{
		ID:   CatID,
		Name: "mock postcategory",
	}

	ThemeCol = models.Theme{
		ID:            ThemeID,
		Name:          "photograph",
		TitlePosition: "title-above",
	}

	PostCol1 = models.Post{
		ID:               PostID1,
		Slug:             MockPostSlug1,
		Name:             "mock post slug 1",
		Style:            "article",
		State:            "published",
		ThemeOrigin:      ThemeID,
		PublishedDate:    time.Now(),
		HeroImageOrigin:  ImgID1,
		CategoriesOrigin: []bson.ObjectId{CatID},
		OgImageOrigin:    ImgID1,
		IsFeatured:       true,
		TopicOrigin:      TopicID,
		RelatedsOrigin:   []bson.ObjectId{PostID2},
	}

	TopicCol = models.Topic{
		ID:                 TopicID,
		Slug:               MockTopicSlug,
		TopicName:          "mock topic slug",
		Title:              "mock title",
		State:              "published",
		RelatedsOrigin:     []bson.ObjectId{PostID1, PostID2},
		LeadingImageOrigin: ImgID1,
		LeadingVideoOrigin: VideoID,
		OgImageOrigin:      ImgID1,
	}
	// insert img1 and  img2
	MgoDB.DB(MgoDBName).C(MgoImgCol).Insert(ImgCol1)

	ImgCol2 = ImgCol1
	ImgCol2.ID = ImgID2
	MgoDB.DB(MgoDBName).C(MgoImgCol).Insert(ImgCol2)

	// insert video
	MgoDB.DB(MgoDBName).C(MgoVideoCol).Insert(VideoCol)

	// insert tag and postcategory
	MgoDB.DB(MgoDBName).C(MgoTagCol).Insert(TagCol)
	MgoDB.DB(MgoDBName).C(MgoCategoriesCol).Insert(CatCol)

	// insert post1 and post2
	MgoDB.DB(MgoDBName).C(MgoPostCol).Insert(PostCol1)

	// insert theme
	MgoDB.DB(MgoDBName).C(MgoThemeCol).Insert(ThemeCol)

	PostCol2 = PostCol1
	PostCol2.ID = PostID2
	PostCol2.Slug = "mock-post-slug-2"
	PostCol2.Name = "mock post slug 2"
	PostCol2.Style = "review"
	PostCol2.PublishedDate = time.Now()
	PostCol2.HeroImageOrigin = ImgID2
	PostCol2.LeadingImagePortraitOrigin = ImgID1
	PostCol2.OgImageOrigin = ImgID2
	PostCol2.IsFeatured = false
	PostCol2.TagsOrigin = []bson.ObjectId{TagID}
	MgoDB.DB(MgoDBName).C(MgoPostCol).Insert(PostCol2)

	// insert topic
	MgoDB.DB(MgoDBName).C(MgoTopicCol).Insert(TopicCol)
}

func SetGormDefaultRecords() {
	// Set an active reporter account
	ms := storage.NewGormStorage(DB)

	ra := models.ReporterAccount{
		Email:         DefaultAccount,
		ActivateToken: DefaultToken,
		ActExpTime:    time.Now().Add(time.Duration(15) * time.Minute),
	}
	_, _ = ms.InsertUserByReporterAccount(ra)

	ms.CreateService(models.ServiceJSON{Name: DefaultService})

	ms.CreateRegistration(DefaultService, models.RegistrationJSON{Email: DefaultAccount, ActivateToken: DefaultToken})
}

type mockMailStrategy struct{}

func (s mockMailStrategy) Send(to, subject, body string) error {
	return nil
}

func SetupGinServer() {
	strategy := mockMailStrategy{}
	cf := controllers.NewControllerFactory(DB, MgoDB, utils.NewEmailSender(strategy))
	Engine = routers.SetupRouter(cf)
}

func RequestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

func GenerateJWT(user models.User) (jwt string) {
	jwt, _ = utils.RetrieveToken(user.ID, user.Email.String)
	return
}

func GetUser(email string) (user models.User) {
	as := storage.NewGormStorage(DB)
	user, _ = as.GetUserByEmail(email)
	return
}

func ServeHTTP(method, path, body, contentType, authorization string) (resp *httptest.ResponseRecorder) {
	var req *http.Request

	req = RequestWithBody(method, path, body)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)

	return
}
