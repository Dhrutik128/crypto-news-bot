package sentiment

import (
	"fmt"
	"github.com/drankou/go-vader/vader"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Compiler struct {
	Items map[string]*Sentiment
	Avg   float64
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

func Do() {
	sia := vader.SentimentIntensityAnalyzer{}
	err := sia.Init("vader_lexicon.txt", "emoji_utf8_lexicon.txt")
	if err != nil {
		log.Fatal(err)
	}

	score := sia.PolarityScores("XRP Back in a wedge? Ready for another rally?")
	fmt.Println(score)
}

func (c *Compiler) Compile() {
	var sum float64
	for _, item := range c.Items {
		sum = sum + item.Sentiment["compound"]
	}
	c.Avg = (c.Avg + sum) / float64(len(c.Items))
}
func (c Compiler) string() string {
	return fmt.Sprintf("%f", c.Avg)
}
