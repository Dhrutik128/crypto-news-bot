package sentiment

import (
	"crypto/sha256"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/olekukonko/tablewriter"
	"strings"
)

type Sentiment struct {
	Feed         string
	FeedItem     *gofeed.Item
	Sentiment    map[string]float64
	Hash         []byte
	Coin         string
	WasBroadcast bool
}

func (s Sentiment) String() string {
	sb := &strings.Builder{}
	table := tablewriter.NewWriter(sb)
	table.Append([]string{s.FeedItem.Title, s.FeedItem.PublishedParsed.String(), fmt.Sprintf("%f", s.Sentiment["compound"])})
	table.SetHeader([]string{"Title", "Published", "Sentiment"})

	table.Render()
	return sb.String()
}
func (s *Sentiment) hash() {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", s.FeedItem.Title)))
	s.Hash = append([]byte("sentiment_"), h.Sum(nil)...)
}
func (s *Sentiment) Key() []byte {
	if len(s.Hash) > 0 {
		return s.Hash
	} else {
		s.hash()
		return s.Hash
	}

}

type sortedNewsFeed []*Sentiment

func (p sortedNewsFeed) Len() int {
	return len(p)
}

func (p sortedNewsFeed) Less(i, j int) bool {
	return p[i].FeedItem.PublishedParsed.Before(p[j].FeedItem.PublishedParsed.Local())
}

func (p sortedNewsFeed) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
