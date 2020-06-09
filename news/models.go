package news

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id"`
	Slug                 string             `bson:"slug" json:"slug"`
	Name                 string             `bson:"name" json:"name"`
	TopicName            string             `bson:"topic_name" json:"topic_name"`
	Title                string             `bson:"title" json:"title"`
	TitlePosition        string             `bson:"title_position" json:"title_position"`
	Subtitle             string             `bson:"subtitle" json:"subtitle"`
	Headline             string             `bson:"headline" json:"headline"`
	State                string             `bson:"state" json:"state"`
	Description          *ContentBody       `bson:"description,omitempty" json:"description,omitempty"`
	TeamDescription      *ContentBody       `bson:"team_description,omitempty" json:"team_description,omitempty"`
	Relateds             []Post             `bson:"relateds" json:"relateds,omitempty"`
	RelatedsFormat       string             `bson:"relateds_format" json:"relateds_format"`
	RelatedsBackground   string             `bson:"relateds_background" json:"relateds_background"`
	LeadingImage         *Image             `bson:"leading_image" json:"leading_image,omitempty"`
	LeadingImagePortrait *Image             `bson:"leading_image_portrait" json:"leading_image_portrait,omitempty"`
	LeadingVideo         *Video             `bson:"leading_video" json:"leading_video,omitempty"`
	OgTitle              string             `bson:"og_title" json:"og_title"`
	OgDescription        string             `bson:"og_description" json:"og_description"`
	OgImage              *Image             `bson:"og_image" json:"og_image,omitempty"`
	PublishedDate        time.Time          `bson:"publishedDate" json:"published_date"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updated_at"`
	Full                 bool               `bson:"-" json:"full"`
}

type ContentBody struct {
	APIData []primitive.M `bson:"apiData" json:"api_data"`
}

type Post struct {
	ID                     primitive.ObjectID `bson:"_id" json:"id"`
	Slug                   string             `bson:"slug" json:"slug"`
	Name                   string             `bson:"name" json:"name"`
	Title                  string             `bson:"title" json:"title"`
	Subtitle               string             `bson:"subtitle" json:"subtitle"`
	State                  string             `bson:"state" json:"state"`
	HeroImage              *Image             `bson:"heroImage" json:"hero_image,omitempty"`
	HeroImageSize          string             `bson:"heroImageSize" json:"hero_image_size"`
	LeadingImagePortrait   *Image             `bson:"leading_image_portrait" json:"leading_image_portrait,omitempty"`
	LeadingImageDecription string             `bson:"leading_image_description" json:"leading_image_description"`
	Brief                  *ContentBody       `bson:"brief,omitempty" json:"brief,omitempty"`
	Categories             []Category         `bson:"categories" json:"categories,omitempty"`
	Style                  string             `bson:"style" json:"style"`
	Theme                  *Theme             `bson:"theme" json:"theme"`
	Copyright              string             `bson:"copyright" json:"copyright"`
	Tags                   []Tag              `bson:"tags" json:"tags,omitempty"`
	OgTitle                string             `bson:"og_title" json:"og_title"`
	OgDescription          string             `bson:"og_description" json:"og_description"`
	OgImage                *Image             `bson:"og_image" json:"og_image,omitempty"`
	IsFeatured             bool               `bson:"isFeatured" json:"is_featured"`
	Topic                  *Topic             `bson:"topics" json:"topics,omitempty"`
	Writters               []Author           `bson:"writters" json:"writters,omitempty"`
	Photographers          []Author           `bson:"photographers" json:"photographers,omitempty"`
	Designers              []Author           `bson:"designers" json:"designers,omitempty"`
	Engineers              []Author           `bson:"engineers" json:"engineers,omitempty"`
	ExtendByline           string             `bson:"extend_byline" json:"extend_byline"`
	LeadingVideo           *Video             `bson:"leading_video" json:"leading_video,omitempty"`
	Content                *ContentBody       `bson:"content,omitempty" json:"content,omitempty"`
	Relateds               []Post             `bson:"relateds" json:"relateds,omitempty"`
	PublishedDate          time.Time          `bson:"publishedDate" json:"published_date"`
	UpdatedAt              time.Time          `bson:"updatedAt" json:"updated_at"`
	Full                   bool               `bson:"-" json:"full"`
	IsExternal             bool               `bson:"is_external" json:"is_external"`
}

type Author struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	JobTitle string             `bson:"job_title" json:"job_title"`
	Name     string             `bson:"name" json:"name"`
}

type Category struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SortOrder uint               `bson:"sort_order" json:"sort_order"`
	Name      string             `bson:"name" json:"name"`
}

// Theme ...
/*type Image struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Description    string             `json:"description,omitempty"`
	Copyright      string             `json:"copyright,omitempty"`
	Height         uint               `json:"height,omitempty"`
	Filetype       string             `json:"filetype,omitempty"`
	Width          uint               `json:"width,omitempty"`
	URL            string             `json:"url,omitempty"`
	ResizedTargets ResizedTargets     `json:"resized_targets,omitempty" bson:"resizedTargets"`
}*/
type Image struct {
	ImageMeta   `bson:"image"`
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Description string             `json:"description,omitempty"`
	Copyright   string             `json:"copyright,omitempty"`
}

type ImageMeta struct {
	Height         uint           `bson:"height" json:"height"`
	Filetype       string         `bson:"filetype" json:"filetype"`
	Width          uint           `bson:"width" json:"width"`
	URL            string         `bson:"url" json:"url"`
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

type Theme struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Name            string             `bson:"name" json:"name"`
	TitlePosition   string             `bson:"title_position" json:"title_position"`
	HeaderPosition  string             `bson:"header_position" json:"header_position"`
	TitleColor      string             `bson:"title_color" json:"title_color"`
	SubtitleColor   string             `bson:"subtitle_color" json:"subtitle_color"`
	TopicColor      string             `bson:"topic_color" json:"topic_color"`
	FontColor       string             `bson:"font_color" json:"font_color"`
	BackgroundColor string             `bson:"bg_color" json:"bg_color"`
	FooterBGColor   string             `bson:"footer_bg_color" json:"footer_bg_color"`
	LogoColor       string             `bson:"logo_color" json:"logo_color"`
}

type Video struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Title    string             `json:"title"`
	Filetype string             `json:"filetype"`
	Size     uint               `json:"size"`
	URL      string             `json:"url"`
}
