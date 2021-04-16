package controllers

import (
	"net/http"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"

	"github.com/twreporter/go-api/globals"
)

// search - search records from algolia webservice
func search(c *gin.Context, indexName string) {
	var err error
	var hitsPerPage int
	var page int
	var res algoliasearch.QueryRes

	filters := c.Query("filters")
	hitsPerPage, _ = strconv.Atoi(c.Query("hitsPerPage"))
	page, _ = strconv.Atoi(c.Query("page"))
	keywords := c.Query("keywords")

	client := algoliasearch.NewClient(globals.Conf.Algolia.ApplicationID, globals.Conf.Algolia.APIKey)
	index := client.InitIndex(indexName)

	res, err = index.Search(keywords, algoliasearch.Map{
		"filters":     filters,
		"hitsPerPage": hitsPerPage,
		"page":        page,
	})

	if err != nil {
		if globals.Conf.Environment == "development" {
			log.Errorf("%+v", errors.WithStack(err))
		} else {
			log.WithField("detail", err).Errorf("%s", f.FormatStack(errors.WithStack(err)))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
	}

	c.JSON(http.StatusOK, res)
}

// SearchAuthors - search authors from algolia webservice
func (nc *NewsController) SearchAuthors(c *gin.Context) {
	search(c, "contacts-index-v2")
}

// SearchPosts - search posts of authors from algolia webservice
func (nc *NewsController) SearchPosts(c *gin.Context) {
	search(c, "posts-index-v2")
}
