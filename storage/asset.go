package storage

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// GetEmbeddedAsset ...
func (m *MongoStorage) GetEmbeddedAsset(entity models.NewsEntity, embedded []string) {
	if embedded != nil {
		for _, ele := range embedded {
			switch ele {
			case "hero_image":
				if id := entity.GetHeroImageOrigin(); id != "" {
					img, err := m.GetImage(id)
					if err == nil {
						entity.SetEmbeddedAsset("HeroImage", &img)
					}
				}
				break
			case "leading_image":
				if id := entity.GetLeadingImageOrigin(); id != "" {
					img, err := m.GetImage(id)
					if err == nil {
						entity.SetEmbeddedAsset("LeadingImage", &img)
					}
				}
				break
			case "leading_image_portrait":
				if id := entity.GetLeadingImagePortraitOrigin(); id != "" {
					img, err := m.GetImage(id)
					if err == nil {
						entity.SetEmbeddedAsset("LeadingImagePortrait", &img)
					}
				}
				break
			case "leading_video":
				if id := entity.GetLeadingVideoOrigin(); id != "" {
					video, err := m.GetVideo(id)
					if err == nil {
						entity.SetEmbeddedAsset("LeadingVideo", &video)
					}
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
			case "relateds_meta":
				ids := entity.GetRelatedsOrigin()
				relateds, err := m.GetRelatedsMeta(ids)
				if err == nil {
					entity.SetEmbeddedAsset("Relateds", relateds)
				}
				break
			case "topic_meta":
				if id := entity.GetTopicOrigin(); id != "" {
					t, err := m.GetTopicMeta(id)
					if err == nil {
						entity.SetEmbeddedAsset("Topic", &t)
					}
				}
				break
			default:
				log.Info(fmt.Sprintf("Embedded element (%v) is not supported: ", ele))
			}
		}
	}
	return
}

func (m *MongoStorage) GetTopicMeta(id bson.ObjectId) (models.Topic, error) {

	query := bson.M{
		"_id": id,
	}

	topics, err := m.GetTopics(query, 0, 0, "-publishedDate", []string{"leading_image", "og_image"})

	if err != nil {
		return models.Topic{}, err
	}

	return topics[0], nil
}

// GetRelatedsMeta ...
func (m *MongoStorage) GetRelatedsMeta(ids []bson.ObjectId) ([]models.PostMeta, error) {

	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	posts, err := m.GetMetaOfPosts(query, 0, 0, "-publishedDate", []string{"hero_image", "og_image"})

	if err != nil {
		return nil, err
	}

	return posts, nil
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

// GetVideo ...
func (m *MongoStorage) GetVideo(id bson.ObjectId) (models.Video, error) {
	var mgoVideo models.MongoVideo

	err := m.GetDocument(id, "videos", &mgoVideo)

	if err != nil {
		return models.Video{}, err
	}

	return mgoVideo.ToVideo(), nil
}

// GetImage ...
func (m *MongoStorage) GetImage(id bson.ObjectId) (models.Image, error) {
	var mgoImg models.MongoImage

	err := m.GetDocument(id, "images", &mgoImg)

	if err != nil {
		return models.Image{}, err
	}

	return mgoImg.ToImage(), nil
}
