package storage

import (
	"fmt"
	"sort"
	"sync"
)

type SentimentCompiler struct {
	Items map[string]*FeedItem
	Avg   float64
	Mutex *sync.Mutex
}

func NewSentimentCompiler() *SentimentCompiler {
	return &SentimentCompiler{Mutex: &sync.Mutex{}, Items: make(map[string]*FeedItem, 0)}
}
func (c SentimentCompiler) sorted() []*FeedItem {
	news := make(sortableFeedItems, 0)
	for _, s := range c.Items {
		if s.Item.PublishedParsed == nil {
			continue
		}
		news = append(news, s)
	}
	sort.Sort(news)
	return news
}
func (c SentimentCompiler) GetNews() []*FeedItem {
	sortedNews := c.sorted()
	if len(sortedNews) > 10 {
		return sortedNews[len(sortedNews)-10:]
	}
	return sortedNews
}

func (c *SentimentCompiler) Compile() {
	var sum float64
	for _, item := range c.Items {
		sum = sum + item.Sentiment["compound"]
	}
	c.Avg = sum / float64(len(c.Items))
}
func (c SentimentCompiler) string() string {
	return fmt.Sprintf("%f", c.Avg)
}
