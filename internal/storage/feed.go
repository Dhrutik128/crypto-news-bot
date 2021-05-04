package storage

import (
	"crypto/sha256"
	"fmt"
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	HashKey     []byte       `json:"hash_key"`
	Subscribers []int        `json:"subscribers"`
	Source      *gofeed.Feed `json:"source"`
	IsDefault   bool         `json:"is_default"`
}

// hash the feed struct using the feed link
func (f *Feed) hash() {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", f.Source.FeedLink)))
	f.HashKey = append([]byte("feed_"), h.Sum(nil)...)
}

// Key for storable
func (f *Feed) Key() []byte {
	if len(f.HashKey) > 0 {
		return f.HashKey
	} else {
		f.hash()
		return []byte(fmt.Sprintf("feed_%d", f.HashKey))
	}
}

// RemoveFeed from storage
func (f *Feed) RemoveFeed(db *DB) error {
	return db.Delete(f)
}

// ImportFeed initial import of feed. check if it is already imported / hash collision
func ImportFeed(feed *Feed, db *DB) error {
	if ok, _ := db.Exists(feed); !ok {
		return SetFeed(feed, db)
	}
	return fmt.Errorf("feed already exists")
}

// SetFeed for updating the feed in storage.
func SetFeed(feed *Feed, db *DB) error {
	items := feed.Source.Items
	feed.Source.Items = nil
	defer func() {
		feed.Source.Items = items
	}()
	return db.Set(feed)
}

func removeInt(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// RemoveUser from feed subscription.
// currently this will not update the storable
func (f *Feed) RemoveUser(user *User) {
	if len(f.Subscribers) > 1 {
		for i, u := range f.Subscribers {
			if u == user.User.ID {
				f.Subscribers = removeInt(f.Subscribers, i)
			}
		}
	} else {
		f.Subscribers = make([]int, 0)
	}
}

// AddUser to feed.
// currently this will not update the storable
func (f *Feed) AddUser(user *User) error {
	alreadyAdded := false
	for _, sub := range f.Subscribers {
		if sub == user.User.ID {
			alreadyAdded = true
			break
		}
	}
	if !alreadyAdded {
		f.Subscribers = append(f.Subscribers, user.User.ID)
		return nil
	}
	return fmt.Errorf("user is already included in feed")

}

// HasUser returns true is user is subscribed to feed
func (f Feed) HasUser(user User) bool {
	for _, u := range f.Subscribers {
		if u == user.User.ID {
			return true
		}
	}
	return false
}
