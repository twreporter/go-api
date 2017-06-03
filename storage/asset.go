package storage

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// GetPostEmbeddedAsset ...
func (m *MongoStorage) GetEmbeddedAsset(entity models.NewsEntity, embedded []string) {
	if embedded != nil {
		for _, ele := range embedded {
			switch ele {
			case "hero_image":
				img, err := m.GetImage(entity.GetHeroImageOrigin())
				if err == nil {
					entity.SetEmbeddedAsset("HeroImage", &img)
				}
				break
			case "leading_image":
				img, err := m.GetImage(entity.GetLeadingImageOrigin())
				if err == nil {
					entity.SetEmbeddedAsset("LeadingImage", &img)
				}
				break
			case "leading_image_portrait":
				img, err := m.GetImage(entity.GetLeadingImagePortraitOrigin())
				if err == nil {
					entity.SetEmbeddedAsset("LeadingImagePortrait", &img)
				}
				break
			case "og_image":
				img, err := m.GetImage(entity.GetOgImageOrigin())
				if err == nil {
					entity.SetEmbeddedAsset("OgImage", &img)
				}
				break
			case "categories":
				categories, _ := m.GetCategories(entity.GetCategoriesOrigin())
				_categories := make([]models.Category, len(categories))
				for i, v := range categories {
					_categories[i] = v
				}
				entity.SetEmbeddedAsset("Categories", _categories)
				break
			case "tags":
				tags, _ := m.GetTags(entity.GetTagsOrigin())
				_tags := make([]models.Tag, len(tags))
				for i, v := range tags {
					_tags[i] = v
				}
				entity.SetEmbeddedAsset("Tags", _tags)
				break
			case "topic_meta":
				var t models.TopicMeta
				err := m.GetDocument(entity.GetTopicMetaOrigin(), "topics", &t)
				if err == nil {
					entity.SetEmbeddedAsset("TopicMeta", &t)
				}
				break
			default:
				log.Info(fmt.Sprintf("Embedded element (%v) is not supported: ", ele))
			}
		}
	}
	return
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

	err := m.GetDocuments(query, 0, 0, "_id", "postcategories", &cats)

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

	err := m.GetDocuments(query, 0, 0, "_id", "tags", &tags)

	if err != nil {
		return nil, m.NewStorageError(err, "GetCategories", "storage.posts.get_tags")
	}

	return tags, nil
}

func (m *MongoStorage) GetVideo(id bson.ObjectId) (models.Video, error) {
	var mgoVideo models.MongoVideo

	if id == "" {
		return models.Video{}, models.NewAppError("GetVideo", "storage.posts.get_video.id_not_provided", "Resource not found", http.StatusNotFound)
	}

	err := m.db.DB("plate").C("videos").FindId(id).One(&mgoVideo)
	if err != nil {
		return models.Video{}, m.NewStorageError(err, "GetVideo", "storage.posts.get_video.error")
	}

	return mgoVideo.ToVideo(), nil
}

// GetImage ...
func (m *MongoStorage) GetImage(id bson.ObjectId) (models.Image, error) {
	var mgoImg models.MongoImage

	err := m.GetDocument(id, "images", &mgoImg)

	if err != nil {
		log.Info("err:", err.Error())
		return models.Image{}, err
	}

	return mgoImg.ToImage(), nil
}
