package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/twreporter/go-api/configs"
	"github.com/twreporter/go-api/controllers"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/routers"
	"github.com/twreporter/go-api/storage"
	"github.com/twreporter/go-api/utils"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Globs globalVariables

var testMongoClient *mongodriver.Client

func init() {
	var defaults = defaultVariables{
		Account:          "developer@twreporter.org",
		Service:          "default_service",
		Token:            "default_token",
		ImgID1:           bson.NewObjectId(),
		ImgID2:           bson.NewObjectId(),
		VideoID:          bson.NewObjectId(),
		PostID1:          bson.NewObjectId(),
		PostID2:          bson.NewObjectId(),
		TopicID:          bson.NewObjectId(),
		TagID:            bson.NewObjectId(),
		CatReviewID:      bson.ObjectIdHex(configs.ReviewListID),
		CatPhotographyID: bson.ObjectIdHex(configs.PhotographyListID),
		ThemeID:          bson.NewObjectId(),

		MockPostSlug1: "mock-post-slug-1",
		MockTopicSlug: "mock-topic-slug",

		ErrorEmailAddress: "error@twreporter.org",
	}

	img1 := models.MongoImage{
		ID:          defaults.ImgID1,
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
	defaults.ImgCol1 = img1
	img2 := img1
	img2.ID = defaults.ImgID2
	defaults.ImgCol2 = img2

	video := models.MongoVideo{
		ID:    defaults.VideoID,
		Title: "mock video title",
		Video: models.MongoVideoAsset{
			Filetype: "video/mp4",
			Size:     1000,
			URL:      "https://www.twreporter.org/videos/mock-video.mp4",
		},
	}
	defaults.VideoCol = video

	tag := models.Tag{
		ID:   defaults.TagID,
		Name: "mock tag",
	}
	defaults.TagCol = tag

	defaults.CatReviewCol = models.Category{
		ID:   defaults.CatReviewID,
		Name: "評論",
	}

	defaults.CatPhotographyCol = models.Category{
		ID:   defaults.CatPhotographyID,
		Name: "攝影",
	}

	theme := models.Theme{
		ID:            defaults.ThemeID,
		Name:          "photograph",
		TitlePosition: "title-above",
	}
	defaults.ThemeCol = theme

	post1 := models.Post{
		ID:               defaults.PostID1,
		Slug:             defaults.MockPostSlug1,
		Name:             "mock post slug 1",
		Style:            "article:v2:default",
		State:            "published",
		ThemeOrigin:      defaults.ThemeID,
		PublishedDate:    time.Now(),
		HeroImageOrigin:  defaults.ImgID1,
		CategoriesOrigin: []bson.ObjectId{defaults.CatPhotographyID},
		OgImageOrigin:    defaults.ImgID1,
		IsFeatured:       true,
		TopicOrigin:      defaults.TopicID,
		RelatedsOrigin:   []bson.ObjectId{defaults.PostID2},
	}
	defaults.PostCol1 = post1

	post2 := models.Post{
		ID:                         defaults.PostID2,
		Slug:                       "mock-post-slug-2",
		Name:                       "mock post slug 2",
		Style:                      "article:v2:default",
		State:                      "published",
		ThemeOrigin:                defaults.ThemeID,
		PublishedDate:              post1.PublishedDate.Add(time.Duration(1) * time.Minute),
		HeroImageOrigin:            defaults.ImgID2,
		CategoriesOrigin:           []bson.ObjectId{defaults.CatReviewID},
		OgImageOrigin:              defaults.ImgID2,
		IsFeatured:                 false,
		LeadingImagePortraitOrigin: defaults.ImgID1,
		TopicOrigin:                defaults.TopicID,
		RelatedsOrigin:             []bson.ObjectId{defaults.PostID2},
		TagsOrigin:                 []bson.ObjectId{defaults.TagID},
	}
	defaults.PostCol2 = post2

	topic := models.Topic{
		ID:                 defaults.TopicID,
		Slug:               defaults.MockTopicSlug,
		TopicName:          "mock topic slug",
		Title:              "mock title",
		State:              "published",
		RelatedsOrigin:     []bson.ObjectId{defaults.PostID1, defaults.PostID2},
		LeadingImageOrigin: defaults.ImgID1,
		LeadingVideoOrigin: defaults.VideoID,
		OgImageOrigin:      defaults.ImgID1,
	}
	defaults.TopicCol = topic

	Globs = globalVariables{
		Defaults: defaults,
	}
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

func generateIDToken(user models.User) (jwt string) {
	jwt, _ = utils.RetrieveV2IDToken(user.ID, user.Email.ValueOrZero(), user.FirstName.ValueOrZero(), user.LastName.ValueOrZero(), 3600)
	return
}

func getReporterAccount(email string) (ra models.ReporterAccount) {
	as := storage.NewGormStorage(Globs.GormDB)
	ra, _ = as.GetReporterAccountData(email)
	return ra
}

func createUser(email string) models.User {
	as := storage.NewGormStorage(Globs.GormDB)

	ra := models.ReporterAccount{
		Email:         email,
		ActivateToken: Globs.Defaults.Token,
		ActExpTime:    time.Now().Add(time.Duration(15) * time.Minute),
	}

	user, _ := as.InsertUserByReporterAccount(ra)

	return user
}

func deleteUser(user models.User) {
	db := Globs.GormDB

	// Remove corresponding reporter account
	db.Unscoped().Delete(user.ReporterAccount)
	db.Unscoped().Delete(user)
}

func getUser(email string) (user models.User) {
	as := storage.NewGormStorage(Globs.GormDB)
	user, _ = as.GetUserByEmail(email)
	return
}

func serveHTTP(method, path, body, contentType, authorization string) (resp *httptest.ResponseRecorder) {
	var req *http.Request

	req = requestWithBody(method, path, body)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	resp = httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)

	return
}

func serveHTTPWithCookies(method, path, body, contentType, authorization string, cookies ...http.Cookie) (resp *httptest.ResponseRecorder) {
	var req *http.Request

	req = requestWithBody(method, path, body)

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	if authorization != "" {
		req.Header.Add("Authorization", authorization)
	}

	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	resp = httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)

	return
}

