package models

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
)

// ImageAsset ...
type ImageAsset struct {
	Height uint   `json:"height,omitempty"`
	Width  uint   `json:"width,omitempty"`
	URL    string `json:"url,omitempty"`
}

// ResizedTargets ...
type ResizedTargets struct {
	Mobile  ImageAsset `json:"mobile,omitempty"`
	Tiny    ImageAsset `json:"tiny,omitempty"`
	Desktop ImageAsset `json:"desktop,omitempty"`
	Tablet  ImageAsset `json:"tablet,omitempty"`
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

// MongoImage is the data structure  returned by Mongo DB
type MongoImage struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Description string
	Copyright   string
	Image       struct {
		Height         uint
		Filetype       string `json:"filetype"`
		Width          uint
		URL            string
		ResizedTargets ResizedTargets `bson:"resizedTargets" json:"resized_targets"`
	}
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

// MongoVideo is the data structure returned by Mongo DB
type MongoVideo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Title string
	Video struct {
		Filetype string `json:"filetype"`
		Size     uint
		URL      string
	}
}

// Video is used to return in response
type Video struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Title    string
	Filetype string `json:"filetype"`
	Size     uint
	URL      string
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
	SortOrder uint          `json:"sort_order"`
	Name      string        `json:"name"`
}

// Tag ...
type Tag struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `json:"name"`
}

// NewsEntity defines the method of structs such `Topic`, `PostMeta` ...etc
type NewsEntity interface {
	SetEmbeddedAsset(string, interface{})
	GetHeroImageOrigin() bson.ObjectId
	GetLeadingImageOrigin() bson.ObjectId
	GetOgImageOrigin() bson.ObjectId
	GetCategoriesOrigin() []bson.ObjectId
	GetTagsOrigin() []bson.ObjectId
	GetTopicMetaOrigin() bson.ObjectId
	GetTopicOrigin() bson.ObjectId
	GetLeadingImagePortraitOrigin() bson.ObjectId
	GetLeadingVideoOrigin() bson.ObjectId
}

// SetEmbeddedAsset - `PostMeta` implements this method to become a `NewsEntity`
func (pm *PostMeta) SetEmbeddedAsset(key string, asset interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Set field %v into PostMeta occurs error", key)
		}
	}()
	// pointer to struct - addressable
	ps := reflect.ValueOf(pm)
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

// GetTopicMetaOrigin ...
func (pm *PostMeta) GetTopicMetaOrigin() bson.ObjectId {
	return pm.TopicMetaOrigin
}

// GetTopicOrigin ...
func (pm *PostMeta) GetTopicOrigin() bson.ObjectId {
	return ""
}

// GetLeadingImagePortraitOrigin ...
func (pm *PostMeta) GetLeadingImagePortraitOrigin() bson.ObjectId {
	return ""
}

// GetLeadingVideoOrigin ...
func (pm *PostMeta) GetLeadingVideoOrigin() bson.ObjectId {
	return ""
}

// SetEmbeddedAsset - `Topic` implements this method to become a `NewsEntity`
func (t *Topic) SetEmbeddedAsset(key string, asset interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Set field %v into PostMeta occurs error", key)
		}
	}()
	// pointer to struct - addressable
	ps := reflect.ValueOf(t)
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

// GetTopicMetaOrigin ...
func (t *Topic) GetTopicMetaOrigin() bson.ObjectId {
	return ""
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
