package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"io/ioutil"

	"gopkg.in/guregu/null.v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/storage"
	"github.com/twreporter/go-api/internal/news"
)

type (
	reqBody struct {
		PostID         null.String `json:"post_id"`
		ReadPostsCount null.Bool   `json:"read_posts_count"`
		ReadPostsSec   null.Int    `json:"read_posts_sec"`
	}
	readPostsData struct {
		UserID         string    `json:"user_id"`
		PostID         string    `json:"post_id"`
		ReadPostsCount null.Bool `json:"read_posts_count"`
		ReadPostsSec   null.Int  `json:"read_posts_sec"`
	}
	resBody struct {
		Status string        `json:"status"`
		Data   readPostsData `json:"data"`
	}
	reqBodyFootprint struct {
		PostID         null.String `json:"post_id"`
	}
	respBodyFootprint struct {
		Status    string                 `json:"status"`
		Records   []news.MetaOfFootprint `json:"records"`
	}
)

const (
	mockPostID = "573422ab3fac3c7322005ae9"
	mockPostSec = 3660
)

func TestSetUserAnalytics_Success(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	analytics := reqBody{
		PostID: null.NewString(mockPostID, true),
		ReadPostsCount: null.NewBool(true, true),
		ReadPostsSec: null.NewInt(mockPostSec, true),
	}
	payload, _ := json.Marshal(analytics)

	// Send request to test SetUserAnalytics function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestSetUserAnalytics_ReadingTimeExceedMaximum(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	analytics := reqBody{
		PostID: null.NewString(mockPostID, true),
		ReadPostsSec: null.NewInt(7201, true),
	}
	payload, _ := json.Marshal(analytics)

	var resBody resBody
	// Send request to test SetUserAnalytics function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))
	resBodyInBytes, _ := ioutil.ReadAll(response.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)

	expectedReadPostSec := null.NewInt(7200, true)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, expectedReadPostSec, resBody.Data.ReadPostsSec)
}

func TestSetUserAnalytics_EmptyPostID(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	analytics := reqBody{
		ReadPostsCount: null.NewBool(true, true),
		ReadPostsSec: null.NewInt(mockPostSec, true),
	}
	payload, _ := json.Marshal(analytics)

	// Send request to test SetUserAnalytics function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestSetUserAnalytics_InvalidUserID(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Send request to test SetUserAnalytics function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics", user.ID + 1), "{}", "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestSetUserReadingFootprint_Success(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	footprint := reqBodyFootprint{
		PostID: null.NewString(mockPostID, true),
	}
	payload, _ := json.Marshal(footprint)

	// Send request to test SetUserReadingFootprint function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics/reading-footprint", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestSetUserReadingFootprint_EmptyPostID(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	footprint := reqBodyFootprint{}
	payload, _ := json.Marshal(footprint)

	// Send request to test SetUserReadingFootprint function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics/reading-footprint", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestSetUserReadingFootprint_InvalidUserID(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Send request to test SetUserReadingFootprint function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics/reading-footprint", user.ID + 1), "{}", "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusForbidden, response.Code)
}

func TestGetReadingFootprint_Success(t *testing.T) {
	var resBody respBodyFootprint

	// Mocking Post
	mockPostObjectID, err := primitive.ObjectIDFromHex(mockPostID)
	if err != nil {
		fmt.Println(err.Error())
	}
	post := news.MetaOfFootprint{ID: mockPostObjectID}
	err = createPost(post)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)
	as := storage.NewAnalyticsGormStorage(Globs.GormDB)
	if _, err := as.UpdateUserReadingFootprint(fmt.Sprint(user.ID), mockPostID); nil != err {
		fmt.Println(err.Error())
	}

	// Send request to test GetUser function
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/users/%d/analytics/reading-footprint", user.ID), "", "", fmt.Sprintf("Bearer %v", jwt))
	resBodyInBytes, _ := ioutil.ReadAll(response.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, 1, len(resBody.Records))
	assert.Equal(t, mockPostObjectID, resBody.Records[0].ID)
	assert.NotEmpty(t, resBody.Records[0].UpdatedAt)
}

func TestGetUserReadingFootprint_InvalidUserID(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Send request to test SetUserReadingFootprint function
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/users/%d/analytics/reading-footprint", user.ID + 1), "{}", "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusForbidden, response.Code)
}
