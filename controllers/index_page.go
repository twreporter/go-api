package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

type IndexPageQueryStruct struct {
	MongoQuery   models.MongoQuery
	Limit        int
	Offset       int
	Sort         string
	Embedded     []string
	ResourceType string
	Full         bool
}

func (nc *NewsController) __GetIndexPageContent(part IndexPageQueryStruct) (interface{}, error) {
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
		log.Info("err:", err)
		return nil, err
	}

	return entities, nil
}

func (nc *NewsController) __GetContentConcurrently(parts map[string]IndexPageQueryStruct) map[string]interface{} {
	var ch chan map[string]interface{} = make(chan map[string]interface{})
	var rtn = make(map[string]interface{})

	for name, part := range parts {
		// concurrently get the sections
		go func(name string, part IndexPageQueryStruct) {
			entities, err := nc.__GetIndexPageContent(part)
			if err == nil {
				ch <- map[string]interface{}{name: entities}
			}
		}(name, part)

		select {
		// read the section content from channel
		case section := <-ch:
			for k, v := range section {
				rtn[k] = v
			}
		// set timeout
		case <-time.After(3 * time.Second):
			close(ch)
			log.Info("The requests for fetching sections index page needed timeouts")
		}
	}

	return rtn
}

// GetIndexPageContents is specifically made for index page.
// It will return the first fourth sections including
// latest, editor picks, latest topic and reviews.
func (nc *NewsController) GetIndexPageContents(c *gin.Context) {
	var parts = map[string]IndexPageQueryStruct{
		constants.LastestSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:        6,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		constants.EditorPicksSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State:      "published",
				IsFeatured: true,
			},
			Limit:        6,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		constants.LatestTopicSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:        1,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "topics",
			Full:         true,
		},
		constants.ReviewsSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
				Style: "review",
			},
			Limit:        4,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		constants.TopicsSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:        4,
			Offset:       1,
			Sort:         "-publishedDate",
			ResourceType: "topics",
			Full:         false,
		},
		constants.PhotoSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
				Style: "photography",
			},
			Limit:        4,
			Offset:       0,
			Sort:         "-publishedDate",
			ResourceType: "posts",
		},
		constants.InfographicSection: IndexPageQueryStruct{
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

	rtn := nc.__GetContentConcurrently(parts)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": rtn})
}

// GetIndexPageContents is specifically made for index page.
// It will return the first fourth sections including
// latest, editor picks, latest topic and reviews.
func (nc *NewsController) GetCategoriesPosts(c *gin.Context) {
	var cats = map[string]string{
		constants.HumanRights:        configs.HumanRightsListID,
		constants.LandEnvironment:    configs.LandEnvironmentListID,
		constants.TransformedJustice: configs.TransformedJusticeListID,
		constants.CultureMovie:       configs.CultureMovieListID,
		constants.PhotoAudio:         configs.PhotoAudioListID,
		constants.International:      configs.InternationalListID,
		constants.Character:          configs.CharacterListID,
		constants.PoliticalSociety:   configs.PoliticalSocietyListID,
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

	rtn := nc.__GetContentConcurrently(parts)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": rtn})
}
