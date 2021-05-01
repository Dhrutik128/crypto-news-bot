package news

import (
	"github.com/gohumble/crypto-news-bot/internal/storage"
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

// this should remove the user from the feed object stored in the news analyzer
// if the removed user was the last one, we must also remove the feed from the analyzer and storage.
// otherwise we would download feeds, that have no subscribers.
func (b *Analyzer) RemoveFeed(source *url.URL, user *storage.User) error {
	// todo -- implement this
	feed := b.Feeds[source.String()]
	if feed != nil {
		feed.RemoveUser(user)
		if len(feed.Subscribers) == 0 {
			// remove the feed from news analyzer, if no user is currently subscribed.
			// this will prevent the bot from downloading feeds without users.
			delete(b.Feeds, source.String())
			err := b.Db.Delete(feed)
			if err != nil {
				return err
			}
		}
		err := storage.SetFeed(feed, b.Db)
		if err != nil {
			return err
		}
	}
	err := user.RemoveFeed(source.String(), b.Db)
	if err != nil {
		return err
	}
	return b.Db.Set(user)

}

// if feed does not exists in the news analyzer, we should fetch the feed, add the current user
// store the feed and run the analytics.
// if feed is already included in the news analyzer, we just add the user and update the feed in storage.
func (b *Analyzer) AddFeed(source *url.URL, user *storage.User) error {
	// case 1 -- new feed
	if b.Feeds[source.String()] == nil {
		feed, err := fetch(source.String())
		if err != nil {
			return err
		}
		// update the feed link if this is not present in feed
		if feed.FeedLink != source.String() {
			feed.FeedLink = source.String()
		}

		f := &storage.Feed{Source: *feed}
		if user != nil {
			f.AddUser(user)

			err := user.AddFeed(source.String(), b.Db)
			if err != nil {
				return err
			}
		}
		err = storage.ImportFeed(f, b.Db)
		if err != nil {
			return err
		}
		b.Feeds[source.String()] = f
		b.categorizeFeed(feed)
		return nil

	}
	if user != nil {
		b.Feeds[source.String()].AddUser(user)
		// case 2 -- feed already exists. user subscribes to existing feed!
		err := user.AddFeed(source.String(), b.Db)
		if err != nil {
			return err
		}
		err = storage.SetFeed(b.Feeds[source.String()], b.Db)
		if err != nil {
			// todo -- remove the fee from user when setFeed fails.
			return err
		}
	}

	return nil
}
