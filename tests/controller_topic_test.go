package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	//log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
)

/*
The whole testing mongodb is set by ./test.go
You should check #SetMgoDefaultRecords function,
if you want to know more about the data set in the testing mongodb
*/

type TopicsResponse struct {
	Status  string         `json:"status"`
	Records []models.Topic `json:"records"`
}

type TopicResponse struct {
	Status string       `json:"status"`
	Record models.Topic `json:"record"`
}

func TestGetATopic(t *testing.T) {
	// Post Not Found //
	resp := ServeHTTP("GET", "/v1/topics/post-not-found", "",
		"", "")
	assert.Equal(t, resp.Code, 404)
	// Post Not Found //

	// Get a post without full url param //
	resp = ServeHTTP("GET", "/v1/topics/"+Globs.Defaults.MockTopicSlug, "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := TopicResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Record.ID, Globs.Defaults.TopicID)
	assert.Equal(t, len(res.Record.Relateds), 0)
	assert.Equal(t, res.Record.LeadingImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, res.Record.OgImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, res.Record.Full, false)
	// Get a post without full url param //

	// Get a post with full url param //
	resp = ServeHTTP("GET", "/v1/topics/"+Globs.Defaults.MockTopicSlug+"?full=true", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = TopicResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Record.ID, Globs.Defaults.TopicID)
	assert.Equal(t, len(res.Record.Relateds), 2)
	assert.Equal(t, res.Record.Relateds[0].ID, Globs.Defaults.PostID1)
	assert.Equal(t, res.Record.Relateds[1].ID, Globs.Defaults.PostID2)
	assert.Equal(t, res.Record.LeadingImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, res.Record.LeadingVideo.ID, Globs.Defaults.VideoID)
	assert.Equal(t, res.Record.OgImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, res.Record.Full, true)
	// Get a post with full url param //
}

func TestGetTopics(t *testing.T) {

	var resp *httptest.ResponseRecorder

	// Start -- Get all the topics //
	resp = ServeHTTP("GET", "/v1/topics", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := TopicsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)

	topic := res.Records[0]
	assert.Equal(t, topic.ID, Globs.Defaults.TopicCol.ID)
	assert.Equal(t, topic.LeadingImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, len(topic.Relateds), 0)
	// End -- Get all the posts //

	// Start -- Get all the full topics //
	resp = ServeHTTP("GET", "/v1/topics?full=true", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = TopicsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)

	topic = res.Records[0]
	assert.Equal(t, topic.ID, Globs.Defaults.TopicCol.ID)
	assert.Equal(t, topic.LeadingImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, topic.LeadingVideo.ID, Globs.Defaults.VideoID)
	assert.Equal(t, len(topic.Relateds), 2)
	// End -- Get all the posts //

	// Start -- Get the topics with slug=mock-topic-slug//
	resp = ServeHTTP("GET", "/v1/topics?where={\"slug\":\"mock-topic-slug\"}", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = TopicsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)
	// End -- Get the topics with slug=mock-topic-slug//

	// Start -- Get no topics  //
	resp = ServeHTTP("GET", "/v1/topics?where={\"slug\":\"wrong-topic-slug\"}", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = TopicsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 0)
	// End -- Get the topics with slug=mock-topic-slug//
}
