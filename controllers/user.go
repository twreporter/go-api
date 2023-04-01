package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/models"
)

// GetUser given userID, this func will try to get the user record, joined with users_mailgroup table, from DB
func (mc *MembershipController) GetUser(c *gin.Context) (int, gin.H, error) {
	userID := c.Param("userID")

	user, err := mc.Storage.GetUserByID(userID)
	if err != nil {
		return toResponse(err)
	}

	return http.StatusOK, gin.H{"status": "ok", "record": user}, nil
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

	// Call CreateMaillistOfUser to save the preferences.Maillist to DB
	if err = mc.Storage.CreateMaillistOfUser(userID, preferences.Maillist); err != nil {
		return toResponse(err)
	}

	return http.StatusCreated, gin.H{"status": "ok", "record": preferences}, nil
}
