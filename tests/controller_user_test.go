package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/storage"
)

type (
  userData struct {
	ID                  uint `json:"id"`
	AgreeDataCollection bool `json:"agree_data_collection"`
	ReadPostsCount      int  `json:"read_posts_count"`
	ReadPostsSec        int  `json:"read_posts_sec"`
  }

  responseBodyUser struct {
	Status string   `json:"status"`
	Data   userData `json:"data"`
  }
)

func TestGetUser_Success(t *testing.T) {
	var resBody responseBodyUser
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)
	as := storage.NewGormStorage(Globs.GormDB)
	if err := as.UpdateUser(models.User{
		ID:                  user.ID,
		AgreeDataCollection: true,
		ReadPostsCount:      19,
		ReadPostsSec:        3360,
	}); nil != err {
		fmt.Println(err.Error())
	}
	updatedUser := getUser(Globs.Defaults.Account)

	// Send request to test GetUser function
	response := serveHTTP(http.MethodGet, fmt.Sprintf("/v2/users/%d", user.ID), "", "", fmt.Sprintf("Bearer %v", jwt))
	fmt.Print(response.Body)
	resBodyInBytes, _ := ioutil.ReadAll(response.Result().Body)
	json.Unmarshal(resBodyInBytes, &resBody)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, updatedUser.AgreeDataCollection, resBody.Data.AgreeDataCollection)
	assert.Equal(t, updatedUser.ReadPostsCount, resBody.Data.ReadPostsCount)
	assert.Equal(t, updatedUser.ReadPostsSec, resBody.Data.ReadPostsSec)
}

func TestSetUser_Success(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	var InterestIDsKeys []string

	for k := range globals.Conf.Mailchimp.InterestIDs {
		InterestIDsKeys = append(InterestIDsKeys, k)
	}

	// Mocking preferences
	preferences := models.UserPreference{
		ReadPreference: []string{"international", "cross_straits"},
		Maillist:       InterestIDsKeys,
	}
	payload, _ := json.Marshal(preferences)

	// Send request to test SetUser function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestSetUser_InvalidMaillist(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking preferences
	preferences := models.UserPreference{
		ReadPreference: []string{"international", "cross_straits"},
		Maillist:       []string{"maillist1", "maillist2", "maillist5"},
	}
	payload, _ := json.Marshal(preferences)

	// Send request to test SetUser function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusBadRequest, response.Code)
}
