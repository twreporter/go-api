package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TopicMeta struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	Name          string        `json:"name"`
	TopicName     string        `json:"topic_name"`
	Title         string        `json:"title"`
	Subtitle      string        `json:"subtitle"`
	Headline      string        `json:"headline"`
	PublishedDate time.Time     `bson:"publishedDate" json:"published_date"`
	Description   string        `json:"description"`
	OgTitle       string        `bson:"og_title" json:"og_title"`
	OgDescription string        `bson:"og_description" json:"og_description"`
}

/*
type Topic struct {
	ID                   bson.ObjectId `json:"id" bson:"_id"`
	Slug                 string        `json:"slug"`
	TopicName            string        `json:"topic_name"`
	Title                string        `json:"title"`
	TitlePosition        string        `json:"title_position"`
	Subtitle             string        `json:"subtitle"`
	Headline             string        `json:"headline"`
	State                string        `json:"state"`
	PublishedDate        time.Time     `json:"published_date"`
	Description          Brief         `json:"description"`
	TeamDescription      Brief         `json:"team_description"`
	Relateds             []PostMeta    `json:"relateds"`
	RelatedsFormat       string        `json:"relateds_format"`
	RelatedsBackground   string        `json:"relateds_background"`
	LeadingImage         Image         `json:"leading_image"`
	LeadingImagePortrait Image         `json:"leading_image_portrait"`
	LeadingVideo         Video         `json:"leading_video"`
	OgTitle              string        `json:"og_title"`
	OgDescription        string        `json:"og_description"`
	OgImage              Image         `json:"og_image"`
}
*/
