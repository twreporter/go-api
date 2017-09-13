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
	ID         int    `json:"ID"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	Host       string `json:"host"`
	IsExternal bool   `json:"is_external"`
	Desc       string `json:"desc"`
	Thumbnail  string `json:"thumbnail"`
	Category   string `json:"category"`
	Authors    string `json:"authors"`
	PubDate    string `json:"published_date"`
}

type Meta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type Response struct {
	Meta      Meta           `json:"meta"`
	Status    string         `json:"status"`
	Bookmarks []BookmarkJSON `json:"records"`
	Bookmark  BookmarkJSON   `json:"record"`
}

var bookmarkJSON = `{"slug":"mock-article-1","title":"mock title 1","host":"www.twreporter.org","is_external":false,"desc": "mock desc 1","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`
var bookmarkJSON2 = `{"slug":"mock-article-2","title":"mock title 2","host":"www.twreporter.org","is_external":false,"desc": "mock desc 2","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`
var bookmarkJSON3 = `{"slug":"mock-article-3","title":"mock title 3","host":"www.twreporter.org","is_external":false,"desc": "mock desc 3","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`

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
	var badBookmarkJSON = `{"slug":"bad-mock-article","title":"mock title","host":"www.twreporter.org","is_external":false,"desc": "mock desc","thumbnail":"www.twreporter.org/images/}`

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
	var path = fmt.Sprintf("/v1/users/%v/bookmarks?offset=0", DefaultID2)
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

	assert.Equal(t, res.Meta.Limit, 10)
	assert.Equal(t, res.Meta.Offset, 0)
	assert.Equal(t, res.Meta.Total, 1)
	assert.Equal(t, res.Bookmarks[0].Slug, "mock-article-1")
	assert.Equal(t, res.Bookmarks[0].IsExternal, false)
	assert.Equal(t, res.Bookmarks[0].Category, "")
	/** END - List bookmarks successfully **/

	/** START - Fail to list bookmark **/
	// user is not existed
	var fakeID uint = 100
	jwt, _ = utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	resp = ServeHTTP("GET", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to list bookmark **/
}

func TestGetABookmarkOfAUser(t *testing.T) {
	var path = fmt.Sprintf("/v1/users/%v/bookmarks/mock-article-3", DefaultID2)
	var jwt = GenerateJWT(GetUser(DefaultID2))

	/** START - Fail to get a bookmark of a user **/
	resp := ServeHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)

	// add a bookmark onto a user
	_ = ServeHTTP("POST", fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID2), bookmarkJSON3, "application/json", fmt.Sprintf("Bearer %v", jwt))

	// still fail to get the bookmark of the user because of host is not provided
	resp = ServeHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)

	/** END - Fail to get a bookmark of a user **/

	/** START - get a bookmark of a user **/
	// add host param
	path = path + "?host=www.twreporter.org"
	// get the bookmark of the user
	resp = ServeHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := Response{}
	json.Unmarshal(body, &res)

	assert.Equal(t, res.Bookmark.Slug, "mock-article-3")
	/** END - get a bookmark of user **/
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
