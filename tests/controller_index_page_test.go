package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/globals"
)

type indexPageResponse struct {
	Status  string                   `json:"status"`
	Records map[string][]interface{} `json:"records"`
}

func TestIndexPage(t *testing.T) {
	var resp *httptest.ResponseRecorder

	const defaultLastestNum = 2
	const defaultPickNum = 1
	const defaultTopicNum = 1
	const defaultReviewNum = 1

	// Start -- Get four sections in the index page first screen //
	resp = serveHTTP("GET", "/v1/index_page", "",
		"", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := indexPageResponse{}
	json.Unmarshal(body, &res)

	latest, ok1 := res.Records[globals.LastestSection]
	assert.True(t, ok1)
	assert.Equal(t, defaultLastestNum, len(latest))

	picks, ok2 := res.Records[globals.EditorPicksSection]
	assert.True(t, ok2)
	assert.Equal(t, defaultPickNum, len(picks))

	topic, ok3 := res.Records[globals.LatestTopicSection]
	assert.True(t, ok3)
	assert.Equal(t, defaultTopicNum, len(topic))

	reviews, ok4 := res.Records[globals.ReviewsSection]
	assert.True(t, ok4)
	assert.Equal(t, defaultReviewNum, len(reviews))
	// End -- Get four sections in the index page first screen //
}
