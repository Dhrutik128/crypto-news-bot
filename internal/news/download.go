package news

import (
	"context"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"github.com/gohumble/crypto-news-bot/internal/safe"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
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

// downloads and processes all previously added feeds again
func (b *Analyzer) downloadAndCategorizeFeeds() {
	wg := &sync.WaitGroup{}
	for _, feed := range b.Feeds {
		wg.Add(1)
		requestContext := context.WithValue(context.Background(), "ref", feed.Source.String())
		b.Pool.GoCtx(safe.NewRoutineWithContext(func(ctx context.Context, routine safe.RoutineCtx) {
			feedUrl, err := url.Parse(feed.Source.FeedLink)
			if err != nil {
				log.WithFields(log.Fields{"feed": feed, "error": err.Error()}).Error("could not parse feed url")
				return
			}
			config.IgnoreError(b.addSource(feedUrl, nil, wg, false))
		}, requestContext))
	}
	wg.Wait()
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
	tryDownload := func() {
		if b.tickerShouldDownloadFeeds() {
			//b.downloadAndCategorizeFeeds(b.getFeeds())
			b.downloadAndCategorizeFeeds()
			b.Db.SetFeedLastDownloadTime(time.Now())
		}
	}
	//tryDownload()
	ticker := time.NewTicker(b.RefreshPeriodDuration)
	quit := make(chan struct{})
	tryDownload()
	go func() {
		for {
			select {
			case <-ticker.C:
				tryDownload()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
