package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Topic ...
type Topic struct {
	ID                         bson.ObjectId   `json:"id" bson:"_id"`
	Slug                       string          `bson:"slug" json:"slug"`
	Name                       string          `bson:"name" json:"name"`
	TopicName                  string          `bson:"topic_name" json:"topic_name"`
	Title                      string          `bson:"title" json:"title"`
	TitlePosition              string          `bson:"title_position" json:"title_position"`
	Subtitle                   string          `bson:"subtitle" json:"subtitle"`
	Headline                   string          `bson:"headline" json:"headline"`
	State                      string          `bson:"state" json:"state"`
	Description                *ContentBody    `bson:"description,omitempty" json:"description,omitempty"`
	TeamDescription            *ContentBody    `bson:"team_description,omitempty" json:"team_description,omitempty"`
	Relateds                   []Post          `bson:"-" json:"relateds,omitempty"`
	RelatedsOrigin             []bson.ObjectId `bson:"relateds,omitempty" json:"-"`
	RelatedsFormat             string          `bson:"relateds_format" json:"relateds_format"`
	RelatedsBackground         string          `bson:"relateds_background" json:"relateds_background"`
	LeadingImage               *Image          `bson:"-" json:"leading_image,omitempty"`
	LeadingImageOrigin         bson.ObjectId   `bson:"leading_image,omitempty" json:"-"`
	LeadingImagePortrait       *Image          `bson:"-" json:"leading_image_portrait,omitempty"`
	LeadingImagePortraitOrigin bson.ObjectId   `bson:"leading_image_portrait,omitempty" json:"-"`
	LeadingVideo               *Video          `bson:"-" json:"leading_video,omitempty"`
	LeadingVideoOrigin         bson.ObjectId   `bson:"leading_video,omitempty" json:"-"`
	OgTitle                    string          `bson:"og_title" json:"og_title"`
	OgDescription              string          `bson:"og_description" json:"og_description"`
	OgImage                    *Image          `bson:"-" json:"og_image,omitempty"`
	OgImageOrigin              bson.ObjectId   `bson:"og_image,omitempty" json:"-"`
	PublishedDate              time.Time       `bson:"publishedDate" json:"published_date"`
	UpdatedAt                  time.Time       `bson:"updatedAt" json:"updated_at"`
	Full                       bool            `bson:"-" json:"full"`
}
