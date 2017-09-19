package storage

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

func (m *MongoStorage) _StringToPscalCase(str string) (pscalCase string) {
	isToUpper := true
	for _, runeValue := range str {
		if isToUpper {
			pscalCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '_' {
				isToUpper = true
			} else {
				pscalCase += string(runeValue)
			}
		}
	}
	return
}

// GetEmbeddedAsset ...
func (m *MongoStorage) GetEmbeddedAsset(entity models.NewsEntity, embedded []string) {
	if embedded != nil {
		for _, ele := range embedded {
			switch ele {
			case "writters", "photographers", "designers", "engineers":
				assetName := m._StringToPscalCase(ele)
				if ids := entity.GetEmbeddedAsset(assetName + "Origin"); ids != nil {
					if len(ids) > 0 {
						authors, err := m.GetAuthors(ids)
						if err == nil {
							entity.SetEmbeddedAsset(assetName, authors)
						}
					}
				}
				break
			case "hero_image", "leading_image", "leading_image_portrait", "og_image":
				assetName := m._StringToPscalCase(ele)
				if ids := entity.GetEmbeddedAsset(assetName + "Origin"); ids != nil {
					if len(ids) > 0 {
						img, err := m.GetImage(ids[0])
						if err == nil {
							entity.SetEmbeddedAsset(assetName, &img)
						}
					}
				}
				break
			case "leading_video":
				if ids := entity.GetEmbeddedAsset("LeadingVideoOrigin"); ids != nil {
					if len(ids) > 0 {
						video, err := m.GetVideo(ids[0])
						if err == nil {
							entity.SetEmbeddedAsset("LeadingVideo", &video)
						}
					}
				}
				break
			case "categories":
				if ids := entity.GetEmbeddedAsset("CategoriesOrigin"); ids != nil {
					if len(ids) > 0 {
						categories, _ := m.GetCategories(ids)
						_categories := make([]models.Category, len(categories))
						for i, v := range categories {
							_categories[i] = v
						}
						entity.SetEmbeddedAsset("Categories", _categories)
					}
				}
				break
			case "tags":
				if ids := entity.GetEmbeddedAsset("TagsOrigin"); ids != nil {
					if len(ids) > 0 {
						tags, _ := m.GetTags(ids)
						_tags := make([]models.Tag, len(tags))
						for i, v := range tags {
							_tags[i] = v
						}
						entity.SetEmbeddedAsset("Tags", _tags)
					}
				}
				break
			case "relateds", "topic_relateds":
				var _relateds []models.Post
				if ids := entity.GetEmbeddedAsset("RelatedsOrigin"); ids != nil {
					if len(ids) > 0 {
						query := models.MongoQuery{
							IDs: models.MongoQueryComparison{
								In: ids,
							},
						}
						var embedded []string
						if ele == "topic_relateds" {
							embedded = []string{"hero_image", "categories", "tags", "og_image"}
						}
						relateds, _, err := m.GetMetaOfPosts(query, 0, 0, "-publishedDate", embedded)
						for _, id := range ids {
							for _, related := range relateds {
								if id.Hex() == related.ID.Hex() {
									_relateds = append(_relateds, related)
								}
							}
						}

						if err == nil {
							entity.SetEmbeddedAsset("Relateds", _relateds)
						}
					}
				}
				break
			case "topic":
				if ids := entity.GetEmbeddedAsset("TopicOrigin"); ids != nil {
					if len(ids) > 0 {
						query := models.MongoQuery{
							IDs: models.MongoQueryComparison{
								In: ids,
							},
						}

						topics, _, err := m.GetMetaOfTopics(query, 0, 0, "-publishedDate", nil)

						if err == nil && len(topics) > 0 {
							entity.SetEmbeddedAsset("Topic", &topics[0])
						}
					}
				}
				break
			case "topic_full":
				if ids := entity.GetEmbeddedAsset("TopicOrigin"); ids != nil {
					if len(ids) > 0 {
						query := models.MongoQuery{
							IDs: models.MongoQueryComparison{
								In: ids,
							},
						}

						topics, _, err := m.GetFullTopics(query, 0, 0, "-publishedDate", nil)

						if err == nil && len(topics) > 0 {
							entity.SetEmbeddedAsset("Topic", &topics[0])
						}
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

// GetCategories ...
func (m *MongoStorage) GetCategories(ids []bson.ObjectId) ([]models.Category, error) {
	var cats []models.Category

	if ids == nil {
		return cats, nil
	}

	query := models.MongoQuery{
		IDs: models.MongoQueryComparison{
			In: ids,
		},
	}

	_, err := m.GetDocuments(query, 0, 0, "_id", "postcategories", &cats)

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

	query := models.MongoQuery{
		IDs: models.MongoQueryComparison{
			In: ids,
		},
	}

	_, err := m.GetDocuments(query, 0, 0, "_id", "tags", &tags)

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

func (m *MongoStorage) GetAuthors(ids []bson.ObjectId) ([]models.Author, error) {
	var authors []models.Author
	var orderedAuthors []models.Author

	if ids == nil {
		return authors, nil
	}

	query := models.MongoQuery{
		IDs: models.MongoQueryComparison{
			In: ids,
		},
	}

	_, err := m.GetDocuments(query, 0, 0, "_id", "contacts", &authors)

	if err != nil {
		return authors, err
	}

	for _, id := range ids {
		for _, author := range authors {
			if author.ID == id {
				orderedAuthors = append(orderedAuthors, author)
			}
		}
	}

	return orderedAuthors, nil
}
