package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twreporter/go-api/models"
)

type UserPreferences struct {
	ReadPreference []string `json:"read_preference"`
	Maillist       []string `json:"maillist"`
}

func TestSetUser_Success(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking preferences
	preferences := UserPreferences{
		ReadPreference: []string{"test1", "test2"},
		Maillist:       []string{"maillist1", "maillist2", "maillist5"},
	}
	payload, _ := json.Marshal(preferences)

	// Send request to test SetUser function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/user/%d", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusCreated, response.Code)
}
