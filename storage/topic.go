package storage

import (
	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

// _GetTopics finds the topics according to query string and also get the embedded assets
func (m *MongoStorage) _GetTopics(mq models.MongoQuery, limit int, offset int, sort string, embedded []string, isFull bool) ([]models.Topic, int, error) {
	var topics []models.Topic

	if globals.Conf.Environment != "development" {
		mq.State = "published"
	}

	total, err := m.GetDocuments(mq, limit, offset, sort, "topics", &topics)

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
func (m *MongoStorage) GetFullTopics(mq models.MongoQuery, limit int, offset int, sort string, embedded []string) ([]models.Topic, int, error) {
	if embedded == nil {
		embedded = []string{"topic_relateds", "leading_image", "leading_image_portrait", "leading_video", "og_image"}
	}

	return m._GetTopics(mq, limit, offset, sort, embedded, true)
}

// GetMetaOfTopics is a type-specific functions implementing the method defined in the NewsStorage.
// It will get full topics having PARTIAL corresponding assets
func (m *MongoStorage) GetMetaOfTopics(mq models.MongoQuery, limit int, offset int, sort string, embedded []string) ([]models.Topic, int, error) {
	if embedded == nil {
		embedded = []string{"leading_image", "leading_image_portrait", "og_image"}
	}

	return m._GetTopics(mq, limit, offset, sort, embedded, false)
}
