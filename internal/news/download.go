package news

import (
	"context"
	"encoding/json"
	"github.com/gohumble/crypto-news-bot/internal/sentiment"
	"github.com/mmcdole/gofeed"
	"github.com/prologic/bitcask"
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
	b.downloadAndCategorizeFeeds(b.getUserFeeds())
}

func (b *Analyzer) downloadAndCategorizeFeeds(feeds []string) {
	if len(feeds) > 0 {
		for _, feed := range feeds {
			fetchedFeed, err := fetch(feed)
			if err != nil {
				log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed, "error": err.Error()}).Error("Failed downloading feed")
				continue
			}
			// add the fresh feeds to slice.
			if fetchedFeed.FeedLink != feed {
				fetchedFeed.FeedLink = feed
			}
			// todo -- to increase efficiency, slice should only be updated, when fetched and stored feeds are not equal
			b.Feeds[feed] = fetchedFeed
			b.categorizeFeed(fetchedFeed)
			log.WithFields(log.Fields{"module": "[DOWNLOAD]", "link": feed}).Info("Downloading Users RSS Feeds")
		}

	}
}

// download all feeds, that users have added
func (b *Analyzer) downloadDefaultFeeds() {
	b.downloadAndCategorizeFeeds(DefaultFeed)
}
func saveSentiment( sentiment *sentiment.Sentiment, db *bitcask.Bitcask) {
	sentimentBytes, err := json.Marshal(sentiment)
	if err != nil {
		return
	}
	db.Put(sentiment.Hash, sentimentBytes)
}
func broadCastSentiment(sentiment *sentiment.Sentiment, broadcastChannel chan BroadCast) {
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
		b.downloadDefaultFeeds()
		b.downloadUserFeeds()

		b.Db.Put([]byte( "lastDownloadTime"), []byte(time.Now().Format(time.RFC3339)))
	}
}

// check if rss feeds should be downloaded
func (b *Analyzer) tickerShouldDownloadFeed() bool {
	// load last download timestamp
	lastDownloadTime, _ := b.Db.Get([]byte( "lastDownloadTime"))
	do := true
	if len(lastDownloadTime) > 0 {
		t, err := time.Parse(time.RFC3339, string(lastDownloadTime))
		if err != nil {
			log.Println(err)
			do = true
		}
		if t.After(time.Now().Add(-(b.RefreshRate))) {
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
