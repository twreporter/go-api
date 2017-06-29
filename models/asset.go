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

type Author struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
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

// NewsEntity defines the method of structs such `Topic`, `Post` ...etc
type NewsEntity interface {
	SetEmbeddedAsset(string, interface{})
	GetEmbeddedAsset(string) []bson.ObjectId
}

func __setEmbeddedAsset(strt interface{}, key string, asset interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Set field %v into Post occurs error", key)
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

// SetEmbeddedAsset - `Post` implements this method to become a `NewsEntity`
func (pm *Post) SetEmbeddedAsset(key string, asset interface{}) {
	__setEmbeddedAsset(pm, key, asset)
}

// GetEmbeddedAsset - `Post` implements this method to become a `NewsEntity`
func (pm *Post) GetEmbeddedAsset(key string) []bson.ObjectId {
	var rtn []bson.ObjectId
	switch key {
	case "WrittersOrigin":
		return pm.WrittersOrigin
	case "PhotographersOrigin":
		return pm.PhotographersOrigin
	case "DesignersOrigin":
		return pm.DesignersOrigin
	case "EngineersOrigin":
		return pm.EngineersOrigin
	case "HeroImageOrigin":
		if pm.HeroImageOrigin != "" {
			return append(rtn, pm.HeroImageOrigin)
		}
		return nil
	case "LeadingImageOrigin":
		return nil
	case "CategoriesOrigin":
		return pm.CategoriesOrigin
	case "TagsOrigin":
		return pm.TagsOrigin
	case "OgImageOrigin":
		if pm.OgImageOrigin != "" {
			return append(rtn, pm.OgImageOrigin)
		}
		return nil
	case "LeadingVideoOrigin":
		if pm.LeadingVideoOrigin != "" {
			return append(rtn, pm.LeadingVideoOrigin)
		}
		return nil
	case "LeadingImagePortraitOrigin":
		return nil
	case "TopicOrigin":
		if pm.TopicOrigin != "" {
			return append(rtn, pm.TopicOrigin)
		}
		return nil
	case "RelatedsOrigin":
		return pm.RelatedsOrigin
	default:
		return nil
	}
}

// SetEmbeddedAsset - `Topic` implements this method to become a `NewsEntity`
func (t *Topic) SetEmbeddedAsset(key string, asset interface{}) {
	__setEmbeddedAsset(t, key, asset)
}

// GetEmbeddedAsset - `Topic` implements this method to become a `NewsEntity`
func (t *Topic) GetEmbeddedAsset(key string) []bson.ObjectId {
	var rtn []bson.ObjectId
	switch key {
	case "WrittersOrigin":
		return nil
	case "PhotographersOrigin":
		return nil
	case "DesignersOrigin":
		return nil
	case "EngineersOrigin":
		return nil
	case "HeroImageOrigin":
		return nil
	case "LeadingImageOrigin":
		if t.LeadingImageOrigin != "" {
			return append(rtn, t.LeadingImageOrigin)
		}
		return nil
	case "CategoriesOrigin":
		return nil
	case "TagsOrigin":
		return nil
	case "OgImageOrigin":
		if t.OgImageOrigin != "" {
			return append(rtn, t.OgImageOrigin)
		}
		return nil
	case "LeadingVideoOrigin":
		if t.LeadingVideoOrigin != "" {
			return append(rtn, t.LeadingVideoOrigin)
		}
		return nil
	case "LeadingImagePortraitOrigin":
		if t.LeadingImagePortraitOrigin != "" {
			return append(rtn, t.LeadingImagePortraitOrigin)
		}
		return nil
	case "TopicOrigin":
		return nil
	case "RelatedsOrigin":
		return t.RelatedsOrigin
	default:
		return nil
	}
}
