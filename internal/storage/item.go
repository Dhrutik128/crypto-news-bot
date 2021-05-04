package storage

import (
	"crypto/sha256"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

// a FeedItem is a analyzed single feed item
type FeedItem struct {
	// feed source
	// todo -- may use storage.Feed here or at least use url.URL instead of string type
	Feed *url.URL `json:"feed"`
	// the feed item. (also included in storage.Feed)
	Item *gofeed.Item `json:"item"`
	// sentiment analysis for this item
	Sentiment map[string]float64 `json:"sentiment"`
	// hash key used for storage
	HashKey []byte `json:"hash_key"`
	// the coin
	Coin string `json:"coin"`
	// may be removed. was meant to prevent multiple broadcasts
	WasBroadcast bool `json:"was_broadcast"`
}

func SaveFeedItem(sentiment *FeedItem, db *DB) {
	err := db.Set(sentiment)
	if err != nil {
		log.WithFields(log.Fields{"module": "[PERSISTANCE]", "error": err.Error()}).Info("failed persisting sentiment")
	}
}
func (s FeedItem) String() string {
	sb := &strings.Builder{}
	table := tablewriter.NewWriter(sb)
	table.Append([]string{s.Item.Title, s.Item.PublishedParsed.String(), fmt.Sprintf("%f", s.Sentiment["compound"])})
	table.SetHeader([]string{"Title", "Published", "Item"})

	table.Render()
	return sb.String()
}
func (s *FeedItem) hash() {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", s.Item.Title)))
	s.HashKey = append([]byte("sentiment_"), h.Sum(nil)...)
}
func (s *FeedItem) Key() []byte {
	if len(s.HashKey) > 0 {
		return s.HashKey
	} else {
		s.hash()
		return s.HashKey
	}

}

type sortedNewsFeed []*FeedItem

func (p sortedNewsFeed) Len() int {
	return len(p)
}

func (p sortedNewsFeed) Less(i, j int) bool {
	return p[i].Item.PublishedParsed.Before(p[j].Item.PublishedParsed.Local())
}

func (p sortedNewsFeed) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
