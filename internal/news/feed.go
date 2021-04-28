package news

import (
	"encoding/json"
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	log "github.com/sirupsen/logrus"
	"net/url"
)

func init() {
	df, err := readCsv("feeds.csv")
	if err != nil {
		panic(err)
	}

	df = df[1:]
	for _, f := range df {
		DefaultFeed = append(DefaultFeed, f[0])
	}

}

var DefaultFeed []string

func (b *Analyzer) getUserFeeds() []string {
	var userFeedsMap = make(map[string]struct{}, 0)
	var userFeeds = make([]string, 0)
	b.Db.Scan([]byte("user_"), func(key []byte) error {
		userBytes, err := b.Db.Get(key)
		if err != nil {
			return err
		}
		user := storage.User{}
		err = json.Unmarshal(userBytes, &user)
		if err != nil {
			log.Println(err)
			return err
		}
		for _, feed := range user.Settings.Feeds {
			userFeedsMap[feed] = struct{}{}
		}
		return nil
	})
	for feed, _ := range userFeedsMap {
		userFeeds = append(userFeeds, feed)
	}
	return userFeeds
}

// add feed to user and if feed does not exists, also add it to the analyzer
func (b *Analyzer) RemoveFeed(source *url.URL, user *storage.User) error {
	return user.RemoveFeed(source.String(), b.Db)
}
func (b *Analyzer) AddFeed(source *url.URL, user *storage.User) error {
	if b.Feeds[source.String()] == nil {
		feed, err := fetch(source.String())
		if err != nil {
			return err
		}
		err = user.AddFeed(source.String(), b.Db)
		if err != nil {
			return err
		}
		b.Feeds[source.String()] = feed
		// update the feed link if this is not present in feed
		if feed.FeedLink != source.String() {
			feed.FeedLink = source.String()
		}
		b.categorizeFeed(feed)

		return nil
	}
	err := user.AddFeed(source.String(), b.Db)
	if err != nil {
		return err
	}
	return fmt.Errorf("source already included in feed list")
}
