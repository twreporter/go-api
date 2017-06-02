package storage

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// GetMetaOfPosts is a type-specific functions implementing the method defined in the NewsStorage.
// It parses query string into bson and finds the posts according to that bson.
func (m *MongoStorage) GetMetaOfPosts(qs string, limit int, offset int, sort string, embedded []string) ([]models.PostMeta, error) {
	var q models.MongoQuery
	var posts []models.PostMeta

	err := models.GetQuery(qs, &q)
	log.Info("q:", q)

	if err != nil {
		log.Info("Parse query param occurs error: ", err.Error())
	}

	err = m.db.DB("plate").C("posts").Find(q).Limit(limit).Skip(offset).Sort(sort).All(&posts)

	for index, post := range posts {
		post = m.GetPostEmbeddedAsset(post, embedded)
		posts[index] = post
	}

	if err != nil {
		return posts, m.NewStorageError(err, "GetPosts", "storage.posts.get_posts")
	}

	return posts, nil
}

// GetPostEmbeddedAsset ...
func (m *MongoStorage) GetPostEmbeddedAsset(post models.PostMeta, embedded []string) models.PostMeta {
	if embedded != nil {
		for _, ele := range embedded {
			switch ele {
			case "hero_image":
				img, err := m.GetImage(post.HeroImageOrigin)
				if err == nil {
					post.HeroImage = &img
				}
				break
			case "og_image":
				img, err := m.GetImage(post.OgImageOrigin)
				if err == nil {
					post.OgImage = &img
				}
				break
			case "categories":
				categories, _ := m.GetCategories(post.CategoriesOrigin)
				post.Categories = make([]models.Category, len(categories))
				for i, v := range categories {
					post.Categories[i] = v
				}
				break
			case "tags":
				tags, _ := m.GetTags(post.TagsOrigin)
				post.Tags = make([]models.Tag, len(tags))
				for i, v := range tags {
					post.Tags[i] = v
				}
				break
			case "topic":
				topic, err := m.GetTopicMeta(post.TopicOrigin)
				if err == nil {
					post.Topic = &topic
				}
				break
			default:
				log.Info(fmt.Sprintf("Embedded element (%v) is not supported: ", ele))
			}
		}
	}
	return post
}

// GetCategories ...
func (m *MongoStorage) GetCategories(ids []bson.ObjectId) ([]models.Category, error) {
	var cats []models.Category

	if ids == nil {
		return cats, nil
	}

	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	err := m.db.DB("plate").C("postcategories").Find(query).All(&cats)
	if err != nil {
		return nil, m.NewStorageError(err, "GetCategories", "storage.posts.get_categories")
	}

	return cats, nil
}

// GetTags ...
func (m *MongoStorage) GetTags(ids []bson.ObjectId) ([]models.Tag, error) {
	var tags []models.Tag

	if ids == nil {
		return tags, nil
	}

	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	err := m.db.DB("plate").C("tags").Find(query).All(&tags)
	if err != nil {
		return nil, m.NewStorageError(err, "GetCategories", "storage.posts.get_tags")
	}

	return tags, nil
}

// GetTopicMeta ...
func (m *MongoStorage) GetTopicMeta(id bson.ObjectId) (models.TopicMeta, error) {
	var tm models.TopicMeta

	if id == "" {
		return tm, models.NewAppError("GetTopicMeta", "storage.posts.get_meta_of_topic.id_not_provided", "Resource not found", http.StatusNotFound)
	}

	err := m.db.DB("plate").C("topics").FindId(id).One(&tm)
	if err != nil {
		return tm, m.NewStorageError(err, "GetTopicMeta", "storage.posts.get_meta_of_topic.error")
	}

	return tm, nil
}

/*
func (g *MongoStorage) GetVideo(id interface{}) (models.Video,error) {
}
*/

// GetImage ...
func (m *MongoStorage) GetImage(id bson.ObjectId) (models.Image, error) {
	var mgoImg models.MongoImage

	if id == "" {
		return models.Image{}, models.NewAppError("GetImage", "storage.posts.get_image.id_not_provided", "Resource not found", http.StatusNotFound)
	}

	err := m.db.DB("plate").C("images").FindId(id).One(&mgoImg)
	if err != nil {
		return models.Image{}, m.NewStorageError(err, "GetImage", "storage.posts.get_image.error")
	}

	return mgoImg.ToImage(), nil
}
