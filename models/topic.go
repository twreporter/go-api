package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TopicMeta struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	Slug          string        `json:"slug"`
	Name          string        `json:"name"`
	TopicName     string        `json:"topic_name"`
	Title         string        `json:"title"`
	Subtitle      string        `json:"subtitle"`
	Headline      string        `json:"headline"`
	PublishedDate time.Time     `bson:"publishedDate" json:"published_date"`
	Description   string        `json:"description"`
	OgTitle       string        `json:"og_title"`
	OgDescription string        `json:"og_description"`
}

type Topic struct {
	ID                         bson.ObjectId   `json:"id" bson:"_id"`
	Slug                       string          `json:"slug"`
	TopicName                  string          `json:"topic_name"`
	Title                      string          `json:"title"`
	TitlePosition              string          `json:"title_position"`
	Subtitle                   string          `json:"subtitle"`
	Headline                   string          `json:"headline"`
	State                      string          `json:"state"`
	PublishedDate              time.Time       `bson:"publishedDate" json:"published_date"`
	Description                Brief           `json:"description,omitempty"`
	TeamDescription            Brief           `json:"team_description,omitempty"`
	Relateds                   []PostMeta      `bson:"-" json:",omitempty"`
	RelatedsOrigin             []bson.ObjectId `bson:"relateds,omitempty" json:"-"`
	RelatedsFormat             string          `json:"relateds_format"`
	RelatedsBackground         string          `json:"relateds_background"`
	LeadingImage               *Image          `bson:"-" json:"leading_image,omitempty"`
	LeadingImageOrigin         bson.ObjectId   `bson:"leading_image,omitempty" json:"-"`
	LeadingImagePortrait       *Image          `bson:"-" json:"leading_image_portrait,omitempty"`
	LeadingImagePortraitOrigin bson.ObjectId   `bson:"leading_image_portrait,omitempty" json:"-"`
	LeadingVideo               Video           `bson:"-" json:"leading_video,omitempty"`
	LeadingVideoOrigin         bson.ObjectId   `bson:"leading_video,omitempty" json:"-"`
	OgTitle                    string          `json:"og_title"`
	OgDescription              string          `json:"og_description"`
	OgImage                    *Image          `bson:"-" json:"og_image,omitempty"`
	OgImageOrigin              bson.ObjectId   `bson:"og_image,omitempty" json:"-"`
}
