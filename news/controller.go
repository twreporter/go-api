package news

import (
	"context"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/configs"
	"twreporter.org/go-api/globals"
)

type newsStorage interface {
	GetPosts(context.Context, *Query) ([]Post, error)
	GetTopics(context.Context, *Query) ([]Topic, error)

	GetPostCount(context.Context, *Filter) (int, error)
	GetTopicCount(context.Context, *Filter) (int, error)
}

type newsController struct {
	Storage newsStorage
}

func NewController(s newsStorage) *newsController {
	return &newsController{s}
}

func (nc *newsController) GetPosts(c *gin.Context) {
	q := NewQuery(FromUrlQueryMap(c.Request.URL.Query()))

	posts, err := nc.Storage.GetPosts(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	total, err := nc.Storage.GetPostCount(c, &q.Filter)
	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": posts, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}})
}

func (nc *newsController) GetAPost(c *gin.Context) {
	q := NewQuery(FromSlug(c.Param("slug")), FromUrlQueryMap(c.Request.URL.Query()))
	q.Limit = 1

	posts, err := nc.Storage.GetPosts(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": posts[0]})
}

func (nc *newsController) GetTopics(c *gin.Context) {
	q := NewQuery(FromUrlQueryMap(c.Request.URL.Query()))

	topics, err := nc.Storage.GetTopics(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	total, err := nc.Storage.GetTopicCount(c, &q.Filter)
	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": topics, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}})
}

func (nc *newsController) GetATopic(c *gin.Context) {
	q := NewQuery(FromSlug(c.Param("slug")), FromUrlQueryMap(c.Request.URL.Query()))
	q.Limit = 1

	topics, err := nc.Storage.GetTopics(c, q)

	if err != nil {
		log.Errorf("%+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": topics[0]})
}

type Job struct {
	Name  string
	Query Query
	Type  string
}

type Result struct {
	Name    string
	Content interface{}
	Error   error
}

const (
	typePost  = "post"
	typeTopic = "topic"
)

func (nc *newsController) GetIndexPage(c *gin.Context) {
	reviewCategory, _ := primitive.ObjectIDFromHex(configs.ReviewListID)
	photoCategory, _ := primitive.ObjectIDFromHex(configs.PhotographyListID)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	jobs := []Job{
		{
			Name: globals.LatestSection,
			Type: typePost,
			Query: Query{
				Filter: Filter{
					State: "published",
				},
				Pagination: Pagination{
					Limit: 6,
				},
				Sort: Sort{
					PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)},
				},
			},
		}, {
			Name: globals.EditorPicksSection,
			Type: typePost,
			Query: Query{
				Filter: Filter{
					State:      "published",
					IsFeatured: null.BoolFrom(true),
				},
				Pagination: Pagination{
					Limit: 6,
				},
				Sort: Sort{
					UpdatedAt: SortOrder{IsAsc: null.BoolFrom(false)},
				},
			},
		}, {
			Name: globals.LatestTopicSection,
			Type: typeTopic,
			Query: Query{
				Filter: Filter{
					State: "published",
				},
				Pagination: Pagination{
					Limit: 1,
				},
				Sort: Sort{
					PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)},
				},
				Full: true,
			},
		}, {
			Name: globals.ReviewsSection,
			Type: typePost,
			Query: Query{
				Filter: Filter{
					State: "published",
					Categories: []primitive.ObjectID{
						reviewCategory,
					},
				},
				Pagination: Pagination{
					Limit: 4,
				},
				Sort: Sort{
					PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)},
				},
			},
		}, {
			Name: globals.PhotoSection,
			Type: typePost,
			Query: Query{
				Filter: Filter{
					State: "published",
					Categories: []primitive.ObjectID{
						photoCategory,
					},
				},
				Pagination: Pagination{
					Limit: 6,
				},
				Sort: Sort{
					PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)},
				},
			},
		}, {
			Name: globals.InfographicSection,
			Type: typePost,
			Query: Query{
				Filter: Filter{
					State: "published",
					Style: "interactive",
				},
				Pagination: Pagination{
					Limit: 6,
				},
				Sort: Sort{
					PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)},
				},
			},
		}, {
			Name: globals.TopicsSection,
			Type: typeTopic,
			Query: Query{
				Filter: Filter{
					State: "published",
				},
				Pagination: Pagination{
					Limit:  4,
					Offset: 1,
				},
				Sort: Sort{
					PublishedDate: SortOrder{IsAsc: null.BoolFrom(false)},
				},
			},
		},
	}

	jobStream := nc.prepareJobs(ctx, jobs)

	workers := make([]<-chan Result, runtime.NumCPU())
	for i, _ := range workers {
		workers[i] = nc.fetchJobs(ctx, jobStream)
	}

	const totalSection = 7
	var err error
	results := make(map[string]interface{})

	defer func() {
		if err != nil {
			log.Errorf("%+v", err)
		}
	}()
	for result := range nc.merge(ctx, workers...) {
		select {
		case <-ctx.Done():
			err = errors.New("timeout cancel")
			break
		default:
			if result.Error != nil {
				err = result.Error
				break
			}
			results[result.Name] = result.Content
			if len(results) == totalSection {
				break
			}
		}
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "records": results})
}

func (nc *newsController) prepareJobs(ctx context.Context, jobs []Job) <-chan Job {
	jobStream := make(chan Job, len(jobs))

	go func() {
		defer close(jobStream)
		for _, job := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobStream <- job:
			}
		}
	}()
	return jobStream
}

func (nc *newsController) fetchJobs(ctx context.Context, jobStream <-chan Job) <-chan Result {
	resultStream := make(chan Result)

	go func() {
		defer close(resultStream)
		for job := range jobStream {
			select {
			case <-ctx.Done():
				return
			default:
				switch job.Type {
				case typePost:
					posts, err := nc.Storage.GetPosts(ctx, &job.Query)
					result := Result{
						Name:    job.Name,
						Content: posts,
						Error:   err,
					}
					resultStream <- result
				case typeTopic:
					topics, err := nc.Storage.GetTopics(ctx, &job.Query)
					result := Result{
						Name:    job.Name,
						Content: topics,
						Error:   err,
					}
					resultStream <- result
				}
			}
		}
	}()
	return resultStream
}

func (nc *newsController) merge(ctx context.Context, resultStream ...<-chan Result) <-chan Result {
	var wg sync.WaitGroup

	wg.Add(len(resultStream))
	fanIn := make(chan Result)

	multiplex := func(result <-chan Result) {
		defer wg.Done()
		for r := range result {
			select {
			case <-ctx.Done():
				return
			case fanIn <- r:
			}
		}
	}

	for _, r := range resultStream {
		go multiplex(r)
	}

	go func() {
		defer close(fanIn)
		wg.Wait()
	}()
	return fanIn
}
