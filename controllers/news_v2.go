package controllers

import (
	"context"
	"net/http"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"twreporter.org/go-api/internal/news"
)

type newsV2Storage interface {
	GetFullPosts(context.Context, *news.Query) ([]news.Post, error)
	GetMetaOfPosts(context.Context, *news.Query) ([]news.MetaOfPost, error)
	GetFullTopics(context.Context, *news.Query) ([]news.Topic, error)
	GetMetaOfTopics(context.Context, *news.Query) ([]news.MetaOfTopic, error)

	GetPostCount(context.Context, *news.Query) (int, error)
	GetTopicCount(context.Context, *news.Query) (int, error)
}

func NewNewsV2Controller(s newsV2Storage) *newsV2Controller {
	return &newsV2Controller{s}
}

type newsV2Controller struct {
	Storage newsV2Storage
}

func (nc *newsV2Controller) GetPosts(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParsePostListQuery(c)

	// TODO(babygoat): config context with proper timeout
	posts, err := nc.Storage.GetMetaOfPosts(c, q)

	if err != nil {
		return
	}

	// TODO(babygoat): config context with proper timeout
	total, err := nc.Storage.GetPostCount(c, q)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": posts, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}}})
}

func (nc *newsV2Controller) GetAPost(c *gin.Context) {
	var post interface{}
	var err error

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseSinglePostQuery(c)

	if q.Full {
		var posts []news.Post
		// TODO(babygoat): config context with proper timeout
		posts, err = nc.Storage.GetFullPosts(c, q)
		if len(posts) > 0 {
			post = posts[0]
		}
	} else {
		var posts []news.MetaOfPost
		// TODO(babygoat): config context with proper timeout
		posts, err = nc.Storage.GetMetaOfPosts(c, q)
		if len(posts) > 0 {
			post = posts[0]
		}
	}

	if err != nil {
		return
	}

	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"slug": "Cannot find the post from the slug"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

func (nc *newsV2Controller) GetTopics(c *gin.Context) {
	var err error

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseTopicListQuery(c)

	// TODO(babygoat): config context with proper timeout
	topics, err := nc.Storage.GetMetaOfTopics(c, q)

	if err != nil {
		return
	}

	// TODO(babygoat): config context with proper timeout
	total, err := nc.Storage.GetTopicCount(c, q)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": topics, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}}})
}

func (nc *newsV2Controller) GetATopic(c *gin.Context) {
	var topic interface{}
	var err error

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseSingleTopicQuery(c)

	if q.Full {
		var topics []news.Topic
		// TODO(babygoat): config context with proper timeout
		topics, err = nc.Storage.GetFullTopics(c, q)
		if len(topics) > 0 {
			topic = topics[0]
		}
	} else {
		var topics []news.MetaOfTopic
		// TODO(babygoat): config context with proper timeout
		topics, err = nc.Storage.GetMetaOfTopics(c, q)
		if len(topics) > 0 {
			topic = topics[0]
		}
	}

	// server side error
	if err != nil {
		return
	}

	if topic == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"slug": "Cannot find the topic from the slug"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": topic})
}

type job struct {
	Name  string
	Type  string
	Query *news.Query
}

type result struct {
	Name    string
	Content interface{}
	Error   error
}

const (
	typePost  = "post"
	typeTopic = "topic"
)

func (nc *newsV2Controller) GetIndexPage(c *gin.Context) {
	// TODO(babygoat): config context with proper timeout
	ctx := c

	jobs := nc.getIndexPageJobs()
	jobStream := nc.preparejobs(ctx, jobs)

	workers := make([]<-chan result, runtime.NumCPU())
	for i, _ := range workers {
		workers[i] = nc.fetchjobs(ctx, jobStream)
	}

	var err error
	results := make(map[string]interface{})

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()
	for result := range nc.merge(ctx, workers...) {
		select {
		case <-ctx.Done():
			err = errors.WithStack(ctx.Err())
			break
		default:
			if result.Error != nil {
				err = result.Error
				break
			}
			results[result.Name] = result.Content
			if len(results) == len(jobs) {
				break
			}
		}
	}
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": results}})

}

func (nc *newsV2Controller) getIndexPageJobs() []job {
	const statePublished = "published"

	// v1 index page section jobs
	jobs := []job{
		{
			Name:  news.LatestSection,
			Type:  typePost,
			Query: news.NewQuery(news.FilterState(statePublished), news.Limit(6)),
		}, {
			Name: news.EditorPicksSection,
			Type: typePost,
			Query: news.NewQuery(
				news.FilterState(statePublished),
				news.FilterIsFeatured(true),
				news.Limit(6),
				news.SortUpdatedAt(false)),
		}, {
			Name:  news.LatestTopicSection,
			Type:  typeTopic,
			Query: news.NewQuery(news.FilterState(statePublished), news.Limit(1)),
		}, {
			Name: news.ReviewsSection,
			Type: typePost,
			Query: news.NewQuery(
				news.FilterState(statePublished),
				news.FilterCategoryIDs(news.Review.ID),
				news.Limit(4)),
		}, {
			Name: news.PhotoSection,
			Type: typePost,
			Query: news.NewQuery(
				news.FilterState(statePublished),
				news.FilterCategoryIDs(news.Photography.ID),
				news.Limit(6)),
		}, {
			Name:  news.InfographicSection,
			Type:  typePost,
			Query: news.NewQuery(news.FilterState(statePublished), news.FilterStyle("interactive"), news.Limit(6)),
		}, {
			Name:  news.TopicsSection,
			Type:  typeTopic,
			Query: news.NewQuery(news.FilterState(statePublished), news.Offset(1), news.Limit(4)),
		},
	}

	// v1 categories in index page
	jobs = append(jobs, func(categories ...news.Category) []job {
		var jobs []job
		for _, v := range categories {
			jobs = append(jobs, job{
				v.Name,
				typePost,
				news.NewQuery(news.FilterState(statePublished), news.Limit(1)),
			})
		}
		return jobs
	}(
		news.HumanRightsAndSociety,
		news.EnvironmentAndEducation,
		news.PoliticsAndEconomy,
		news.CultureAndArt,
		news.International,
		news.LivingAndMedicalCare,
	)...)

	return jobs
}

func (nc *newsV2Controller) preparejobs(ctx context.Context, jobs []job) <-chan job {
	jobStream := make(chan job, len(jobs))

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

func (nc *newsV2Controller) fetchjobs(ctx context.Context, jobStream <-chan job) <-chan result {
	resultStream := make(chan result)

	go func() {
		defer close(resultStream)
		for job := range jobStream {
			select {
			case <-ctx.Done():
				return
			default:
				switch job.Type {
				case typePost:
					posts, err := nc.Storage.GetMetaOfPosts(ctx, job.Query)
					result := result{
						Name:    job.Name,
						Content: posts,
						Error:   err,
					}
					resultStream <- result
				case typeTopic:
					topics, err := nc.Storage.GetMetaOfTopics(ctx, job.Query)
					result := result{
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

func (nc *newsV2Controller) merge(ctx context.Context, resultStream ...<-chan result) <-chan result {
	var wg sync.WaitGroup

	wg.Add(len(resultStream))
	fanIn := make(chan result)

	multiplex := func(result <-chan result) {
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

func (nc *newsV2Controller) helperCleanup(c *gin.Context, err error) {
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{"status": "error", "message": "Query upstream server timeout."})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unexpected error."})
		}
		log.Errorf("%+v", err)
	}
}
