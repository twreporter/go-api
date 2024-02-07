package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"context"
	"time"

	"gopkg.in/guregu/null.v3"
	"github.com/gin-gonic/gin"
	"github.com/twreporter/go-api/storage"
	"github.com/twreporter/go-api/models"
	"github.com/twreporter/go-api/internal/news"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewAnalyticsController(gs storage.AnalyticsGormStorage, ms storage.AnalyticsMongoStorage) *AnalyticsController {
	return &AnalyticsController{gs, ms}
}

type AnalyticsController struct {
	gs storage.AnalyticsGormStorage
	ms storage.AnalyticsMongoStorage
}

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
	reqBodyFootprint struct {
		PostID         null.String `json:"post_id"`
	}
)

func (ac *AnalyticsController) SetUserAnalytics(c *gin.Context) (int, gin.H, error) {
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
		isExisted, err = ac.gs.UpdateUserReadingPostCount(userID, req.PostID.String)
		if err != nil {
			return toResponse(err)
		}
		resp.ReadPostsCount = null.NewBool(true, true)
	}

	if null.Int.IsZero(req.ReadPostsSec) == false {
		// update read post time
		err = ac.gs.UpdateUserReadingPostTime(userID, req.PostID.String, int(req.ReadPostsSec.Int64))
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

func (ac *AnalyticsController) GetUserAnalyticsReadingFootprint(c *gin.Context) (int, gin.H, error) {
	// parameter validation
	userID := c.Param("userID")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	if limit == 0 {
		limit = 10
	}

	// get footprint posts of target user
	footprints, total, err := ac.gs.GetFootprintsOfAUser(userID, limit, offset)
	if err != nil {
		return toResponse(err)
	}

	// fetch posts meta from mongo db
	type footprintMeta struct {
		bookmarkID string
		updatedAt  time.Time
	}
	postIds := make([]string, len(footprints))
	bookmarkMap := make(map[primitive.ObjectID]footprintMeta)
	for index, footprint := range footprints {
		postIds[index] = footprint.PostID
		objectID, err := primitive.ObjectIDFromHex(footprint.PostID)
		if err != nil {
			continue;
		}
		bookmarkMap[objectID] = footprintMeta{bookmarkID:footprint.BookmarkID, updatedAt:footprint.UpdatedAt}
	}
	var posts []news.MetaOfFootprint
	posts, err = ac.ms.GetPostsOfIDs(context.Background(), postIds)
	if err != nil {
		return toResponse(err)
	}

	// add bookmarks in posts
	for index, post := range posts {
		posts[index].BookmarkID = bookmarkMap[post.ID].bookmarkID
		posts[index].UpdatedAt = bookmarkMap[post.ID].updatedAt
	}

	return http.StatusOK, gin.H{"status": "ok", "records": posts, "meta": models.MetaOfResponse{
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}}, nil
}

func (ac *AnalyticsController) SetUserAnalyticsReadingFootprint(c *gin.Context) (int, gin.H, error) {
	var req reqBodyFootprint
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

	isExisted, err = ac.gs.UpdateUserReadingFootprint(userID, req.PostID.String)
	if err != nil {
		return toResponse(err)
	}

	if isExisted {
		return http.StatusOK, gin.H{"status": "success"}, nil
	}
	return http.StatusCreated, gin.H{"status": "success"}, nil
}
