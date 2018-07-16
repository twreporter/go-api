package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/controllers"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
	"twreporter.org/go-api/utils"
)

var Globs globalVariables

func init() {
	var defaults = defaultVariables{
		Account: "developer@twreporter.org",
		Service: "default_service",
		Token:   "default_token",
		ImgID1:  bson.NewObjectId(),
		ImgID2:  bson.NewObjectId(),
		VideoID: bson.NewObjectId(),
		PostID1: bson.NewObjectId(),
		PostID2: bson.NewObjectId(),
		TopicID: bson.NewObjectId(),
		TagID:   bson.NewObjectId(),
		CatID:   bson.NewObjectId(),
		ThemeID: bson.NewObjectId(),

		MockPostSlug1: "mock-post-slug-1",
		MockTopicSlug: "mock-topic-slug",
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

	cat := models.Category{
		ID:   defaults.CatID,
		Name: "mock postcategory",
	}
	defaults.CatCol = cat

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
		Style:            "article",
		State:            "published",
		ThemeOrigin:      defaults.ThemeID,
		PublishedDate:    time.Now(),
		HeroImageOrigin:  defaults.ImgID1,
		CategoriesOrigin: []bson.ObjectId{defaults.CatID},
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
		Style:                      "review",
		State:                      "published",
		ThemeOrigin:                defaults.ThemeID,
		PublishedDate:              post1.PublishedDate.Add(time.Duration(1) * time.Minute),
		HeroImageOrigin:            defaults.ImgID2,
		CategoriesOrigin:           []bson.ObjectId{defaults.CatID},
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

func generateJWT(user models.User) (jwt string) {
	jwt, _ = utils.RetrieveToken(user.ID, user.Email.String)
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

func setupGinServer(gormDB *gorm.DB, mgoDB *mgo.Session) *gin.Engine {
	// set up data storage
	gs := storage.NewGormStorage(gormDB)

	// init controllers
	mc := controllers.NewMembershipController(gs)
	fc := controllers.Facebook{Storage: gs}
	gc := controllers.Google{Storage: gs}

	ms := storage.NewMongoStorage(mgoDB)
	nc := controllers.NewNewsController(ms)

	cf := &controllers.ControllerFactory{
		Controllers: make(map[string]controllers.Controller),
	}
	cf.SetController(constants.MembershipController, mc)
	cf.SetController(constants.FacebookController, fc)
	cf.SetController(constants.GoogleController, gc)
	cf.SetController(constants.NewsController, nc)

	engine := gin.Default()
	routerGroup := engine.Group("/v1")
	{
		menuitems := new(controllers.MenuItemsController)
		routerGroup.GET("/ping", menuitems.Retrieve)
	}

	routerGroup = cf.SetRoute(routerGroup)

	return engine
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/ping", nil)
	resp := httptest.NewRecorder()
	Globs.GinEngine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestMain(m *testing.M) {
	// set default config
	utils.Cfg.SetDefaults()
	// set default mongo database to 'mgo'
	utils.Cfg.MongoDBSettings.DBName = mgoDBName

	// set up DB environment
	gormDB, mgoDB := setUpDBEnvironment()

	Globs.GormDB = gormDB
	Globs.MgoDB = mgoDB

	// set up gin server
	engine := setupGinServer(gormDB, mgoDB)
	Globs.GinEngine = engine

	defer Globs.GormDB.Close()
	defer Globs.MgoDB.Close()

	retCode := m.Run()
	os.Exit(retCode)
}
