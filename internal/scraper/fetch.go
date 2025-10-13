package scraper

import (
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

func fetchFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()

	fp.Client = &http.Client{Timeout: time.Second * 10}

	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, nil
}
