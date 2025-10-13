package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/mehmettalhairmak/rss-aggregator/internal/database"
)

// User represents a user in the API
// Different from database model - used for API responses
// Not password_hash is NOT included here for security reasons
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}

// Feed represents an RSS feed in the API
type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

// FeedFollow represents a feed follow relationship in the API
// Represents a user's feed subscription
type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

// DatabaseUserToUser converts a database user to an API user
// Note: password_hash is NOT included for security
func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Email:     dbUser.Email.String, // sql.NullString -> string conversion
	}
}

// DatabaseFeedToFeed converts a database feed to an API feed
func DatabaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

func DatabasePostToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		PublishedAt: dbPost.PublishedAt,
		FeedID:      dbPost.FeedID,
	}
}

func DatabaseAllPostToAllPost(dbPosts []database.Post) []Post {
	posts := make([]Post, 0, len(dbPosts))
	for _, post := range dbPosts {
		posts = append(posts, DatabasePostToPost(post))
	}
	return posts
}

// DatabaseAllFeedToAllFeed converts multiple database feeds to API feeds
func DatabaseAllFeedToAllFeed(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, 0, len(dbFeeds)) // Initialize with capacity for performance
	for _, feed := range dbFeeds {
		feeds = append(feeds, DatabaseFeedToFeed(feed))
	}
	return feeds
}

// DatabaseFeedFollowToFeedFollow converts a database feed follow to an API feed follow
func DatabaseFeedFollowToFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.FeedID,
	}
}

// DatabaseAllFeedFollowToAllFeedFollow converts multiple database feed follows to API feed follows
func DatabaseAllFeedFollowToAllFeedFollow(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, 0, len(dbFeedFollows))
	for _, feedFollow := range dbFeedFollows {
		feedFollows = append(feedFollows, DatabaseFeedFollowToFeedFollow(feedFollow))
	}
	return feedFollows
}
