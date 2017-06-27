package storage

import (
	"twreporter.org/go-api/models"
)

// GetMetaOfPosts is a type-specific functions implementing the method defined in the NewsStorage.
// It parses query string into bson and finds the posts according to that bson.
func (m *MongoStorage) GetMetaOfPosts(qs interface{}, limit int, offset int, sort string, embedded []string) ([]models.Post, error) {
	var posts []models.Post

	if embedded == nil {
		embedded = []string{"hero_image", "categories", "tags", "topic_meta", "og_image"}
	}

	err := m.GetDocuments(qs, limit, offset, sort, "posts", &posts)

	if err != nil {
		return posts, err
	}

	for index := range posts {
		m.GetEmbeddedAsset(&posts[index], embedded)
	}

	return posts, nil
}
