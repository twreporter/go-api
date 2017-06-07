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

type TopicsResponse struct {
	Status  string         `json:"status"`
	Records []models.Topic `json:"records"`
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
	assert.Equal(t, topic.ID, TopicCol.ID)
	assert.Equal(t, topic.LeadingImage.ID, ImgID1)
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
	assert.Equal(t, topic.ID, TopicCol.ID)
	assert.Equal(t, topic.LeadingImage.ID, ImgID1)
	assert.Equal(t, topic.LeadingVideo.ID, VideoID)
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
