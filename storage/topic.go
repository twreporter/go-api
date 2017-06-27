package storage

import (
	"twreporter.org/go-api/models"
	//log "github.com/Sirupsen/logrus"
)

// GetTopics is a type-specific functions implementing the method defined in the NewsStorage.
// It parses query string into bson and finds the topics with embedded assets according to that bson.
func (m *MongoStorage) GetTopics(qs interface{}, limit int, offset int, sort string, embedded []string) ([]models.Topic, int, error) {
	var topics []models.Topic
	if embedded == nil {
		embedded = []string{"relateds_meta", "leading_image", "leading_image_portrait", "leading_video", "og_image"}
	}

	total, err := m.GetDocuments(qs, limit, offset, sort, "topics", &topics)

	if err != nil {
		return topics, 0, err
	}

	for index := range topics {
		m.GetEmbeddedAsset(&topics[index], embedded)
	}

	return topics, total, nil
}
