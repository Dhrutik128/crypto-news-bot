package news

import (
	"context"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"time"
)

// download source and check if its a valid feed
func fetch(source string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	ctx, c := context.WithTimeout(context.Background(), 15*time.Second)
	defer c()
	feed, err := fp.ParseURLWithContext(source, ctx)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func (b *Analyzer) getFeeds() []string {
	keys := make([]string, 0, len(b.Feeds))
	for k := range b.Feeds {
		keys = append(keys, k)
	}
	return keys
}
func (b *Analyzer) downloadAndCategorizeFeeds() {
	for _, feed := range b.Feeds {
		go func(feed *storage.Feed) {
			if len(feed.Subscribers) > 0 {
				// TODO -- check here if the feeds last download timestamp is older that x
				log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed.Source.FeedLink}).Info("Downloading RSS Feeds")
				fetchedFeed, err := fetch(feed.Source.FeedLink)
				if err != nil {
					log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed.Source.FeedLink, "error": err.Error()}).Error("Failed downloading feed")
					return
				}
				feed.DownloadTimestamp = time.Now()
				if fetchedFeed.FeedLink == "" {
					fetchedFeed.FeedLink = feed.Source.FeedLink
				}
				b.categorizeFeed(fetchedFeed)
			} else {
				log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed.Source.FeedLink}).Info("skipping feed. no subscriber")
			}

		}(feed)
	}
}

// check if rss feeds should be downloaded
func (b *Analyzer) tickerShouldDownloadFeeds() bool {
	// load last download timestamp
	lastDownloadTime := b.Db.GetFeedLastDownloadTime()
	do := true
	if !lastDownloadTime.IsZero() {
		if lastDownloadTime.After(time.Now().Add(-(b.RefreshPeriodDuration))) {
			do = false
		}
	}
	return do
}

// first try to download all user feeds, then start a download ticker based on configurable refresh rate
func (b *Analyzer) startFeedDownloadTicker() {
	ticker := time.NewTicker(b.RefreshPeriodDuration)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if b.tickerShouldDownloadFeeds() {
					//b.downloadAndCategorizeFeeds(b.getFeeds())
					b.downloadAndCategorizeFeeds()
					b.Db.SetFeedLastDownloadTime(time.Now())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
