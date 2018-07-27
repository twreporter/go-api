package storage

import (
	"fmt"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// _StringToPscalCase - change WORD to pscal-case
// EX: leading_image_portrait -> LeadingImagePortrait
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

// GetEmbeddedAsset - get full embedded assets of the instance of models.NewsEntity
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
						var imgs []models.MongoImage
						err := m._GetAssetsByIDs(ids, "images", &imgs)
						if err == nil && len(imgs) > 0 {
							img := imgs[0].ToImage()
							entity.SetEmbeddedAsset(assetName, &img)
						}
					}
				}
				break
			case "leading_video":
				if ids := entity.GetEmbeddedAsset("LeadingVideoOrigin"); ids != nil {
					if len(ids) > 0 {
						var videos []models.MongoVideo
						err := m._GetAssetsByIDs(ids, "videos", &videos)
						if err == nil && len(videos) > 0 {
							video := videos[0].ToVideo()
							entity.SetEmbeddedAsset("LeadingVideo", &video)
						}
					}
				}
				break
			case "theme":
				if ids := entity.GetEmbeddedAsset("ThemeOrigin"); ids != nil {
					if len(ids) > 0 {
						var themes []models.Theme
						err := m._GetAssetsByIDs(ids, "themes", &themes)
						if err == nil && len(themes) > 0 {
							entity.SetEmbeddedAsset("Theme", &themes[0])
						}
					}
				}
				break
			case "categories":
				if ids := entity.GetEmbeddedAsset("CategoriesOrigin"); ids != nil {
					if len(ids) > 0 {
						var categories []models.Category
						err := m._GetAssetsByIDs(ids, "postcategories", &categories)
						if err == nil {
							entity.SetEmbeddedAsset("Categories", categories)
						}
					}
				}
				break
			case "tags":
				if ids := entity.GetEmbeddedAsset("TagsOrigin"); ids != nil {
					if len(ids) > 0 {
						var tags []models.Tag
						err := m._GetAssetsByIDs(ids, "tags", &tags)
						if err == nil {
							entity.SetEmbeddedAsset("Tags", tags)
						}
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

// _GetAssetsByIDs - get whole assets by their objectIDs
func (m *MongoStorage) _GetAssetsByIDs(ids []bson.ObjectId, collectionName string, v interface{}) error {
	if ids == nil {
		return nil
	}

	query := models.MongoQuery{
		IDs: models.MongoQueryComparison{
			In: ids,
		},
	}

	_, err := m.GetDocuments(query, 0, 0, "_id", collectionName, v)

	if err != nil {
		return m.NewStorageError(err, "MongoStorage._GetAssetsByIDs", fmt.Sprintf("get assets %s by ids %v occurs error", collectionName, ids))
	}

	return nil
}

// GetAuthors - get whole author documents by their objectIDs
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
