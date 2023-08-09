package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/twreporter/go-api/configs/constants"
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
			"key":     role.Key,
		}
	}

	var activated *time.Time
	if user.Activated.Valid && !user.Activated.Time.IsZero() {
		activated = &user.Activated.Time
	}

	mailGroups := make([]string, 0)
	for _, group := range user.MailGroups {
		for key, value := range globals.Conf.Mailchimp.InterestIDs {
			if value == group.MailgroupID {
				mailGroups = append(mailGroups, key)
				break
			}
		}
	}
	readPreferenceArr := make([]string, 0)
	if user.ReadPreference.Valid {
		readPreferenceArr = strings.Split(user.ReadPreference.String, ",")
	}

	return http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"user_id":           userID,
		"first_name":        user.FirstName.String,
		"last_name":         user.LastName.String,
		"email":             user.Email.String,
		"registration_date": user.RegistrationDate.Time,
		"activated":         activated,
		"roles":             roles,
		"read_preference":   readPreferenceArr,
		"maillist":          mailGroups,
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
	maillists := make([]string, 0)
	for _, maillist := range preferences.Maillist {
		convertedMaillist, exists := globals.Conf.Mailchimp.InterestIDs[maillist]
		if !exists {
			return http.StatusBadRequest, gin.H{"status": "error", "message": "invalid maillist value"}, errors.New("Invalid maillist value")
		}
		maillists = append(maillists, convertedMaillist)
	}

	// send explorer role email for first time user
	user, err := mc.Storage.GetUserByID(fmt.Sprint(userID))
	if err != nil {
		return http.StatusInternalServerError, gin.H{"status": "error", "message": "user does not exist"}, err
	}

	if !user.Activated.Valid || user.Activated.Time.IsZero() {
		roleCheck, roleCheckErr := mc.Storage.HasRole(user, constants.RoleExplorer)
		if roleCheckErr != nil {
			log.Println("Error checking role:", roleCheckErr)
		}

		log.WithFields(log.Fields{
			"user.Activated.Valid":         user.Activated.Valid,
			"user.Activated.Time.IsZero()": user.Activated.Time.IsZero(),
			"roleCheck":                    roleCheck,
			"sendAssignRoleMail":           !user.Activated.Valid && user.Activated.Time.IsZero() && roleCheck,
		}).Info("SetUser Activated Role check")

		if roleCheck {
			go mc.sendAssignRoleMail(constants.RoleExplorer, user.Email.String)
		}
	}

	// Call UpdateReadPreferenceOfUser to save the preferences.ReadPreference to DB
	if err = mc.Storage.UpdateReadPreferenceOfUser(userID, preferences.ReadPreference); err != nil {
		return toResponse(err)
	}

	// Call CreateMaillistOfUser to save the preferences.Maillist to DB
	if err = mc.Storage.CreateMaillistOfUser(userID, maillists); err != nil {
		return toResponse(err)
	}

	return http.StatusCreated, gin.H{"status": "ok", "record": preferences}, nil
}
