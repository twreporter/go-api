package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ContentBody ...
type ContentBody struct {
	HTML    string   `json:"html"`
	APIData []bson.M `bson:"apiData" json:"api_data"`
}

// Post ...
type Post struct {
	ID                  bson.ObjectId   `bson:"_id" json:"id"`
	Slug                string          `bson:"slug" json:"slug"`
	Name                string          `bson:"name" json:"name"`
	Title               string          `bson:"title" json:"title"`
	Subtitle            string          `bson:"subtitle" json:"subtitle"`
	State               string          `bson:"state" json:"state"`
	HeroImage           *Image          `bson:"-" json:"hero_image,omitempty"`
	HeroImageOrigin     bson.ObjectId   `bson:"heroImage" json:"-"`
	Brief               *ContentBody    `bson:"brief,omitempty" json:"brief,omitempty"`
	Categories          []Category      `bson:"-" json:"categories,omitempty"`
	CategoriesOrigin    []bson.ObjectId `bson:"categories,omitempty" json:"-"`
	Style               string          `bson:"style" json:"style"`
	Copyright           string          `bson:"copyright" json:"copyright"`
	Tags                []Tag           `bson:"-" json:"tags,omitempty"`
	TagsOrigin          []bson.ObjectId `bson:"tags,omitempty" json:"-"`
	OgTitle             string          `bson:"og_title" json:"og_title"`
	OgDescription       string          `bson:"og_description" json:"og_description"`
	OgImage             *Image          `bson:"-" json:"og_image,omitempty"`
	OgImageOrigin       bson.ObjectId   `bson:"og_image" json:"-"`
	IsFeatured          bool            `bson:"isFeatured" json:"is_featured"`
	PublishedDate       time.Time       `bson:"publishedDate" json:"published_date"`
	Topic               *Topic          `bson:"-" json:"topic,omitempty"`
	TopicOrigin         bson.ObjectId   `bson:"topics,omitempty" json:"-"`
	Writters            []Author        `bson:"-" json:"writters,omitempty"`
	WrittersOrigin      []bson.ObjectId `bson:"writters,omitempty" json:"-"`
	Photographers       []Author        `bson:"-" json:"photographers,omitempty"`
	PhotographersOrigin []bson.ObjectId `bson:"photographers,omitempty" json:"-"`
	Designers           []Author        `bson:"-" json:"designers,omitempty"`
	DesignersOrigin     []bson.ObjectId `bson:"designers,omitempty" json:"-"`
	Engineers           []Author        `bson:"-" json:"engineers,omitempty"`
	EngineersOrigin     []bson.ObjectId `bson:"engineers,omitempty" json:"-"`
	ExtendByline        string          `bson:"extend_byline" json:"extend_byline"`
	LeadingVideo        *Video          `bson:"-" json:"leading_video,omitempty"`
	LeadingVideoOrigin  bson.ObjectId   `bson:"leading_video,omitempty" json:"-"`
	Content             *ContentBody    `bson:"content,omitempty" json:"content,omitempty"`
	Relateds            []Post          `bson:"-" json:"relateds,omitempty"`
	RelatedsOrigin      []bson.ObjectId `bson:"relateds,omitempty" json:"-"`
}
