package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
)

func StartScraping(db *database.Queries, concurrency int, interval time.Duration) {
	log.Printf("Starting scraping with interval %v", interval)

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		log.Println("Ticker triggered: Fetching feeds...")

		feeds, err := db.GetFeeds(context.Background())
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}

		fmt.Printf("Found %d feeds to fetch \n", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
		log.Println("All feeds fetched successfully for this cycle.")
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	log.Printf("Scraping feed %s", feed.Name)

	parsedFeed, errorParsedFeed := fetchFeed(feed.Url)
	if errorParsedFeed != nil {
		log.Printf("Error fetching feed '%s': %v", feed.Name, errorParsedFeed)
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
				log.Printf("Post already exists, skipping: %s", item.Title)
				continue
			}
			log.Printf("Error creating post: %v", errCreatePost)
		} else {
			log.Printf("Successfully created post: %s", item.Title)
		}
	}
}
