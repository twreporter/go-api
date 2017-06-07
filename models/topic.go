package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// TopicMeta ...
type TopicMeta struct {
	ID                 bson.ObjectId `bson:"_id" json:"id"`
	Slug               string        `json:"slug"`
	Name               string        `json:"name"`
	TopicName          string        `json:"topic_name"`
	Title              string        `json:"title"`
	Subtitle           string        `json:"subtitle"`
	Headline           string        `json:"headline"`
	PublishedDate      time.Time     `bson:"publishedDate" json:"published_date"`
	Description        string        `json:"description"`
	LeadingImage       *Image        `bson:"-" json:"leading_image,omitempty"`
	LeadingImageOrigin bson.ObjectId `bson:"leading_image,omitempty" json:"-"`
	OgTitle            string        `json:"og_title"`
	OgDescription      string        `json:"og_description"`
	OgImage            *Image        `bson:"-" json:"og_image,omitempty"`
	OgImageOrigin      bson.ObjectId `bson:"og_image,omitempty" json:"-"`
}

// Topic ...
type Topic struct {
	ID                         bson.ObjectId   `json:"id" bson:"_id"`
	Slug                       string          `json:"slug"`
	TopicName                  string          `bson:"topic_name"`
	Title                      string          `bson:"title"`
	TitlePosition              string          `bson:"title_position,omitempty"`
	Subtitle                   string          `bson:"subtitle"`
	Headline                   string          `bson:"headline"`
	State                      string          `bson:"state"`
	PublishedDate              time.Time       `bson:"publishedDate" json:"published_date"`
	Description                Brief           `bson:"description,omitempty"`
	TeamDescription            Brief           `bson:"team_description,omitempty"`
	Relateds                   []PostMeta      `bson:"-" json:",omitempty"`
	RelatedsOrigin             []bson.ObjectId `bson:"relateds,omitempty" json:"-"`
	RelatedsFormat             string          `bson:"relateds_format"`
	RelatedsBackground         string          `bson:"relateds_background"`
	LeadingImage               *Image          `bson:"-" json:"leading_image,omitempty"`
	LeadingImageOrigin         bson.ObjectId   `bson:"leading_image,omitempty" json:"-"`
	LeadingImagePortrait       *Image          `bson:"-" json:"leading_image_portrait,omitempty"`
	LeadingImagePortraitOrigin bson.ObjectId   `bson:"leading_image_portrait,omitempty" json:"-"`
	LeadingVideo               *Video          `bson:"-" json:"leading_video,omitempty"`
	LeadingVideoOrigin         bson.ObjectId   `bson:"leading_video,omitempty" json:"-"`
	OgTitle                    string          `bson:"og_title"`
	OgDescription              string          `bson:"og_description"`
	OgImage                    *Image          `bson:"-" json:"og_image,omitempty"`
	OgImageOrigin              bson.ObjectId   `bson:"og_image,omitempty" json:"-"`
}
