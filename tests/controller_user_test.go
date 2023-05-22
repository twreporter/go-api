package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twreporter/go-api/models"
)

func TestSetUser_Success(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	var InterestIDsKeys []string

	for k := range models.InterestIDs {
		InterestIDsKeys = append(InterestIDsKeys, k)
	}

	// Mocking preferences
	preferences := models.UserPreference{
		ReadPreference: []string{"international", "cross_straits"},
		Maillist:       InterestIDsKeys,
	}
	payload, _ := json.Marshal(preferences)

	// Send request to test SetUser function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/user/%d", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

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
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/user/%d", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusBadRequest, response.Code)
}
