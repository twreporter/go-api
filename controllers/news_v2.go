package controllers

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	f "github.com/twreporter/logformatter"
	"github.com/twreporter/go-api/globals"
	"github.com/twreporter/go-api/internal/news"
)

type newsV2Storage interface {
	GetFullPosts(context.Context, *news.Query) ([]news.Post, error)
	GetMetaOfPosts(context.Context, *news.Query) ([]news.MetaOfPost, error)
	GetFullTopics(context.Context, *news.Query) ([]news.Topic, error)
	GetMetaOfTopics(context.Context, *news.Query) ([]news.MetaOfTopic, error)
	GetAuthors(context.Context, *news.Query) ([]news.Author, error)

	GetTags(context.Context, *news.Query) ([]news.Tag, error)

	GetPostCount(context.Context, *news.Query) (int64, error)
	GetTopicCount(context.Context, *news.Query) (int64, error)
	GetAuthorCount(context.Context, *news.Query) (int64, error)

	CheckCategorySetValid(context.Context, *news.Query) (bool, error)
}

type newsV2SqlStorage interface {
	GetBookmarksOfPosts(context.Context, string, []news.MetaOfPost) ([]news.MetaOfPost, error)
}

func NewNewsV2Controller(s newsV2Storage, client news.AlgoliaSearcher, sqls newsV2SqlStorage) *newsV2Controller {
	return &newsV2Controller{s, client, sqls}
}

type newsV2Controller struct {
	Storage     newsV2Storage
	indexClient news.AlgoliaSearcher
	SqlStorage  newsV2SqlStorage
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

	toggleBookmark, _ := strconv.Atoi(c.Query("toggleBookmark"))
	if toggleBookmark == 1 {
		c.Writer.Header().Set("Cache-Control", "no-store")
	}
	authUserID := c.Request.Context().Value(globals.AuthUserIDProperty)
	if authUserID != nil && toggleBookmark == 1 {
		authUserIdString := fmt.Sprintf("%v", authUserID)
		if _, err := nc.SqlStorage.GetBookmarksOfPosts(ctx, authUserIdString, posts); err != nil {
			log.WithField("detail", err).Errorf("%s", f.FormatStack(err))
		}
	}

	total, err := nc.Storage.GetPostCount(ctx, q)
	if err != nil {
		return
	}

	categorySetIsValid, err := nc.Storage.CheckCategorySetValid(ctx, q)
	if err != nil || !categorySetIsValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "category & subcategory is not consistent",
		})
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

func (nc *newsV2Controller) GetTags(c *gin.Context) {
	var err error

	ctx, cancel := context.WithTimeout(c, globals.Conf.News.PostPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseTagListQuery(c)

	tags, err := nc.Storage.GetTags(ctx, q)

	if err != nil {
		return
	}

	if tags == nil {
		tags = make([]news.Tag, 0)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": tags, "meta": gin.H{
		"latest_order": q.Filter.LatestOrder,
		"offset":       q.Offset,
		"limit":        q.Limit,
	}}})
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
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": results})

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
				news.WithFilterCategoryIDs(news.Opinion.Key),
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

	// v2 categories in index page
	for _, v := range []news.CategorySet{
		news.World,
		news.Humanrights,
		news.PoliticsAndSociety,
		news.Health,
		news.Econ,
		news.Culture,
		news.Education,
		news.Environment,
	} {
		jobs = append(jobs, job{
			v.Name,
			typePost,
			news.NewQuery(news.WithFilterCategorySet(v.Key), news.WithLimit(1)),
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

func (nc *newsV2Controller) GetAuthors(c *gin.Context) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), globals.Conf.News.AuthorPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseAuthorListQuery(c)

	var authors []news.Author
	var total int64
	var authorIDs []string
	authorIDs, total, err = news.GetRankedAuthorIDs(ctx, nc.indexClient, q)
	switch {
	// Return early if timeout occurs
	case errors.Is(err, context.DeadlineExceeded):
		return
	// Fallback to database query if algolia search unavailable(e.g. quota exceeds)
	// Note that empty search result will not produce any error
	case err != nil:
		if authors, err = nc.Storage.GetAuthors(ctx, q); err != nil {
			return
		}

		if total, err = nc.Storage.GetAuthorCount(ctx, q); err != nil {
			return
		}
	// Proceeds to database query with ranked author IDs to assemble the API response if result is available
	case len(authorIDs) > 0:
		queryForResponse := &news.Query{
			Filter: news.Filter{
				IDs: authorIDs,
			},
		}
		if authors, err = nc.Storage.GetAuthors(ctx, queryForResponse); err != nil {
			return
		}
		// Create lookup map for preserving the authors order as ranked result
		lookupIds := make(map[string]int)
		for index, id := range authorIDs {
			lookupIds[id] = index
		}
		// Sort the data w.r.t the map
		sort.SliceStable(authors, func(i, j int) bool {
			return lookupIds[authors[i].ID.Hex()] < lookupIds[authors[j].ID.Hex()]
		})
	}

	if len(authors) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"records": authors, "meta": gin.H{
		"total":  total,
		"offset": q.Offset,
		"limit":  q.Limit,
	}}})
}

func (nc *newsV2Controller) GetAuthorByID(c *gin.Context) {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), globals.Conf.News.AuthorPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseSingleAuthorQuery(c)

	authors, err := nc.Storage.GetAuthors(ctx, q)

	if err != nil {
		return
	}

	if len(authors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "fail", "data": gin.H{"author_id": "Cannot find the author from the author_id"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": authors[0]})
}
func (nc *newsV2Controller) GetPostsByAuthor(c *gin.Context) {
	var err error

	ctx, cancel := context.WithTimeout(c, globals.Conf.News.PostPageTimeout)
	defer cancel()

	defer func() {
		if err != nil {
			nc.helperCleanup(c, err)
		}
	}()

	q := news.ParseAuthorPostListQuery(c)

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
