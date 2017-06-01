package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ImageAsset ...
type ImageAsset struct {
	Height uint   `json:"height"`
	Width  uint   `json:"width"`
	URL    string `json:"url"`
}
type ResizedTargets struct {
	Mobile  ImageAsset `json:"mobile"`
	Tiny    ImageAsset `json:"tiny"`
	Desktop ImageAsset `json:"desktop"`
	Tablet  ImageAsset `json:"tablet"`
}

// Image ...
type Image struct {
	ID             bson.ObjectId  `json:"id" bson:"_id"`
	Description    string         `json:"description"`
	Copyright      string         `json:"copyright"`
	Height         uint           `json:"height"`
	Filetype       string         `json:"filetype"`
	Width          uint           `json:"width"`
	URL            string         `json:"url"`
	ResizedTargets ResizedTargets `json:"resized_targets"`
}

type MongoImage struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Description string
	Copyright   string
	Image       struct {
		Height         uint
		Filetype       string `bson:"filetype"`
		Width          uint
		URL            string
		ResizedTargets ResizedTargets `bson:"resizedTargets" json:"resized_targets"`
	}
}

func (mi *MongoImage) ToImage() (img Image) {
	img.ID = mi.ID
	img.Description = mi.Description
	img.Copyright = mi.Copyright
	img.Height = mi.Image.Height
	img.Width = mi.Image.Width
	img.Filetype = mi.Image.Filetype
	img.URL = mi.Image.URL
	img.ResizedTargets = mi.Image.ResizedTargets
	return
}

// Brief ...
type Brief struct {
	HTML    string   `json:"html"`
	APIData []bson.M `bson:"apiData" json:"api_data"`
}

// Category ...
type Category struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	SortOrder uint          `json:"sort_order"`
	Name      string        `json:"name"`
}

// Tag ...
type Tag struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `json:"name"`
}

// PostMeta ...
type PostMeta struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	Slug          string        `json:"slug"`
	Name          string        `json:"name"`
	Subtitle      string        `json:"subtitle"`
	State         string        `json:"state"`
	HeroImage     interface{}   `bson:"heroImage" json:"hero_image"`
	Brief         `json:"brief"`
	Categories    []interface{} `bson:"categories" json:"categories"`
	Style         string        `json:"style"`
	Bookmark      string        `json:"bookmark"`
	Copyright     string        `json:"copyright"`
	Tags          []interface{} `json:"tags"`
	OgDescription string        `bson:"og_description" json:"og_description"`
	OgImage       interface{}   `bson:"og_image" json:"og_image"`
	IsFeatured    bool          `bson:"isFeatured" json:"is_featured"`
	PublishedDate time.Time     `bson:"publishedDate" json:"published_date"`
	Topic         interface{}   `bson:"topics,omitempty" json:"topic"`
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
