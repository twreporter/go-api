package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"

	log "github.com/Sirupsen/logrus"
)

type bookmarkJSON struct {
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

type meta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type Response struct {
	Meta      meta           `json:"meta"`
	Status    string         `json:"status"`
	Bookmarks []bookmarkJSON `json:"records"`
	Bookmark  bookmarkJSON   `json:"record"`
}

var bookmarkJSONStr = `{"slug":"mock-article-1","title":"mock title 1","host":"www.twreporter.org","is_external":false,"desc": "mock desc 1","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`
var bookmarkJSONStr2 = `{"slug":"mock-article-2","title":"mock title 2","host":"www.twreporter.org","is_external":false,"desc": "mock desc 2","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`

func TestBookmarkAuthorization(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var user models.User
	var path string

	user = getUser(Globs.Defaults.Account)
	path = fmt.Sprintf("/v1/users/%v/bookmarks", user.ID)

	/** START - Fail to pass Authorization **/
	// without Authorization header
	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 401)

	// wrong jwt in Authorization header
	resp = serveHTTP("GET", path, "", "", "")
	assert.Equal(t, resp.Code, 401)
	/** END - Fail to pass Authorization **/

	// Pass Authroization
	resp = serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", generateJWT(user)))
	assert.Equal(t, resp.Code, 200)
}

