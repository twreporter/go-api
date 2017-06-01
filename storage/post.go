package storage

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// GetPosts is type-specific functions to implement the method defined in the NewsStorage.
// It parses query string into bson and finds the posts according to that bson.
func (g *MongoStorage) GetMetaOfPosts(qs string, limit int, offset int, embedded []string) ([]models.PostMeta, error) {
	var q models.MongoQuery
	var posts []models.PostMeta

	err := models.GetQuery(qs, &q)
	log.Info("q:", q)

	if err != nil {
		log.Info("Parse query param occurs error: ", err.Error())
	}

	err = g.db.DB("plate").C("posts").Find(q).Limit(limit).Skip(offset).All(&posts)

	for index, post := range posts {
		if post.HeroImage != nil && post.HeroImage != "" {
			post = g.GetPostEmbeddedAsset(post, embedded)
			posts[index] = post
		}
	}

	if err != nil {
		return posts, g.NewStorageError(err, "GetPosts", "storage.posts.get_posts")
	}

	return posts, nil
}

func (g *MongoStorage) GetPostEmbeddedAsset(post models.PostMeta, embedded []string) models.PostMeta {
	if embedded != nil {
		for _, ele := range embedded {
			switch ele {
			case "hero_image":
				heroImage, err := g.GetImage(post.HeroImage)
				if err == nil {
					post.HeroImage = heroImage
				}
				break
			case "og_image":
				image, err := g.GetImage(post.OgImage)
				if err == nil {
					post.OgImage = image
				}
				break
			case "categories":
				categories, _ := g.GetCategories(post.Categories)
				post.Categories = make([]interface{}, len(categories))
				for i, v := range categories {
					post.Categories[i] = v
				}
				break
			case "tags":
				tags, _ := g.GetTags(post.Tags)
				post.Tags = make([]interface{}, len(tags))
				for i, v := range tags {
					post.Tags[i] = v
				}
				break
			case "topic":
				topic, err := g.GetTopicMeta(post.Topic)
				if err == nil {
					post.Topic = topic
				}
				break
			default:
				log.Info(fmt.Sprintf("Embedded element (%v) is not supported: ", ele))
			}
		}
	}
	return post
}

func (g *MongoStorage) GetCategories(ids []interface{}) ([]models.Category, error) {
	var cats []models.Category

	if ids == nil {
		return cats, nil
	}

	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	err := g.db.DB("plate").C("postcategories").Find(query).All(&cats)
	if err != nil {
		return nil, g.NewStorageError(err, "GetCategories", "storage.posts.get_categories")
	}

	return cats, nil
}

func (g *MongoStorage) GetTags(ids []interface{}) ([]models.Tag, error) {
	var tags []models.Tag

	if ids == nil {
		return tags, nil
	}

	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	err := g.db.DB("plate").C("tags").Find(query).All(&tags)
	if err != nil {
		return nil, g.NewStorageError(err, "GetCategories", "storage.posts.get_tags")
	}

	return tags, nil
}

func (g *MongoStorage) GetTopicMeta(id interface{}) (models.TopicMeta, error) {
	var tm models.TopicMeta

	if id == nil || id == "" {
		return tm, models.NewAppError("GetTopicMeta", "storage.posts.get_meta_of_topic.id_not_provided", "Resource not found", http.StatusNotFound)
	}

	err := g.db.DB("plate").C("topics").FindId(id).One(&tm)
	if err != nil {
		return tm, g.NewStorageError(err, "GetTopicMeta", "storage.posts.get_meta_of_topic.error")
	}

	return tm, nil
}

/*
func (g *MongoStorage) GetVideo(id interface{}) (models.Video,error) {
}
*/

func (g *MongoStorage) GetImage(id interface{}) (models.Image, error) {
	var mgoImg models.MongoImage

	if id == nil || id == "" {
		return models.Image{}, models.NewAppError("GetImage", "storage.posts.get_image.id_not_provided", "Resource not found", http.StatusNotFound)
	}

	err := g.db.DB("plate").C("images").FindId(id).One(&mgoImg)
	if err != nil {
		return models.Image{}, g.NewStorageError(err, "GetImage", "storage.posts.get_image.error")
	}

	return mgoImg.ToImage(), nil
}
