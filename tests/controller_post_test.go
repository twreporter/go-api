package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	//log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
)

/*
The whole testing mongodb is set up by `setMgoDefaultRecords` function in `pre_test_environment_setup.go`
*/

type postsResponse struct {
	Status  string        `json:"status"`
	Records []models.Post `json:"records"`
}

type postResponse struct {
	Status string      `json:"status"`
	Record models.Post `json:"record"`
}

func TestGetAPost(t *testing.T) {
	// Post Not Found //
	resp := serveHTTP("GET", "/v1/posts/post-not-found", "",
		"", "")
	assert.Equal(t, resp.Code, 404)
	// Post Not Found //

	// Get a post without full url param //
	resp = serveHTTP("GET", "/v1/posts/"+Globs.Defaults.MockPostSlug1, "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := postResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Record.ID, Globs.Defaults.PostID1)
	assert.Equal(t, len(res.Record.Relateds), 0)
	assert.Equal(t, res.Record.Full, false)
	// Get a post without full url param //

	// Get a post with full url param //
	resp = serveHTTP("GET", "/v1/posts/"+Globs.Defaults.MockPostSlug1+"?full=true", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, res.Record.ID, Globs.Defaults.PostID1)
	assert.Equal(t, len(res.Record.Relateds), 1)
	assert.Equal(t, res.Record.Relateds[0].ID, Globs.Defaults.PostID2)
	assert.Equal(t, res.Record.Full, true)
	// Get a post with full url param //
}

func TestGetPosts(t *testing.T) {

	var resp *httptest.ResponseRecorder

	// Start -- Get all the posts //
	resp = serveHTTP("GET", "/v1/posts", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 2)

	post := res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol2.ID)
	assert.Equal(t, post.OgImage.ID, Globs.Defaults.ImgID2)
	assert.Equal(t, post.HeroImage.ID, Globs.Defaults.ImgID2)
	assert.Equal(t, post.Topic.ID, Globs.Defaults.TopicID)
	assert.Equal(t, len(post.Tags), 1)
	assert.Equal(t, len(post.Categories), 1)
	assert.Equal(t, post.Tags[0].ID, Globs.Defaults.TagID)
	assert.Equal(t, post.Categories[0].ID, Globs.Defaults.CatID)
	assert.Equal(t, post.IsFeatured, false)
	assert.Equal(t, post.Full, false)

	post = res.Records[1]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol1.ID)
	assert.Equal(t, post.OgImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, post.Topic.ID, Globs.Defaults.TopicID)
	assert.Equal(t, len(post.Tags), 0)
	assert.Equal(t, len(post.Categories), 1)
	assert.Equal(t, post.Categories[0].ID, Globs.Defaults.CatID)
	assert.Equal(t, post.IsFeatured, true)
	assert.Equal(t, post.Full, false)
	assert.Equal(t, post.Theme.ID, Globs.Defaults.ThemeID)
	// End -- Get all the posts //

	// Start -- Get posts with isFeature=true //
	resp = serveHTTP("GET", "/v1/posts?where={\"is_featured\":true}", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)

	post = res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol1.ID)
	// End -- Get posts with isFeature=true //

	// Start -- Get posts with style=review //
	resp = serveHTTP("GET", "/v1/posts?where={\"style\":\"review\"}", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)

	post = res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol2.ID)
	// End -- Get posts with style=review //

	// Start -- Get posts with slug=mock-post-slug-2 //
	resp = serveHTTP("GET", "/v1/posts?where={\"slug\":\"mock-post-slug-2\"}", "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)

	post = res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol2.ID)
	// End -- Get posts with slug=mock-post-slug-2 //

	// Start -- Get posts containing TagID //
	resp = serveHTTP("GET", fmt.Sprintf("/v1/posts?where={\"tags\":{\"in\":[\"%v\"]}}", Globs.Defaults.TagID.Hex()), "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 1)

	post = res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol2.ID)
	// End -- Get posts containing TagID //

	// Start -- Get posts containing CatID //
	resp = serveHTTP("GET", fmt.Sprintf("/v1/posts?where={\"postcategories\":{\"in\":[\"%v\"]}}", Globs.Defaults.CatID.Hex()), "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, len(res.Records), 2)
	// End -- Get posts containing CatID //
}
