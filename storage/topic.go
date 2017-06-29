package storage

import (
	"twreporter.org/go-api/models"
	//log "github.com/Sirupsen/logrus"
)

// _GetTopics finds the topics according to query string and also get the embedded assets
func (m *MongoStorage) _GetTopics(qs interface{}, limit int, offset int, sort string, embedded []string, isFull bool) ([]models.Topic, int, error) {
	var topics []models.Topic
	total, err := m.GetDocuments(qs, limit, offset, sort, "topics", &topics)

	if err != nil {
		return topics, 0, err
	}

	for index := range topics {
		m.GetEmbeddedAsset(&topics[index], embedded)
		if isFull {
			topics[index].Full = isFull
		}
	}

	return topics, total, nil
}

// GetFullTopics is a type-specific functions implementing the method defined in the NewsStorage.
// It will get full topics having ALL the corresponding assets
func (m *MongoStorage) GetFullTopics(qs interface{}, limit int, offset int, sort string, embedded []string) ([]models.Topic, int, error) {
	if embedded == nil {
		embedded = []string{"relateds", "leading_image", "leading_image_portrait", "leading_video", "og_image"}
	}

	return m._GetTopics(qs, limit, offset, sort, embedded, true)
}

// GetMetaOfTopics is a type-specific functions implementing the method defined in the NewsStorage.
// It will get full topics having PARTIAL corresponding assets
func (m *MongoStorage) GetMetaOfTopics(qs interface{}, limit int, offset int, sort string, embedded []string) ([]models.Topic, int, error) {
	if embedded == nil {
		embedded = []string{"leading_image", "leading_image_portrait", "og_image"}
	}

	return m._GetTopics(qs, limit, offset, sort, embedded, false)
}
