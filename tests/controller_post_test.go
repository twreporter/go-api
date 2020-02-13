package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/configs"
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
	assert.Equal(t, http.StatusNotFound, resp.Code)
	// Post Not Found //

	// Get a post without full url param //
	resp = serveHTTP("GET", "/v1/posts/"+Globs.Defaults.MockPostSlug1, "",
		"", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := postResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, Globs.Defaults.PostID1, res.Record.ID)
	assert.Equal(t, 0, len(res.Record.Relateds))
	assert.Equal(t, false, res.Record.Full)
	// Get a post without full url param //

	// Get a post with full url param //
	resp = serveHTTP("GET", "/v1/posts/"+Globs.Defaults.MockPostSlug1+"?full=true", "",
		"", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, Globs.Defaults.PostID1, res.Record.ID)
	assert.Equal(t, 1, len(res.Record.Relateds))
	assert.Equal(t, Globs.Defaults.PostID2, res.Record.Relateds[0].ID)
	assert.Equal(t, true, res.Record.Full)
	// Get a post with full url param //
}

func TestGetPosts(t *testing.T) {

	var resp *httptest.ResponseRecorder

	// Start -- Get all the posts //
	resp = serveHTTP("GET", "/v1/posts", "",
		"", "")
	assert.Equal(t, http.StatusOK, resp.Code)
	body, _ := ioutil.ReadAll(resp.Result().Body)
	res := postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, 2, len(res.Records))

	post := res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol2.ID)
	assert.Equal(t, post.OgImage.ID, Globs.Defaults.ImgID2)
	assert.Equal(t, post.HeroImage.ID, Globs.Defaults.ImgID2)
	assert.Equal(t, post.Topic.ID, Globs.Defaults.TopicID)
	assert.Equal(t, len(post.Tags), 1)
	assert.Equal(t, 1, len(post.Categories))
	assert.Equal(t, post.Tags[0].ID, Globs.Defaults.TagID)
	assert.Equal(t, Globs.Defaults.CatReviewID, post.Categories[0].ID)
	assert.Equal(t, post.IsFeatured, false)
	assert.Equal(t, post.Full, false)

	post = res.Records[1]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol1.ID)
	assert.Equal(t, post.OgImage.ID, Globs.Defaults.ImgID1)
	assert.Equal(t, post.Topic.ID, Globs.Defaults.TopicID)
	assert.Equal(t, len(post.Tags), 0)
	assert.Equal(t, 1, len(post.Categories))
	assert.Equal(t, Globs.Defaults.CatPhotographyID, post.Categories[0].ID)
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

	// Start -- Get posts with Review category //
	resp = serveHTTP("GET", fmt.Sprintf("/v1/posts?where={\"categories\":{\"in\":[\"%v\"]}}", configs.ReviewListID), "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, 1, len(res.Records))

	post = res.Records[0]
	assert.Equal(t, post.ID, Globs.Defaults.PostCol2.ID)
	// End -- Get posts with style=review //

	// Start -- Get posts with Photography category //
	resp = serveHTTP("GET", fmt.Sprintf("/v1/posts?where={\"categories\":{\"in\":[\"%v\"]}}", configs.PhotographyListID), "",
		"", "")
	assert.Equal(t, resp.Code, 200)
	body, _ = ioutil.ReadAll(resp.Result().Body)
	res = postsResponse{}
	json.Unmarshal(body, &res)
	assert.Equal(t, 1, len(res.Records))
	// End -- Get posts containing Photography category //

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
}
