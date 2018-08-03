package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	//log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/constants"
)

type indexPageResponse struct {
	Status  string                   `json:"status"`
	Records map[string][]interface{} `json:"records"`
}

func TestIndexPage(t *testing.T) {

	var resp *httptest.ResponseRecorder

	// Start -- Get four sections in the index page first screen //
	resp = serveHTTP("GET", "/v1/index_page", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := indexPageResponse{}
	json.Unmarshal(body, &res)

	latest, ok1 := res.Records[constants.LastestSection]
	assert.True(t, ok1)
	assert.Equal(t, len(latest), 2)

	picks, ok2 := res.Records[constants.EditorPicksSection]
	assert.True(t, ok2)
	assert.Equal(t, len(picks), 1)

	topic, ok3 := res.Records[constants.LatestTopicSection]
	assert.True(t, ok3)
	assert.Equal(t, len(topic), 1)

	reviews, ok4 := res.Records[constants.ReviewsSection]
	assert.True(t, ok4)
	assert.Equal(t, len(reviews), 1)
	// End -- Get four sections in the index page first screen //
}
