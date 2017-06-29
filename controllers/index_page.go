package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
}

func (nc *NewsController) __GetIndexPageContent(part IndexPageQueryStruct) (interface{}, error) {
	var entities interface{}
	var err error
	if part.ResourceType == "topics" {
		entities, _, err = nc.Storage.GetFullTopics(part.MongoQuery, part.Limit, part.Offset, part.Sort, nil)
	} else {
		entities, _, err = nc.Storage.GetMetaOfPosts(part.MongoQuery, part.Limit, part.Offset, part.Sort, nil)
	}

	if err != nil {
		log.Info("err:", err)
		return nil, err
	}

	return entities, nil
}

// GetIndexPageContents is specifically made for index page.
// It will return the first fourth sections including
// latest, editor picks, latest topic and reviews.
func (nc *NewsController) GetIndexPageContents(c *gin.Context) {
	var ch chan map[string]interface{} = make(chan map[string]interface{})
	var rtn = make(map[string]interface{})

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
	}

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

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": rtn})
}
