package models

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
)

// ImageAsset ...
type ImageAsset struct {
	Height uint   `bson:"height" json:"height"`
	Width  uint   `bson:"width" json:"width"`
	URL    string `bson:"url" json:"url"`
}

// ResizedTargets ...
type ResizedTargets struct {
	Mobile  ImageAsset `bson:"mobile" json:"mobile"`
	Tiny    ImageAsset `bson:"tiny" json:"tiny"`
	Desktop ImageAsset `bson:"desktop" json:"desktop"`
	Tablet  ImageAsset `bson:"tablet" json:"tablet"`
}

// Image is used to return in response
type Image struct {
	ID             bson.ObjectId  `json:"id,omitempty" bson:"_id,omitempty"`
	Description    string         `json:"description,omitempty"`
	Copyright      string         `json:"copyright,omitempty"`
	Height         uint           `json:"height,omitempty"`
	Filetype       string         `json:"filetype,omitempty"`
	Width          uint           `json:"width,omitempty"`
	URL            string         `json:"url,omitempty"`
	ResizedTargets ResizedTargets `json:"resized_targets,omitempty"`
}

type MongoImageAsset struct {
	Height         uint           `bson:"height" json:"height"`
	Filetype       string         `bson:""filetype json:"filetype"`
	Width          uint           `bson:"width" json:"width"`
	URL            string         `bson:"url" json:"url"`
	ResizedTargets ResizedTargets `bson:"resizedTargets" json:"resized_targets"`
}

// MongoImage is the data structure  returned by Mongo DB
type MongoImage struct {
	ID          bson.ObjectId   `bson:"_id"`
	Description string          `bson:"description"`
	Copyright   string          `bson:"copyright"`
	Image       MongoImageAsset `bson:"image"`
}

// ToImage transforms MongoImage to Image
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

type MongoVideoAsset struct {
	Filetype string `bson:"filetype"`
	Size     uint   `bson:"size"`
	URL      string `bson:"url"`
}

// MongoVideo is the data structure returned by Mongo DB
type MongoVideo struct {
	ID    bson.ObjectId   `bson:"_id"`
	Title string          `bson:"title"`
	Video MongoVideoAsset `bson:"video"`
}

// Video is used to return in response
type Video struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Title    string        `json:"title"`
	Filetype string        `json:"filetype"`
	Size     uint          `json:"size"`
	URL      string        `json:"url"`
}

// ToVideo Transform MongoVideo to Video
func (mv *MongoVideo) ToVideo() (video Video) {
	video.ID = mv.ID
	video.Title = mv.Title
	video.Filetype = mv.Video.Filetype
	video.Size = mv.Video.Size
	video.URL = mv.Video.URL
	return
}

// Category ...
type Category struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	SortOrder uint          `bson:"sort_order" json:"sort_order"`
	Name      string        `bson:"name" json:"name"`
}

// Tag ...
type Tag struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

// NewsEntity defines the method of structs such `Topic`, `PostMeta` ...etc
type NewsEntity interface {
	SetEmbeddedAsset(string, interface{})
	GetHeroImageOrigin() bson.ObjectId
	GetLeadingImageOrigin() bson.ObjectId
	GetOgImageOrigin() bson.ObjectId
	GetCategoriesOrigin() []bson.ObjectId
	GetTagsOrigin() []bson.ObjectId
	GetTopicOrigin() bson.ObjectId
	GetLeadingImagePortraitOrigin() bson.ObjectId
	GetLeadingVideoOrigin() bson.ObjectId
	GetRelatedsOrigin() []bson.ObjectId
}

func __setEmbeddedAsset(strt interface{}, key string, asset interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Set field %v into PostMeta occurs error", key)
		}
	}()
	// pointer to struct - addressable
	ps := reflect.ValueOf(strt)
	// struct
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// exported field
		f := s.FieldByName(key)
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() {
				// change value of N
				f.Set(reflect.ValueOf(asset))
			}
		}
	}
}

// SetEmbeddedAsset - `PostMeta` implements this method to become a `NewsEntity`
func (pm *PostMeta) SetEmbeddedAsset(key string, asset interface{}) {
	__setEmbeddedAsset(pm, key, asset)
}

// GetHeroImageOrigin ...
func (pm *PostMeta) GetHeroImageOrigin() bson.ObjectId {
	return pm.HeroImageOrigin
}

// GetLeadingImageOrigin ...
func (pm *PostMeta) GetLeadingImageOrigin() bson.ObjectId {
	return ""
}

// GetOgImageOrigin ...
func (pm *PostMeta) GetOgImageOrigin() bson.ObjectId {
	return pm.OgImageOrigin
}

// GetCategoriesOrigin ...
func (pm *PostMeta) GetCategoriesOrigin() []bson.ObjectId {
	return pm.CategoriesOrigin
}

// GetTagsOrigin ...
func (pm *PostMeta) GetTagsOrigin() []bson.ObjectId {
	return pm.TagsOrigin
}

// GetTopicOrigin ...
func (pm *PostMeta) GetTopicOrigin() bson.ObjectId {
	return pm.TopicOrigin
}

// GetLeadingImagePortraitOrigin ...
func (pm *PostMeta) GetLeadingImagePortraitOrigin() bson.ObjectId {
	return ""
}

// GetLeadingVideoOrigin ...
func (pm *PostMeta) GetLeadingVideoOrigin() bson.ObjectId {
	return ""
}

// GetRelatedsOrigin ...
func (pm *PostMeta) GetRelatedsOrigin() []bson.ObjectId {
	return nil
}

// SetEmbeddedAsset - `Topic` implements this method to become a `NewsEntity`
func (t *Topic) SetEmbeddedAsset(key string, asset interface{}) {
	__setEmbeddedAsset(t, key, asset)
}

// GetHeroImageOrigin ...
func (t *Topic) GetHeroImageOrigin() bson.ObjectId {
	return ""
}

// GetLeadingImageOrigin ...
func (t *Topic) GetLeadingImageOrigin() bson.ObjectId {
	return t.LeadingImageOrigin
}

// GetOgImageOrigin ...
func (t *Topic) GetOgImageOrigin() bson.ObjectId {
	return t.OgImageOrigin
}

// GetCategoriesOrigin ...
func (t *Topic) GetCategoriesOrigin() []bson.ObjectId {
	return nil
}

// GetTagsOrigin ...
func (t *Topic) GetTagsOrigin() []bson.ObjectId {
	return nil
}

// GetTopicOrigin ...
func (t *Topic) GetTopicOrigin() bson.ObjectId {
	return ""
}

// GetLeadingImagePortraitOrigin ...
func (t *Topic) GetLeadingImagePortraitOrigin() bson.ObjectId {
	return t.LeadingImagePortraitOrigin
}

// GetLeadingVideoOrigin ...
func (t *Topic) GetLeadingVideoOrigin() bson.ObjectId {
	return t.LeadingVideoOrigin
}

// GetRelatedsOrigin ...
func (t *Topic) GetRelatedsOrigin() []bson.ObjectId {
	return t.RelatedsOrigin
}
