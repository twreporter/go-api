package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/models"
)

// GetUser given userID, this func will try to get the user record
func (mc *MembershipController) GetUser(c *gin.Context) (int, gin.H, error) {
	userID := c.Param("userID")

	user, err := mc.Storage.GetUserByID(userID)
	if err != nil {
		return toResponse(err)
	}

	roles := make([]gin.H, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = gin.H{
			"id":      role.ID, // does frontend need ID?
			"name":    role.Name,
			"name_en": role.NameEn,
			"key":     role.Key,
		}
	}

	var activated *time.Time
	if user.Activated.Valid && !user.Activated.Time.IsZero() {
		activated = &user.Activated.Time
	}

	readPreferenceArr := make([]string, 0)
	if user.ReadPreference.Valid {
		readPreferenceArr = strings.Split(user.ReadPreference.String, ",")
	}

	return http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"user_id":                userID,
		"first_name":             user.FirstName.String,
		"last_name":              user.LastName.String,
		"email":                  user.Email.String,
		"registration_date":      user.RegistrationDate.Time,
		"activated":              activated,
		"roles":                  roles,
		"read_preference":        readPreferenceArr,
		"agree_data_collection":  user.AgreeDataCollection,
		"read_posts_count":       user.ReadPostsCount,
		"read_posts_sec":         user.ReadPostsSec,
		"is_showofflinedonation": user.IsShowOfflineDonation,
	},
	}, nil
}

// SetUser given userID and POST body, this func will try to create record in the related table,
// and build the relationship between records and user
func (mc *MembershipController) SetUser(c *gin.Context) (int, gin.H, error) {
	userID := c.Param("userID")
	var preferences models.UserPreference
	err := c.BindJSON(&preferences)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	// Call UpdateReadPreferenceOfUser to save the preferences.ReadPreference to DB
	if err = mc.Storage.UpdateReadPreferenceOfUser(userID, preferences.ReadPreference); err != nil {
		return toResponse(err)
	}

	// Call UpdateUser to save preferences.IsShowOfflineDonation to DB
	if preferences.IsShowOfflineDonation.Valid {
		matchedUser, err := mc.Storage.GetUserByID(userID)
		if err != nil {
			return toResponse(err)
		}
		matchedUser.IsShowOfflineDonation = preferences.IsShowOfflineDonation.Bool
		if err = mc.Storage.UpdateUser(matchedUser); err != nil {
			return toResponse(err)
		}
	}

	return http.StatusCreated, gin.H{"status": "ok", "record": preferences}, nil
}
