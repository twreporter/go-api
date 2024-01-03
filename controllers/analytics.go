package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
)

type (
	reqBody struct {
		PostID         null.String `json:"post_id"`
		ReadPostsCount null.Bool   `json:"read_posts_count"`
		ReadPostsSec   null.Int    `json:"read_posts_sec"`
	}
	respBody struct {
		UserID         string    `json:"user_id"`
		PostID         string    `json:"post_id"`
		ReadPostsCount null.Bool `json:"read_posts_count"`
		ReadPostsSec   null.Int  `json:"read_posts_sec"`
	}
)

func (mc *MembershipController) SetUserAnalytics(c *gin.Context) (int, gin.H, error) {
	var req reqBody
	var resp respBody
	var isExisted bool
	var err error
	userID := c.Param("userID")
	if err = c.BindJSON(&req); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()}, nil
	}

	if req.PostID.Valid == false {
		return http.StatusBadRequest, gin.H{"status": "fail", "message": "post_id is required"}, nil
	}
	if req.ReadPostsSec.Valid && req.ReadPostsSec.Int64 < 0 {
		return http.StatusBadRequest, gin.H{"status": "fail", "message": "read_posts_sec cannot be negative"}, nil
	}
	resp.UserID = userID
	resp.PostID = req.PostID.String

	if null.Bool.ValueOrZero(req.ReadPostsCount) == true {
		isExisted, err = mc.Storage.UpdateUserReadingPostCount(userID, req.PostID.String)
		if err != nil {
			return toResponse(err)
		}
		resp.ReadPostsCount = null.NewBool(true, true)
	}

	if null.Int.IsZero(req.ReadPostsSec) == false {
		// update read post time
		err = mc.Storage.UpdateUserReadingPostTime(userID, req.PostID.String, int(req.ReadPostsSec.Int64))
		if err != nil {
			return toResponse(err)
		}
		isExisted = false
		resp.ReadPostsSec = req.ReadPostsSec
	}

	if isExisted {
		return http.StatusOK, gin.H{"status": "success", "data": resp}, nil
	}
	return http.StatusCreated, gin.H{"status": "success", "data": resp}, nil
}