type mockMailStrategy struct{}

func (s mockMailStrategy) Send(to, subject, body string) error {
	if to == Globs.Defaults.ErrorEmailAddress {
		return errors.New("mail service works abnormally")
	}
	return nil
}

func setupGinServer(gormDB *gorm.DB, mgoDB *mgo.Session, client *mongodriver.Client) *gin.Engine {
	mailSvc := mockMailStrategy{}
	cf := controllers.NewControllerFactory(gormDB, mgoDB, mailSvc, client)
	engine := routers.SetupRouter(cf)
	return engine
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	resp := httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestMain(m *testing.M) {
	var err error
	var l net.Listener

	fmt.Println("load default config")
	if globals.Conf, err = configs.LoadDefaultConf(); err != nil {
		panic(fmt.Sprintf("Can not load default config, but got err=%+v", err))
	}

	// set up DB environment
	gormDB, mgoDB, client := setUpDBEnvironment()

	Globs.GormDB = gormDB
	Globs.MgoDB = mgoDB
	testMongoClient = client

	// set up gin server
	engine := setupGinServer(gormDB, mgoDB, client)

	Globs.GinEngine = engine

	defer Globs.GormDB.Close()
	defer Globs.MgoDB.Close()
	defer func() { testMongoClient.Disconnect(context.Background()) }()

	// start server for testing
	// the reason why we start the server
	// is because we send HTTP request internally between controllers
	ts := httptest.NewUnstartedServer(engine)
	if l, err = net.Listen("tcp", "127.0.0.1:8080"); err != nil {
		panic(err)
	}
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	retCode := m.Run()
	os.Exit(retCode)
}
