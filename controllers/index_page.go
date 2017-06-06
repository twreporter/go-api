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
		entities, err = nc.Storage.GetTopics(part.MongoQuery, part.Limit, part.Offset, part.Sort, part.Embedded)
	} else {
		entities, err = nc.Storage.GetMetaOfPosts(part.MongoQuery, part.Limit, part.Offset, part.Sort, part.Embedded)
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
	var ch chan map[string]interface{} = make(chan map[string]interface{}, 1)
	var rtn = make(map[string]interface{})

	var parts = map[string]IndexPageQueryStruct{
		constants.LastestSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:    6,
			Offset:   0,
			Sort:     "-publishedDate",
			Embedded: []string{"hero_image", "categories", "topic_meta", "og_image"},
			//Fn:       nc.Storage.GetMetaOfPosts,
			ResourceType: "posts",
		},
		constants.EditorPicksSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State:      "published",
				IsFeatured: true,
			},
			Limit:    6,
			Offset:   0,
			Sort:     "-publishedDate",
			Embedded: []string{"hero_image", "categories", "topic_meta", "og_image"},
			//Fn:       nc.Storage.GetMetaOfPosts,
			ResourceType: "posts",
		},
		constants.LatestTopicSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
			},
			Limit:    1,
			Offset:   0,
			Sort:     "-publishedDate",
			Embedded: []string{"relateds_meta", "leading_image", "leading_image_portrait", "og_image"},
			//Fn:       nc.Storage.GetTopics,
			ResourceType: "topics",
		},
		constants.ReviewsSection: IndexPageQueryStruct{
			MongoQuery: models.MongoQuery{
				State: "published",
				Style: "review",
			},
			Limit:    4,
			Offset:   0,
			Sort:     "-publishedDate",
			Embedded: []string{"hero_image", "og_image"},
			//Fn:       nc.Storage.GetMetaOfPosts,
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
