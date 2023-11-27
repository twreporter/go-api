package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twreporter/go-api/models"
	"gopkg.in/guregu/null.v3"
)

type (
	reqBody struct {
		PostID         null.String `json:"post_id"`
		ReadPostsCount null.Bool   `json:"read_posts_count"`
		ReadPostsSec   null.Int    `json:"read_posts_sec"`
	}
)

func TestSetUserAnalytics_Success(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	analytics := reqBody{
		PostID: null.NewString("e85bcd2u", true),
		ReadPostsCount: null.NewBool(true, true),
		ReadPostsSec: null.NewInt(3660, true),
	}
	payload, _ := json.Marshal(analytics)

	// Send request to test SetUserAnalytics function
	response := serveHTTP(http.MethodPost, fmt.Sprintf("/v2/users/%d/analytics", user.ID), string(payload), "application/json", fmt.Sprintf("Bearer %v", jwt))

	assert.Equal(t, http.StatusCreated, response.Code)
}

func TestSetUserAnalytics_EmptyPostID(t *testing.T) {
	// Mocking user
	var user models.User = getUser(Globs.Defaults.Account)
	jwt := generateIDToken(user)

	// Mocking request body
	analytics := reqBody{
		ReadPostsCount: null.NewBool(true, true),
		ReadPostsSec: null.NewInt(3660, true),
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
