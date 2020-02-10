package controllers

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"twreporter.org/go-api/globals"
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
		log.Errorf("%+v", errors.WithStack(err))
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
