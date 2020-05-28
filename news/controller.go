package news

import "github.com/gin-gonic/gin"

type newsStorage interface {
	GetMetaOfPosts() ([]Post, error)
	GetFullPosts() ([]Post, error)
	GetMetaOfTopics() ([]Topic, error)
	GetFullOfTopics() ([]Topic, error)

	Count() (int, error)
}

type newsController struct {
	Storage newsStorage
}

func NewController(s newsStorage) *newsController {
	return &newsController{s}
}

func (nc *newsController) GetPosts(c *gin.Context) {
}

func (nc *newsController) GetAPost(c *gin.Context) {
}

func (nc *newsController) GetTopics(c *gin.Context) {
}

func (nc *newsController) GetATopic(c *gin.Context) {
}
