package storage

import (
	"twreporter.org/go-api/models"
)

// _GetPosts finds the posts according to query string and also get the embedded assets
func (m *MongoStorage) _GetPosts(qs interface{}, limit int, offset int, sort string, embedded []string, isFull bool) ([]models.Post, int, error) {
	var posts []models.Post
	total, err := m.GetDocuments(qs, limit, offset, sort, "posts", &posts)

	if err != nil {
		return posts, 0, err
	}

	for index := range posts {
		m.GetEmbeddedAsset(&posts[index], embedded)
		if isFull == false {
			posts[index].Content = nil
		}

		posts[index].Full = isFull
	}

	return posts, total, nil
}

// GetMetaOfPosts is a type-specific functions implementing the method defined in the NewsStorage.
// It finds the posts according to query string and only return the metadata of posts.
func (m *MongoStorage) GetMetaOfPosts(qs interface{}, limit int, offset int, sort string, embedded []string) ([]models.Post, int, error) {
	if embedded == nil {
		embedded = []string{"hero_image", "categories", "tags", "topic", "og_image"}
	}

	return m._GetPosts(qs, limit, offset, sort, embedded, false)
}

// GetFullPosts is a type-specific functions implementing the method defined in the NewsStorage.
// It finds the posts according to query string.
func (m *MongoStorage) GetFullPosts(qs interface{}, limit int, offset int, sort string, embedded []string) ([]models.Post, int, error) {
	if embedded == nil {
		embedded = []string{"hero_image", "leading_video", "categories", "tags", "topic", "og_image", "writters", "photographers", "designers", "engineers", "relateds"}
	}

	return m._GetPosts(qs, limit, offset, sort, embedded, true)
}
