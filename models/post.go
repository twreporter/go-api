package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Brief ...
type Brief struct {
	HTML    string   `json:"html"`
	APIData []bson.M `bson:"apiData" json:"api_data"`
}

// PostMeta ...
type PostMeta struct {
	ID               bson.ObjectId   `bson:"_id" json:"id"`
	Slug             string          `json:"slug"`
	Name             string          `json:"name"`
	Subtitle         string          `json:"subtitle"`
	State            string          `json:"state"`
	HeroImage        *Image          `bson:"-" json:"hero_image,omitempty"`
	HeroImageOrigin  bson.ObjectId   `bson:"heroImage" json:"-"`
	Brief            *Brief          `json:"brief,omitempty"`
	Categories       []Category      `bson:"-" json:"categories,omitempty"`
	CategoriesOrigin []bson.ObjectId `bson:"categories,omitempty" json:"-"`
	Style            string          `json:"style"`
	Copyright        string          `json:"copyright"`
	Tags             []Tag           `bson:"-" json:"tags,omitempty"`
	TagsOrigin       []bson.ObjectId `bson:"tags,omitempty" json:"-"`
	OgDescription    string          `json:"og_description"`
	OgImage          *Image          `bson:"-" json:"og_image,omitempty"`
	OgImageOrigin    bson.ObjectId   `bson:"og_image" json:"-"`
	IsFeatured       bool            `bson:"isFeatured" json:"is_featured"`
	PublishedDate    time.Time       `bson:"publishedDate" json:"published_date"`
	Topic            *Topic          `bson:"-" json:"topic,omitempty"`
	TopicOrigin      bson.ObjectId   `bson:"topics" json:"-"`
}

/*
type Post struct {
	Writters          []Contact
	Photographers     []Contact
	Designers         []Contact
	Engineers         []Contact
	ExtendByline      string
	LeadingVideo      Video
	LeadingVideoID    uint
	HeroImage         Image
	HeroImageID       uint
	BriefDraftState   string     `gorm:"text"`
	BriefHTML         string     `gorm:"text"`
	BriefAPIData      string     `gorm:"text"`
	ContentDraftState string     `gorm:"text"`
	ContentHTML       string     `gorm:"text"`
	ContentAPIData    string     `gorm:"text"`
	Categories        []Category `gorm:"many2many:posts_categories"`
	Style             string     `gorm:"default:article;size:30;not null;index"`
	Bookmark          string     `gorm:"default:untitled;size:30"`
	BookmarkOrder     uint       `gorm:"default:1;"`
	RelatedBookmarks  []Post     `gorm:"foreignkey:post_id;associationforeignkey:relatedbookmarkpost_id;many2many:posts_relatedbookmarkposts;"`
	Topic             Topic
	TopicID           uint
	Copyright         string `gorm:"default:Copyrighted"`
	Tags              []Tag  `gorm:"many2many:posts_tags"`
	Relateds          []Post `gorm:"foreignkey:post_id;associationforeignkey:relatedpost_id;many2many:posts_relatedposts;"`
	OgTitle           string
	OgDescription     string
	OgImage           Image
	OgImageID         uint
	IsFeatured        bool `gorm:"default:0;index"`
}
*/
