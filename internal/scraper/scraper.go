package scraper

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/logger"
)

func StartScraping(db *database.Queries, concurrency int, interval time.Duration) {
	logger.Infof("Starting RSS scraping with interval %v", interval)

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		logger.Info("Ticker triggered: Fetching feeds...")

		// Get feeds ordered by priority (high priority first, oldest updated first)
		feeds, err := db.GetFeedsByPriority(context.Background())
		if err != nil {
			logger.ErrorErr(err, "Error fetching feeds")
			continue
		}

		logger.Infof("Found %d feeds to fetch (prioritized)", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
		logger.Debug("All feeds fetched successfully for this cycle")
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	logger.Debugf("Scraping feed: %s", feed.Name)

	parsedFeed, errorParsedFeed := fetchFeed(feed.Url)
	if errorParsedFeed != nil {
		logger.ErrorErr(errorParsedFeed, "Failed to fetch feed")
		return
	}

	for _, item := range parsedFeed.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt := time.Time{}
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		} else {
			publishedAt = time.Now()
		}

		_, errCreatePost := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if errCreatePost != nil {
			if pqErr, ok := errCreatePost.(*pq.Error); ok && pqErr.Code == "23505" {
				logger.Debugf("Post already exists, skipping: %s", item.Title)
				continue
			}
			logger.ErrorErr(errCreatePost, "Failed to create post")
		} else {
			logger.Debugf("Successfully created post: %s", item.Title)
		}
	}
}
