package news

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetaOfTopic struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Slug  string             `bson:"slug" json:"slug"`
	Title string             `bson:"title" json:"title"`
	// TODO: rename the bson field to short_title
	ShortTitle           string               `bson:"topic_name" json:"short_title"`
	PublishedDate        time.Time            `bson:"publishedDate" json:"published_date"`
	OgDescription        string               `bson:"og_description" json:"og_description"`
	OgImage              *Image               `bson:"og_image" json:"og_image,omitempty"`
	LeadingImage         *Image               `bson:"leading_image" json:"leading_image,omitempty"`
	LeadingImagePortrait *Image               `bson:"leading_image_portrait" json:"leading_image_portrait,omitempty"`
	Relateds             []primitive.ObjectID `bson:"relateds" json:"relateds,omitempty"`
	Full                 bool                 `bson:"-" json:"full"`
}

type Topic struct {
	// Use inline tag for unflattened the response document to unmarshal into embedded struct
	// https://godoc.org/go.mongodb.org/mongo-driver/bson#hdr-Structs
	MetaOfTopic        `bson:",inline"`
	RelatedsBackground string       `bson:"relateds_background" json:"relateds_background"`
	RelatedsFormat     string       `bson:"relateds_format" json:"relateds_format"`
	TitlePosition      string       `bson:"title_position" json:"title_position"`
	LeadingVideo       *Video       `bson:"leading_video" json:"leading_video,omitempty"`
	Headline           string       `bson:"headline" json:"headline"`
	Subtitle           string       `bson:"subtitle" json:"subtitle"`
	Description        *ContentBody `bson:"description,omitempty" json:"description,omitempty"`
	TeamDescription    *ContentBody `bson:"team_description,omitempty" json:"team_description,omitempty"`
	OgTitle            string       `bson:"og_title" json:"og_title"`
}

type ContentBody struct {
	APIData []primitive.M `bson:"apiData" json:"api_data"`
}

type MetaOfPost struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	Style                string             `bson:"style" json:"style"`
	Slug                 string             `bson:"slug" json:"slug"`
	LeadingImagePortrait *Image             `bson:"leading_image_portrait" json:"leading_image_portrait,omitempty"`
	HeroImage            *Image             `bson:"heroImage" json:"hero_image,omitempty"`
	OgImage              *Image             `bson:"og_image" json:"og_image,omitempty"`
	OgDescription        string             `bson:"og_description" json:"og_description"`
	Title                string             `bson:"title" json:"title"`
	Subtitle             string             `bson:"subtitle" json:"subtitle"`
	Categories           []category         `bson:"categories" json:"categories,omitempty"`
	PublishedDate        time.Time          `bson:"publishedDate" json:"published_date"`
	IsExternal           bool               `bson:"is_external" json:"is_external"`
	Tags                 []Tag              `bson:"tags" json:"tags,omitempty"`
	Full                 bool               `bson:"-" json:"full"`
}

type Post struct {
	// Use inline tag for unflattened the response document to unmarshal into embedded struct
	// https://godoc.org/go.mongodb.org/mongo-driver/bson#hdr-Structs
	MetaOfPost             `bson:",inline"`
	Brief                  *ContentBody         `bson:"brief,omitempty" json:"brief,omitempty"`
	Content                *ContentBody         `bson:"content,omitempty" json:"content,omitempty"`
	Copyright              string               `bson:"copyright" json:"copyright"`
	Designers              []MetaOfAuthor       `bson:"designers" json:"designers,omitempty"`
	Engineers              []MetaOfAuthor       `bson:"engineers" json:"engineers,omitempty"`
	ExtendByline           string               `bson:"extend_byline" json:"extend_byline"`
	LeadingImageDecription string               `bson:"leading_image_description" json:"leading_image_description"`
	OgTitle                string               `bson:"og_title" json:"og_title"`
	Photographers          []MetaOfAuthor       `bson:"photographers" json:"photographers,omitempty"`
	Relateds               []primitive.ObjectID `bson:"relateds" json:"relateds,omitempty"`
	// TODO: rename the bson field to `topic`
	// Define inline struct here so that the nested projection
	// won't be performed on topic to exclude the unwanted fields.
	Topic struct {
		Title string `bson:"title" json:"title"`
		// TODO: rename the bson field to short_title
		ShortTitle string               `bson:"topic_name" json:"short_title"`
		Slug       string               `bson:"slug" json:"slug"`
		Relateds   []primitive.ObjectID `bson:"relateds" json:"relateds,omitempty"`
	} `bson:"topics" json:"topic,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updated_at"`
	// TODO: rename the bson field to `writers`
	Writers       []MetaOfAuthor `bson:"writters" json:"writers,omitempty"`
	HeroImageSize string         `bson:"heroImageSize" json:"hero_image_size"`
}

type MetaOfAuthor struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	JobTitle string             `bson:"job_title" json:"job_title"`
	Name     string             `bson:"name" json:"name"`
}

type Author struct {
	MetaOfAuthor `bson:",inline"`
	Email        string    `bson:"email" json:"email"`
	Bio          string    `bson:"bio" json:"bio"`
	Thumbnail    *Image    `bson:"thumbnail" json:"thumbnail"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

type category struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SortOrder uint               `bson:"sort_order" json:"sort_order"`
	Name      string             `bson:"name" json:"name"`
}

type Image struct {
	ImageMeta   `bson:"image"`
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Description string             `json:"description,omitempty"`
}

type ImageMeta struct {
	Filetype       string         `bson:"filetype" json:"filetype"`
	ResizedTargets ResizedTargets `bson:"resizedTargets" json:"resized_targets"`
}

type ImageAsset struct {
	Height uint   `bson:"height" json:"height"`
	Width  uint   `bson:"width" json:"width"`
	URL    string `bson:"url" json:"url"`
}

type ResizedTargets struct {
	Mobile  ImageAsset `bson:"mobile" json:"mobile"`
	Tiny    ImageAsset `bson:"tiny" json:"tiny"`
	Desktop ImageAsset `bson:"desktop" json:"desktop"`
	Tablet  ImageAsset `bson:"tablet" json:"tablet"`
	W400    ImageAsset `bson:"w400" json:"w400"`
}

type Tag struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type VideoMeta struct {
	Filetype string `bson:"filetype" json:"filetype"`
	Size     uint   `bson:"size" json:"size"`
	URL      string `bson:"url" json:"url"`
}

type Video struct {
	VideoMeta `bson:"video"`
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title"`
}