func TestCreateABookmarkOfAUser(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var badBookmarkJSON = `{"slug":"bad-mock-article","title":"mock title","host":"www.twreporter.org","is_external":false,"desc": "mock desc","thumbnail":"www.twreporter.org/images/}`
	var user models.User
	var path string

	user = getUser(Globs.Defaults.Account)
	path = fmt.Sprintf("/v1/users/%v/bookmarks", user.ID)

	var jwt = generateJWT(user)

	/** START - Add bookmark successfully **/
	resp = serveHTTP("POST", path, bookmarkJSONStr, "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 201)

	// add another bookmark
	resp = serveHTTP("POST", path, bookmarkJSONStr2, "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 201)
	/** END - Add bookmark successfully **/

	/** START - Fail to add bookmark **/
	// malformed JSON
	resp = serveHTTP("POST", path, badBookmarkJSON, "application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 400)

	// user is not existed
	var fakeID uint = 100
	jwt, _ = utils.RetrieveV1Token(fakeID, "test@twreporter.org")
	log.Info("jwt:", jwt)
	resp = serveHTTP("POST", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), bookmarkJSONStr,
		"application/json", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to add bookmark **/
}

func TestGetBookmarksOfAUser(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var user models.User
	var path string
	var jwt string

	user = getUser(Globs.Defaults.Account)
	path = fmt.Sprintf("/v1/users/%v/bookmarks?offset=0", user.ID)
	jwt = generateJWT(user)

	/** START - List bookmarks successfully **/
	// List empty array of bookmarks of the user
	resp = serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 200)

	// List non-empty array of bookmarks of the user
	// add a bookmark into the user
	_ = serveHTTP("POST", path, bookmarkJSONStr, "application/json", fmt.Sprintf("Bearer %v", jwt))
	// get bookmarks of the user
	resp = serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := Response{}
	json.Unmarshal(body, &res)

	assert.Equal(t, res.Meta.Limit, 10)
	assert.Equal(t, res.Meta.Offset, 0)
	assert.NotZero(t, res.Meta.Total)
	/** END - List bookmarks successfully **/

	/** START - Fail to list bookmark **/
	// user is not existed
	var fakeID uint = 100
	jwt, _ = utils.RetrieveV1Token(fakeID, "test@twreporter.org")

	resp = serveHTTP("GET", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to list bookmark **/
}

func TestGetABookmarkOfAUser(t *testing.T) {
	var user models.User
	var path string
	var jwt string

	user = getUser(Globs.Defaults.Account)
	jwt = generateJWT(user)
	path = fmt.Sprintf("/v1/users/%v/bookmarks/mock-article-2", user.ID)

	/** START - Fail to get a bookmark of a user **/
	resp := serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)

	// add a bookmark onto a user
	_ = serveHTTP("POST", fmt.Sprintf("/v1/users/%v/bookmarks", user.ID), bookmarkJSONStr2, "application/json", fmt.Sprintf("Bearer %v", jwt))

	// still fail to get the bookmark of the user because of host is not provided
	resp = serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 404)

	/** END - Fail to get a bookmark of a user **/

	/** START - get a bookmark of a user **/
	// add host param
	path = path + "?host=www.twreporter.org"
	// get the bookmark of the user
	resp = serveHTTP("GET", path, "", "", fmt.Sprintf("Bearer %v", jwt))
	assert.Equal(t, resp.Code, 200)

	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := Response{}
	json.Unmarshal(body, &res)

	assert.Equal(t, res.Bookmark.Slug, "mock-article-2")
	/** END - get a bookmark of user **/
}

func TestDeleteBookmark(t *testing.T) {
	var resp *httptest.ResponseRecorder
	var user models.User

	user = getUser(Globs.Defaults.Account)

	type deleteBookMarkInfo struct {
		User       models.User
		BookMarkID uint
	}

	const mockBookMarkID = 1
	const invalidBookMarkID = 1000
	const unknownUserID = 100
	const unknownUserEmail = "test@twreporter.org"
	cases := []struct {
		name         string
		mockBookMark *models.Bookmark
		deleteInfo   deleteBookMarkInfo
		cleanupStmt  string
		respCode     int
	}{
		{
			name: "StatusCode=StatusNotFound,Invalid Bookmark ID",
			mockBookMark: &models.Bookmark{
				ID:         mockBookMarkID,
				Slug:       "mock-article-3",
				Title:      "mock title 3",
				Host:       "www.twreporter.org",
				IsExternal: false,
				Desc:       "mock desc 3",
				Thumbnail:  "www.twreporter.org/images/mock-image.jpg",
				Users: []models.User{
					user,
				},
			},
			deleteInfo: deleteBookMarkInfo{
				User:       user,
				BookMarkID: invalidBookMarkID,
			},
			cleanupStmt: fmt.Sprintf("DELETE FROM bookmarks where id = '%d';", mockBookMarkID),
			respCode:    http.StatusNotFound,
		},
		{
			name:         "StatusCode=StatusNotFound,Unknown User",
			mockBookMark: nil,
			deleteInfo: deleteBookMarkInfo{
				User: models.User{
					ID:    unknownUserID,
					Email: null.StringFrom(unknownUserEmail),
				},
				BookMarkID: mockBookMarkID,
			},
			cleanupStmt: "",
			respCode:    http.StatusNotFound,
		},
		{
			name: "StatusCode=StatusNoContent,Invalid Bookmark ID",
			mockBookMark: &models.Bookmark{
				ID:         mockBookMarkID,
				Slug:       "mock-article-3",
				Title:      "mock title 3",
				Host:       "www.twreporter.org",
				IsExternal: false,
				Desc:       "mock desc 3",
				Thumbnail:  "www.twreporter.org/images/mock-image.jpg",
				Users: []models.User{
					user,
				},
			},
			deleteInfo: deleteBookMarkInfo{
				User:       user,
				BookMarkID: mockBookMarkID,
			},
			cleanupStmt: fmt.Sprintf("DELETE FROM bookmarks where id = '%d';", mockBookMarkID),
			respCode:    http.StatusNoContent,
		},
	}

	db := Globs.GormDB
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			jwt := generateJWT(tc.deleteInfo.User)

			if tc.mockBookMark != nil {
				reqBody, _ := json.Marshal(tc.mockBookMark)
				serveHTTP("POST", fmt.Sprintf("/v1/users/%d/bookmarks", tc.mockBookMark.Users[0].ID), string(reqBody), "application/json", fmt.Sprintf("Bearer %v", jwt))
			}

			if tc.cleanupStmt != "" {
				defer func() {
					db.Exec(tc.cleanupStmt)
				}()
			}
			resp = serveHTTP("DELETE", fmt.Sprintf("/v1/users/%d/bookmarks/%d", tc.deleteInfo.User.ID, tc.deleteInfo.BookMarkID), "", "", fmt.Sprintf("Bearer %v", jwt))
			assert.Equal(t, tc.respCode, resp.Code)
		})
	}
}
