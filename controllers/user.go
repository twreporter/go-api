package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/models"
)

// GetUser given userID, this func will try to get the user record, joined with users_mailgroup table, from DB
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
		}
	}

	return http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"email":             user.Email.String,
		"registration_date": user.RegistrationDate.Time,
		"roles":             roles,
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

	// Convert maillist values using the mapping array
	for i, maillist := range preferences.Maillist {
		convertedMaillist, exists := globals.Conf.Mailchimp.InterestIDs[maillist]
		if !exists {
			return http.StatusBadRequest, gin.H{"status": "error", "message": "invalid maillist value"}, errors.New("Invalid maillist value")
		}
		preferences.Maillist[i] = convertedMaillist
	}

	// Call UpdateReadPreferenceOfUser to save the preferences.ReadPreference to DB
	if err = mc.Storage.UpdateReadPreferenceOfUser(userID, preferences.ReadPreference); err != nil {
		return toResponse(err)
	}

	// Call CreateMaillistOfUser to save the preferences.Maillist to DB
	if err = mc.Storage.CreateMaillistOfUser(userID, preferences.Maillist); err != nil {
		return toResponse(err)
	}

	return http.StatusCreated, gin.H{"status": "ok", "record": preferences}, nil
}
