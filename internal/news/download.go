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

// download all feeds, that users have added
func (b *Analyzer) downloadUserFeeds() {
	//b.downloadAndCategorizeFeeds(b.getUserFeeds())

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
				log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed}).Info("Downloading RSS Feeds")
				fetchedFeed, err := fetch(feed.Source.Link)
				if err != nil {
					log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed, "error": err.Error()}).Error("Failed downloading feed")
					return
				}
				// add the fresh feeds to slice.
				if fetchedFeed.FeedLink != feed.Source.FeedLink {
					fetchedFeed.FeedLink = feed.Source.FeedLink
				}
				// todo -- to increase efficiency, slice should only be updated, when fetched and stored feeds are not equal
				/*b.Mutex.Lock()
				b.Feeds[feed] = fetchedFeed
				b.Mutex.Unlock()*/
				b.categorizeFeed(fetchedFeed)
			} else {
				log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed}).Info("skipping feed. no subscriber")
			}

		}(feed)
	}
}

func broadCastSentiment(sentiment *storage.Sentiment, broadcastChannel chan BroadCast) {
	if !sentiment.WasBroadcast {
		if sentiment.FeedItem.PublishedParsed != nil {
			if sentiment.FeedItem.PublishedParsed.After(time.Now().Add(-(time.Hour * 24))) {
				broadcastChannel <- BroadCast{Sentiment: sentiment}
				// prevents sending same feed item in broadcast for another coin subscription
				sentiment.WasBroadcast = true
			}
		}
	}
}

// download feeds and set lastDownloadTime
func (b *Analyzer) tickerTryDownload() {
	if b.tickerShouldDownloadFeed() {
		//b.downloadAndCategorizeFeeds(b.getFeeds())
		b.downloadAndCategorizeFeeds()
		b.Db.SetFeedLastDownloadTime(time.Now())
	}
}

// check if rss feeds should be downloaded
func (b *Analyzer) tickerShouldDownloadFeed() bool {
	// load last download timestamp
	lastDownloadTime := b.Db.GetFeedLastDownloadTime()
	do := true
	if !lastDownloadTime.IsZero() {
		if lastDownloadTime.After(time.Now().Add(-(b.RefreshRate))) {
			do = false
		}
	}
	return do
}

// first try to download all user feeds, then start a download ticker based on configurable refresh rate
func (b *Analyzer) startFeedDownloadTicker() {

	b.tickerTryDownload()
	ticker := time.NewTicker(b.RefreshRate)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				b.tickerTryDownload()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
