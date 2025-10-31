package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
	"github.com/mehmettalhairmak/rss-aggregator/internal/logger"
	"github.com/mehmettalhairmak/rss-aggregator/internal/realtime"
	"github.com/rs/zerolog"
)

type Scraper struct {
	DB     *database.Queries
	Logger zerolog.Logger
	Hub    *realtime.Hub
}

func NewScraper(db *database.Queries, log zerolog.Logger, hub *realtime.Hub) *Scraper {
	return &Scraper{
		DB:     db,
		Logger: log,
		Hub:    hub,
	}
}

func (s *Scraper) StartScraping(db *database.Queries, interval time.Duration) {
	s.Logger.Info().Msgf("Starting RSS scraping with interval %v", interval)

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for range ticker.C {
		s.Logger.Info().Msg("Ticker triggered: Fetching feeds...")

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
			go s.scrapeFeed(db, wg, feed)
		}
		wg.Wait()
		s.Logger.Debug().Msg("All feeds fetched successfully for this cycle")
	}
}

func (s *Scraper) scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	logger.Debugf("Scraping feed: %s", feed.Name)

	parsedFeed, errorParsedFeed := fetchFeed(feed.Url)
	if errorParsedFeed != nil {
		s.Logger.Error().Err(errorParsedFeed).Msg("Failed to fetch feed")
		return
	}

	newPostCount := 0

	for _, item := range parsedFeed.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		var publishedAt time.Time
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
				s.Logger.Debug().Msgf("Post already exists, skipping: %s", item.Title)
				continue
			}
			s.Logger.Error().Err(errCreatePost).Msg("Failed to create post")
		} else {
			newPostCount++
			s.Logger.Debug().Msgf("Successfully created post: %s", item.Title)
		}
	}

	if newPostCount > 0 {
		s.sendNewPostSignal(context.Background(), feed, newPostCount)
	}
}

func (s *Scraper) sendNewPostSignal(ctx context.Context, feed database.Feed, newCount int) {
	followers, err := s.DB.GetFollowersByFeedID(ctx, feed.ID)
	if err != nil {
		s.Logger.Error().Err(err).Msgf("Scraper failed to get followers for feed %s", feed.ID)
		return
	}

	signals := make(map[uuid.UUID][]byte)
	signalPayload := []byte(fmt.Sprintf(
		`{"type": "NEW_POST_AVAILABLE", "feed_id": "%s", "feed_name": "%s", "count": %d}`,
		feed.ID.String(), feed.Name, newCount,
	))

	for _, follower := range followers {
		signals[follower] = signalPayload
	}

	if len(signals) > 0 {
		s.Hub.SendSignal(signals)
		s.Logger.Info().
			Int("followers_count", len(signals)).
			Str("feed_id", feed.ID.String()).
			Msg("New post signal published to Hub.")
	}
}
