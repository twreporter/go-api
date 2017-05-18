package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/storage"
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
	Bookmarks []BookmarkJSON `json:"bookmarks"`
}

func generateJWT(user models.User) (jwt string) {
	jwt, _ = utils.RetrieveToken(user.ID, user.Privilege, user.FirstName.String, user.LastName.String, user.Email.String)
	return
}

func getUser(userId string) (user models.User) {
	as := storage.NewGormUserStorage(DB)
	user, _ = as.GetUserByID(userId)
	return
}

func TestAuthorization(t *testing.T) {
	var path = fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID)

	/** START - Fail to pass Authorization **/
	// without Authorization header
	req, _ := http.NewRequest("GET", path, nil)
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 401)

	// wrong jwt in Authorization header
	req, _ = http.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", generateJWT(getUser(DefaultID2))))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 401)
	/** END - Fail to pass Authorization **/

	// Pass Authroization
	req, _ = http.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", generateJWT(getUser(DefaultID))))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
}

func TestCreateBookmarkByUser(t *testing.T) {
	var badBookmarkJSON = `{"href":"www.twreporter.org/a/a-mock-article","title":"mock title","desc": "mock desc","thumbnail":"www.twreporter.org/images/}`

	var bookmarkJSON = `{"href":"www.twreporter.org/a/a-mock-article","title":"mock title","desc": "mock desc","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`

	var path = fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID)
	var jwt = generateJWT(getUser(DefaultID))

	/** START - Add bookmark successfully **/
	req := RequestWithBody("POST", path, bookmarkJSON)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 201)
	/** END - Add bookmark successfully **/

	/** START - Fail to ddd bookmark **/
	// malformed JSON
	req = RequestWithBody("POST", path, badBookmarkJSON)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 400)

	// user is not existed
	var fakeID uint = 100
	jwt, _ = utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	req = RequestWithBody("POST", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), bookmarkJSON)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to add bookmark **/
}

func TestListBookmarkByUser(t *testing.T) {
	var bookmarkJSON = `{"href":"www.twreporter.org/a/a-mock-article","title":"mock title","desc": "mock desc","thumbnail":"www.twreporter.org/images/mock-image.jpg"}`
	var path = fmt.Sprintf("/v1/users/%v/bookmarks", DefaultID2)
	var jwt = generateJWT(getUser(DefaultID2))

	/** START - List bookmarks successfully **/
	// List empty array of bookmarks of the user
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)

	// List non-empty array of bookmarks of the user
	// add a bookmark into the user
	req = RequestWithBody("POST", path, bookmarkJSON)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)

	// get bookmarks of the user
	req, _ = http.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	body, _ := ioutil.ReadAll(resp.Result().Body)

	res := Response{}
	json.Unmarshal(body, &res)

	assert.Equal(t, res.Bookmarks[0].ID, 1)
	assert.Equal(t, res.Bookmarks[0].Href, "www.twreporter.org/a/a-mock-article")
	assert.Equal(t, resp.Code, 200)
	/** END - List bookmarks successfully **/

	/** START - Fail to list bookmark **/
	// user is not existed
	var fakeID uint = 100
	req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/users/%v/bookmarks", fakeID), nil)
	jwt, _ = utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to list bookmark **/
}

func TestDeleteBookmark(t *testing.T) {
	/** START - Delete bookmark successfully **/
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/users/%v/bookmarks/1", DefaultID), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", generateJWT(getUser(DefaultID))))
	resp := httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 204)
	/** END - Delete bookmark successfully **/

	/** START - Fail to delete bookmark **/
	// user is not existed
	var fakeID uint = 100
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/users/%v/bookmarks/1", fakeID), nil)
	jwt, _ := utils.RetrieveToken(fakeID, 0, "", "", "test@twreporter.org")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
	resp = httptest.NewRecorder()
	Engine.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 404)
	/** END - Fail to list bookmark **/
}
