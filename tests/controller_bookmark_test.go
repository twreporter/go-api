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
	"github.com/twreporter/go-api/models"
)

type meta struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type response struct {
	Meta      meta              `json:"meta"`
	Status    string            `json:"status"`
	Bookmarks []models.Bookmark `json:"records"`
	Bookmark  models.Bookmark   `json:"record"`
}

func TestBookmarkAuthorization(t *testing.T) {
	user := getUser(Globs.Defaults.Account)
	path := fmt.Sprintf("/v1/users/%v/bookmarks", user.ID)

	for _, tc := range []struct {
		name       string
		credential string
		resultCode int
	}{
		{
			name:       "StatusCode=StatusUnauthorized,Malicious JWT value",
			credential: "MaliciousJWT",
			resultCode: http.StatusUnauthorized,
		},
		{
			name:       "StatusCode=StatusOK,Valid authentication",
			credential: "Bearer " + generateIDToken(user),
			resultCode: http.StatusOK,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := serveHTTP("GET", path, "", "", tc.credential)
			assert.Equal(t, tc.resultCode, resp.Code)
		})
	}

}

func TestCreateABookmarkOfAUser(t *testing.T) {
	const (
		fakeID = 100
	)

	user := getUser(Globs.Defaults.Account)
	bookmarkJSON, _ := json.Marshal(models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"})

	for _, tc := range []struct {
		name        string
		credential  string
		path        string
		payload     string
		cleanupStmt string
		resultCode  int
	}{
		{
			name:       "StatusCode=StatusBadRequest,Malformed request body",
			credential: "Bearer " + generateIDToken(user),
			path:       fmt.Sprintf("/v1/users/%v/bookmarks", user.ID),
			payload:    `{"InvalidJSONPayload"`,
			resultCode: http.StatusBadRequest,
		},
		{
			name:       "StatusCode=StatusUnauthorized,Invalid user id",
			credential: "InvalidJWT",
			path:       fmt.Sprintf("/v1/users/%v/bookmarks", fakeID),
			payload:    string(bookmarkJSON),
			resultCode: http.StatusUnauthorized,
		},
		{
			name:        "StatusCode=StatusCreated,Boookmark created",
			credential:  "Bearer " + generateIDToken(user),
			path:        fmt.Sprintf("/v1/users/%v/bookmarks", user.ID),
			payload:     string(bookmarkJSON),
			cleanupStmt: "SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE bookmarks; TRUNCATE TABLE users_bookmarks; SET FOREIGN_KEY_CHECKS=1",
			resultCode:  http.StatusCreated,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cleanupStmt != "" {
				defer func() {
					Globs.GormDB.Exec(tc.cleanupStmt)
				}()
			}

			resp := serveHTTP("POST", tc.path, tc.payload, "application/json", tc.credential)
			assert.Equal(t, tc.resultCode, resp.Code)
		})
	}
}

