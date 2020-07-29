package controllers

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"twreporter.org/go-api/globals"
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

	ctx, cancel := context.WithTimeout(c, globals.Conf.News.PostPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParsePostListQuery(c)

	posts, err := nc.Storage.GetMetaOfPosts(ctx, q)

	if err != nil {
		return
	}

	total, err := nc.Storage.GetPostCount(ctx, q)
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

	ctx, cancel := context.WithTimeout(c, globals.Conf.News.PostPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseSinglePostQuery(c)

	if q.Full {
		var posts []news.Post
		posts, err = nc.Storage.GetFullPosts(ctx, q)
		if len(posts) > 0 {
			post = posts[0]
		}
	} else {
		var posts []news.MetaOfPost
		posts, err = nc.Storage.GetMetaOfPosts(ctx, q)
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

	ctx, cancel := context.WithTimeout(c, globals.Conf.News.TopicPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseTopicListQuery(c)

	topics, err := nc.Storage.GetMetaOfTopics(ctx, q)

	if err != nil {
		return
	}

	total, err := nc.Storage.GetTopicCount(ctx, q)
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

	ctx, cancel := context.WithTimeout(c, globals.Conf.News.TopicPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseSingleTopicQuery(c)

	if q.Full {
		var topics []news.Topic
		topics, err = nc.Storage.GetFullTopics(ctx, q)
		if len(topics) > 0 {
			topic = topics[0]
		}
	} else {
		var topics []news.MetaOfTopic
		topics, err = nc.Storage.GetMetaOfTopics(ctx, q)
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
	ctx, cancel := context.WithTimeout(c, globals.Conf.News.IndexPageTimeout)
	defer cancel()

	jobs := nc.getIndexPageJobs()

	var err error
	results := make(map[string]interface{})

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	for result := range nc.fetchjobs(ctx, jobs) {
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
	// v1 index page section jobs
	jobs := []job{
		{
			Name:  news.LatestSection,
			Type:  typePost,
			Query: news.NewQuery(news.WithLimit(6)),
		}, {
			Name: news.EditorPicksSection,
			Type: typePost,
			Query: news.NewQuery(
				news.WithFilterIsFeatured(true),
				news.WithLimit(6),
				news.WithSortUpdatedAt(false)),
		}, {
			Name:  news.LatestTopicSection,
			Type:  typeTopic,
			Query: news.NewQuery(news.WithLimit(1)),
		}, {
			Name: news.ReviewsSection,
			Type: typePost,
			Query: news.NewQuery(
				news.WithFilterCategoryIDs(news.Review.ID),
				news.WithLimit(4)),
		}, {
			Name: news.PhotoSection,
			Type: typePost,
			Query: news.NewQuery(
				news.WithFilterCategoryIDs(news.Photography.ID),
				news.WithLimit(6)),
		}, {
			Name:  news.InfographicSection,
			Type:  typePost,
			Query: news.NewQuery(news.WithFilterStyle("interactive"), news.WithLimit(6)),
		}, {
			Name:  news.TopicsSection,
			Type:  typeTopic,
			Query: news.NewQuery(news.WithOffset(1), news.WithLimit(4)),
		},
	}

	// v1 categories in index page
	for _, v := range []news.Category{
		news.HumanRightsAndSociety,
		news.EnvironmentAndEducation,
		news.PoliticsAndEconomy,
		news.CultureAndArt,
		news.International,
		news.LivingAndMedicalCare,
	} {
		jobs = append(jobs, job{
			v.Name,
			typePost,
			news.NewQuery(news.WithFilterCategoryIDs(v.ID), news.WithLimit(1)),
		})
	}

	return jobs
}

func (nc *newsV2Controller) fetchjobs(ctx context.Context, jobs []job) <-chan result {
	resultStream := make(chan result)

	var wg sync.WaitGroup

	wg.Add(len(jobs))
	for _, j := range jobs {
		go func(job job) {
			defer wg.Done()
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
		}(j)
	}

	go func() {
		defer close(resultStream)
		wg.Wait()
	}()
	return resultStream
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
