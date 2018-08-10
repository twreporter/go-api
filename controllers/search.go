package controllers

import (
	"net/http"
	"strconv"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/utils"
)

// _Search - search records from algolia webservice
func (nc *NewsController) _Search(c *gin.Context, indexName string) {
	var err error
	var hitsPerPage int
	var page int
	var res algoliasearch.QueryRes

	filters := c.Query("filters")
	hitsPerPage, err = strconv.Atoi(c.Query("hitsPerPage"))
	page, err = strconv.Atoi(c.Query("page"))
	keywords := c.Query("keywords")

	client := algoliasearch.NewClient(utils.Cfg.AlgoliaSettings.ApplicationID, utils.Cfg.AlgoliaSettings.APIKey)
	index := client.InitIndex(indexName)

	res, err = index.Search(keywords, algoliasearch.Map{
		"filters":     filters,
		"hitsPerPage": hitsPerPage,
		"page":        page,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal server error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// SearchAuthors - search authors from algolia webservice
func (nc *NewsController) SearchAuthors(c *gin.Context) {
	nc._Search(c, "contacts-index-v2")
}

// SearchPosts - search posts of authors from algolia webservice
func (nc *NewsController) SearchPosts(c *gin.Context) {
	nc._Search(c, "posts-index-v2")
}
