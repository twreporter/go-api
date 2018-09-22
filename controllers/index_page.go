package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// IndexPageQueryStruct struct for querying.
type IndexPageQueryStruct struct {
	MongoQuery   models.MongoQuery
	Limit        int
	Offset       int
	Sort         string
	Embedded     []string
	ResourceType string
	Full         bool
}

// _GetIndexPageContent ...
func (nc *NewsController) _GetIndexPageContent(part IndexPageQueryStruct) (interface{}, error) {
	var entities interface{}
	var err error
	if part.ResourceType == "topics" {
		if part.Full {
			entities, _, err = nc.Storage.GetFullTopics(part.MongoQuery, part.Limit, part.Offset, part.Sort, nil)
		} else {
			entities, _, err = nc.Storage.GetMetaOfTopics(part.MongoQuery, part.Limit, part.Offset, part.Sort, nil)
		}
	} else {
		if part.Full {
			entities, _, err = nc.Storage.GetFullPosts(part.MongoQuery, part.Limit, part.Offset, part.Sort, nil)
		} else {
			entities, _, err = nc.Storage.GetMetaOfPosts(part.MongoQuery, part.Limit, part.Offset, part.Sort, nil)
		}
	}

	if err != nil {
		return nil, err
	}

	return entities, nil
}

// _GetContentConcurrently ...
func (nc *NewsController) _GetContentConcurrently(parts map[string]IndexPageQueryStruct) map[string]interface{} {
	var ch = make(chan map[string]interface{}, len(parts))
	var rtn = make(map[string]interface{})

	for name, part := range parts {
		// concurrently get the sections
		go func(name string, part IndexPageQueryStruct) {
			entities, err := nc._GetIndexPageContent(part)
			if err == nil {
				ch <- map[string]interface{}{name: entities}
			}
		}(name, part)
	}

	for i := 0; i < len(parts); i++ {
		select {
		// read the section content from channel
		case section := <-ch:
			for k, v := range section {
				rtn[k] = v
			}
		case <-time.After(configs.TimeoutOfIndexPageController * time.Second):
			log.Info("The requests for fetching section timeouts")
		}
	}

	return rtn
}

// GetIndexPageContents is specifically made for index page.
// It will return the first fourth sections including
// latest, editor picks, latest topic and reviews.
func (nc *NewsController) GetIndexPageContents(c *gin.Context) {
	var rtn map[string]interface{}
	var ch = make(chan map[string]interface{})
	var parts = map[string]IndexPageQueryStruct{
		globals.LastestSection: {
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:        6,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		globals.EditorPicksSection: {
			MongoQuery: models.MongoQuery{
				State:      "published",
				IsFeatured: true,
			},
			Limit:        6,
			Offset:       0,
			Sort:         "-updatedAt",
			ResourceType: "posts",
		},
		globals.LatestTopicSection: {
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:        1,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "topics",
			Full:         true,
		},
		globals.ReviewsSection: {
			MongoQuery: models.MongoQuery{
				State: "published",
				Style: "review",
			},
			Limit:        4,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		globals.TopicsSection: {
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:        4,
			Offset:       1,
			Sort:         "-publishedDate",
			ResourceType: "topics",
			Full:         false,
		},
		globals.PhotoSection: {
			MongoQuery: models.MongoQuery{
				State: "published",
				Style: "photography",
			},
			Limit:        6,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		globals.InfographicSection: {
			MongoQuery: models.MongoQuery{
				State: "published",
				Style: "interactive",
			},
			Limit:        6,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
	}

	go func(parts map[string]IndexPageQueryStruct) {
		ch <- nc._GetContentConcurrently(parts)
	}(parts)
	select {
	// read the section content from channel
	case rtn = <-ch:
	case <-time.After(configs.TimeoutOfIndexPageController * time.Second):
		log.Info("The requests for fetching sections index page needed timeouts")
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": rtn})
}

// GetCategoriesPosts is specifically made for index page.
// It will return the posts of all the categories.
func (nc *NewsController) GetCategoriesPosts(c *gin.Context) {
	var rtn map[string]interface{}
	var ch = make(chan map[string]interface{})
	var cats = map[string]string{
		globals.HumanRightsAndSociety:   configs.HumanRightsAndSocietyListID,
		globals.EnvironmentAndEducation: configs.EnvironmentAndEducationListID,
		globals.PoliticsAndEconomy:      configs.PoliticsAndEconomyListID,
		globals.CultureAndArt:           configs.CultureAndArtListID,
		globals.International:           configs.InternationalListID,
		globals.LivingAndMedicalCare:    configs.LivingAndMedicalCareListID,
	}

	var parts = make(map[string]IndexPageQueryStruct)

	for name, ID := range cats {
		parts[name] = IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
				Categories: models.MongoQueryComparison{
					In: []bson.ObjectId{
						bson.ObjectIdHex(ID),
					},
				},
			},
			Limit:        5,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		}
	}

	go func(parts map[string]IndexPageQueryStruct) {
		ch <- nc._GetContentConcurrently(parts)
	}(parts)
	select {
	// read the section content from channel
	case rtn = <-ch:
	case <-time.After(configs.TimeoutOfIndexPageController * time.Second):
		log.Info("The requests for fetching sections index page needed timeouts")
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": rtn})
}
