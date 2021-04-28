package sentiment

import (
	"fmt"
	"sort"
	"sync"
)

type Compiler struct {
	Items map[string]*Sentiment
	Avg   float64
	Mutex *sync.Mutex
}

func NewCompiler() *Compiler {
	return &Compiler{Mutex: &sync.Mutex{}, Items: make(map[string]*Sentiment, 0)}
}
func (sc Compiler) sorted() []*Sentiment {
	news := make(sortedNewsFeed, 0)
	for _, s := range sc.Items {
		if s.FeedItem.PublishedParsed == nil {
			continue
		}
		news = append(news, s)
	}
	sort.Sort(news)
	return news
}
func (sc Compiler) GetNews() []*Sentiment {
	sortedNews := sc.sorted()
	if len(sortedNews) > 10 {
		return sortedNews[len(sortedNews)-10:]
	}
	return sortedNews
}

func (c *Compiler) Compile() {
	if len(c.Items) > 0 {
		var sum float64
		c.Mutex.Lock()
		for _, item := range c.Items {
			sum = sum + item.Sentiment["compound"]
		}
		c.Mutex.Unlock()
		c.Avg = sum / float64(len(c.Items))
	}
}
func (c Compiler) string() string {
	return fmt.Sprintf("%f", c.Avg)
}
