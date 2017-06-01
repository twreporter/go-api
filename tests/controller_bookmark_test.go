package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/utils"
)

type BookmarkJSON struct {
	ID        int    `json:"ID"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
	Href      string `json:"Href"`
	Title     string `json:"Title"`
	Desc      struct {
		String string `json:"String"`
		Valid  bool   `json:"Valid"`
	} `json:"Desc"`
	Thumbnail struct {
		String string `json:"String"`
		Valid  bool   `json:"Valid"`
	} `json:"Thumbnail"`
}

type Response struct {
	Status    string         `json:"status"`
	Bookmarks []BookmarkJSON `json:"records"`
}

func TestBookmarkAuthorization(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var path = fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID)

	/** START - Fail to pass Authorization **/
	// without Authorization header
	resp = ServeHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 401)

	// wrong jwt in Authorization header
	resp = ServeHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 401)
	/** END - Fail to pass Authorization **/

	// Pass Authroization
	resp = ServeHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 200)
}

func TestCreateABookmarkOfAUser(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var badBookmarkJSON = `{"href":"www.twreporter.org/a/a-mock-article","title":"mock title","desc": "mock desc","thumbnail":"www.twreporter.org/images/}`

	var bookmarkJSON = `{"href":"www.twreporter.org/a/a-mock-article","title":"mock title","desc": "mock desc","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`
	var bookmarkJSON2 = `{"href":"www.twreporter.org/a/another-mock-article","title":"another mock title","desc": "another mock desc","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`

	var path = fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID)
	var jwt = GenerateJWT(GetUser(DefaultID))

	/** START - Add bookmark successfully **/
	resp = ServeHTTP("POST", path, bookmarkJSON, "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 201)

	// add another bookmark
	resp = ServeHTTP("POST", path, bookmarkJSON2, "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 201)
	/** END - Add bookmark successfully **/

	/** START - Fail to add bookmark **/
	// malformed JSON
	resp = ServeHTTP("POST", path, badBookmarkJSON, "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 400)

	// user is not existed
	var fakeID uint = 100
	jwt, _ = utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	resp = ServeHTTP("POST", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), bookmarkJSON,
		"application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to add bookmark **/
}

func TestGetBookmarksOfAUser(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var bookmarkJSON = `{"href":"www.twreporter.org/a/a-mock-article","title":"mock title","desc": "mock desc","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`
	var path = fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID2)
	var jwt = GenerateJWT(GetUser(DefaultID2))

	/** START - List bookmarks successfully **/
	// List empty array of bookmarks of the user
	resp = ServeHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 200)

	// List non-empty array of bookmarks of the user
	// add a bookmark into the user
	_ = ServeHTTP("POST", path, bookmarkJSON, "application/json", fmt.Sprintf("Bearer %v", jwt))

	// get bookmarks of the user
	resp = ServeHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := Response{}
	json.Unmarshal(body, &res)

	assert.Equal(t, res.Bookmarks[0].ID, 1)
	assert.Equal(t, res.Bookmarks[0].Href, "www.twreporter.org/a/a-mock-article")
	/** END - List bookmarks successfully **/

	/** START - Fail to list bookmark **/
	// user is not existed
	var fakeID uint = 100
	jwt, _ = utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	resp = ServeHTTP("GET", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to list bookmark **/
}

func TestDeleteBookmark(t *testing.T) {
	var resp *httptest.ResponseRecorder

	/** START - Delete bookmark successfully **/
	resp = ServeHTTP("DELETE", fmt.Sprintf("/v1/users/%v/bookmarks/1", DefaultID), "", "", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 204)
	/** END - Delete bookmark successfully **/

	/** START - Fail to delete bookmark **/
	// delete the bookmark again
	resp = ServeHTTP("DELETE", fmt.Sprintf("/v1/users/%v/bookmarks/1", DefaultID), "", "", fmt.Sprintf("Bearer %v", GenerateJWT(GetUser(DefaultID))))
	assert.Equal(t, resp.Code, 404)

	// user is not existed
	var fakeID uint = 100
	jwt, _ := utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	resp = ServeHTTP("DELETE", fmt.Sprintf("/v1/users/%v/bookmarks/1", fakeID), "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to list bookmark **/
}