func TestGetBookmarksOfAUser(t *testing.T) {
	const (
		fakeID        = 100
		defaultLimit  = 10
		defaultOffset = 0
	)

	user := getUser(Globs.Defaults.Account)

	bookmarks := []models.Bookmark{
		models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
		models.Bookmark{Slug: "mock-slug-2", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
	}

	for _, tc := range []struct {
		name        string
		path        string
		credential  string
		bookmarks   []models.Bookmark
		cleanupStmt string
		resultCode  int
		resultMeta  *meta
	}{
		{
			name:       "StatusCode=StatusUnauthorized,List with invalid jwt",
			path:       fmt.Sprintf("/v1/users/%v/bookmarks", fakeID),
			credential: "INVALIDJWT",
			resultCode: http.StatusUnauthorized,
		},
		{
			name:        "StatusCode=StatusOK,List bookmarks of a user",
			path:        fmt.Sprintf("/v1/users/%v/bookmarks", user.ID),
			credential:  "Bearer " + generateIDToken(user),
			bookmarks:   bookmarks,
			cleanupStmt: "SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE bookmarks; TRUNCATE TABLE users_bookmarks; SET FOREIGN_KEY_CHECKS=1",
			resultCode:  http.StatusOK,
			resultMeta: &meta{
				Limit:  defaultLimit,
				Offset: defaultOffset,
				Total:  len(bookmarks),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.cleanupStmt != "" {
				defer func() {
					Globs.GormDB.Exec(tc.cleanupStmt)
				}()
			}

			if len(tc.bookmarks) > 0 {
				for _, v := range tc.bookmarks {
					s, _ := json.Marshal(v)
					serveHTTP("POST", tc.path, string(s), "application/json", tc.credential)
				}
			}

			resp := serveHTTP("GET", tc.path, "", "", tc.credential)
			assert.Equal(t, tc.resultCode, resp.Code)

			if tc.resultMeta != nil {
				body, _ := ioutil.ReadAll(resp.Result().Body)

				res := response{}
				json.Unmarshal(body, &res)

				assert.Equal(t, tc.resultMeta.Limit, res.Meta.Limit)
				assert.Equal(t, tc.resultMeta.Offset, res.Meta.Offset)
				assert.Equal(t, tc.resultMeta.Total, res.Meta.Total)
			}
		})

	}
}

func TestGetABookmarkOfAUser(t *testing.T) {
	type (
		// abstraction of the bookmarks state in database
		userBookmarks struct {
			UserEmail string
			Bookmarks []models.Bookmark
		}
	)

	for _, tc := range []struct {
		name                string
		testUserEmails      []string
		userBookmarks       []userBookmarks
		testTargetUserEmail string
		testTargetSlug      string
		testTargetHost      string
		cleanupStmt         string
		resultCode          int
		resultBody          *response
	}{
		{
			name: "StatusCode=StatusNotFound,A user does not have any bookmark",
			testUserEmails: []string{
				"testUserA@twreporter.org",
				"testUserB@twreporter.org",
			},
			userBookmarks: []userBookmarks{
				userBookmarks{UserEmail: "testUserA@twreporter.org"},
				userBookmarks{
					UserEmail: "testUserB@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
						models.Bookmark{Slug: "mock-slug-2", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
			},
			testTargetUserEmail: "testUserA@twreporter.org",
			testTargetSlug:      "mock-slug-1",
			cleanupStmt:         "SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE bookmarks; TRUNCATE TABLE users_bookmarks; SET FOREIGN_KEY_CHECKS=1",
			resultCode:          http.StatusNotFound,
		},
		{
			name: "StatusCode=StatusNotFound,A user does not own the bookmark",
			testUserEmails: []string{
				"testUserA@twreporter.org",
				"testUserB@twreporter.org",
			},
			userBookmarks: []userBookmarks{
				userBookmarks{
					UserEmail: "testUserA@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-2", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
				userBookmarks{
					UserEmail: "testUserB@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
						models.Bookmark{Slug: "mock-slug-2", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
			},
			testTargetUserEmail: "testUserA@twreporter.org",
			testTargetSlug:      "mock-slug-1",
			testTargetHost:      "mockhost",
			cleanupStmt:         "SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE bookmarks; TRUNCATE TABLE users_bookmarks; SET FOREIGN_KEY_CHECKS=1",
			resultCode:          http.StatusNotFound,
		},
		{
			name: "StatusCode=StatusNotFound,Query does not provide a valid host parameter",
			testUserEmails: []string{
				"testUserA@twreporter.org",
				"testUserB@twreporter.org",
			},
			userBookmarks: []userBookmarks{
				userBookmarks{
					UserEmail: "testUserA@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
				userBookmarks{
					UserEmail: "testUserB@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
						models.Bookmark{Slug: "mock-slug-2", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
			},
			testTargetUserEmail: "testUserA@twreporter.org",
			testTargetSlug:      "mock-slug-1",
			cleanupStmt:         "SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE bookmarks; TRUNCATE TABLE users_bookmarks; SET FOREIGN_KEY_CHECKS=1",
			resultCode:          http.StatusNotFound,
		},
		{
			name: "StatusCode=StatusOK,Successful get a bookmark of a user",
			testUserEmails: []string{
				"testUserA@twreporter.org",
				"testUserB@twreporter.org",
			},
			userBookmarks: []userBookmarks{
				userBookmarks{
					UserEmail: "testUserA@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
				userBookmarks{
					UserEmail: "testUserB@twreporter.org",
					Bookmarks: []models.Bookmark{
						models.Bookmark{Slug: "mock-slug-1", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
						models.Bookmark{Slug: "mock-slug-2", Host: "mockhost", Title: "mocktitle", Thumbnail: "mockthumb"},
					},
				},
			},
			testTargetUserEmail: "testUserA@twreporter.org",
			testTargetSlug:      "mock-slug-1",
			testTargetHost:      "mockhost",
			cleanupStmt:         "SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE bookmarks; TRUNCATE TABLE users_bookmarks; SET FOREIGN_KEY_CHECKS=1",
			resultCode:          http.StatusOK,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Create test users
			for _, v := range tc.testUserEmails {
				user := createUser(v)
				defer func() { deleteUser(user) }()
			}

			// Defer cleanup statment for cleaning up bookmark tables
			if tc.cleanupStmt != "" {
				defer func() { Globs.GormDB.Exec(tc.cleanupStmt) }()
			}

			// Insert bookmarks to corresponding users according the state
			for _, u := range tc.userBookmarks {
				user := getUser(u.UserEmail)
				credential := "Bearer " + generateIDToken(user)
				for _, b := range u.Bookmarks {
					s, _ := json.Marshal(b)
					serveHTTP("POST", fmt.Sprintf("/v1/users/%v/bookmarks", user.ID), string(s), "application/json", credential)
				}
			}

			targetUser := getUser(tc.testTargetUserEmail)
			resp := serveHTTP("GET", fmt.Sprintf("/v1/users/%v/bookmarks/%s?host=%s", targetUser.ID, tc.testTargetSlug, tc.testTargetHost), "", "", "Bearer "+generateIDToken(targetUser))
			assert.Equal(t, tc.resultCode, resp.Code)
			if tc.resultBody != nil {
				body, _ := ioutil.ReadAll(resp.Result().Body)
				res := response{}
				json.Unmarshal(body, &res)
				assert.Equal(t, tc.testTargetSlug, res.Bookmark.Slug)
			}
		})
	}
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
			jwt := generateIDToken(tc.deleteInfo.User)

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
