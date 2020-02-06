package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"twreporter.org/go-api/models"
)

// GetAuthors receive HTTP GET method request, and return the authors.
// `limit`, `offset` and `sort` are the url query params,
// which define the rule we retrieve authors from storage.
func (nc *NewsController) GetAuthors(c *gin.Context) (int, gin.H, error) {
	const defaultLimit = 20
	const defaultSort = "updatedAt"
	var authors []models.FullAuthor
	var err error
	var total int

	_, _, limit, offset, sort, _ := nc.GetQueryParam(c)

	if limit == 0 {
		limit = defaultLimit
	}

	if sort == "" {
		sort = defaultSort
	}

	authors, total, err = nc.Storage.GetFullAuthors(limit, offset, sort)

	if err != nil {
		return 0, gin.H{}, err
	}

	return http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"records": authors,
			"meta": models.MetaOfResponse{
				Total:  total,
				Offset: offset,
				Limit:  limit,
			},
		},
	}, nil
}
