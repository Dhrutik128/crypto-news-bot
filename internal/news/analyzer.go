package news

import (
	"encoding/json"
	"fmt"
	"github.com/drankou/go-vader/vader"
	"github.com/gohumble/crypto-news-bot/internal/sentiment"
	"github.com/mmcdole/gofeed"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
	"sync"
	"time"
)

type Analyzer struct {
	RefreshRate       time.Duration
	Feeds             map[string]*gofeed.Feed
	SentimentCompiler map[string]*sentiment.Compiler
	Mutex             sync.Mutex
	SentimentAnalyzer vader.SentimentIntensityAnalyzer
	Channels          Channels
	Sources           [][]string
	Db                *bitcask.Bitcask
}

func NewAnalyzer(db *bitcask.Bitcask, refreshRate time.Duration) *Analyzer {
	// initialize the vader sentiment analyzer first
	sia := vader.SentimentIntensityAnalyzer{}
	err := sia.Init("vader_lexicon.txt", "emoji_utf8_lexicon.txt")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// currently there is no more need for this channels, but they are still included.
	//Either will be used in future or get deprecated.
	c := Channels{
		BroadCastChannel: make(chan BroadCast, 10),
		SentimentChannel: make(chan *sentiment.Sentiment, 200),
		FeedChannel:      make(chan *gofeed.Feed, 200)}

	return &Analyzer{
		RefreshRate:       refreshRate,
		Db:                db,
		Mutex:             sync.Mutex{},
		Feeds:             make(map[string]*gofeed.Feed, 0),
		SentimentCompiler: make(map[string]*sentiment.Compiler, 0),
		SentimentAnalyzer: sia,
		Channels:          c,
	}

}

type BroadCast struct {
	User      *telebot.User
	Sentiment *sentiment.Sentiment
}

var KeyWords = [][]string{
	{"BTC", "bitcoin", "Bitcoin", "BITCOIN", "Satoshi"},
	{"XRP", "ripple", "xrp", "Ripple", "RIPPLE"},
	{"XLM", "Stellar Lumens"},
	{"BCH", "Bitcoin Cash"},
	{"DOGE", "Dogecoin", "dogecoin", "DOGE Coin"},
	{"ETH", "Ethereum"},
	{"BNB", "Binance Coin"},
	{"LTC", "Litecoin", "lite coin", "LiteCoin"},
	{"ZEC", "Zcash"},
	{"FORTH", "Ampleforth Governance Token"},
	{"FIL", "Filecoin"},
	{"Uniswap", "uniswap", "UniSwap"},
	{"BTG", "Bitcoin Gold", "BITCOIN GOLD"},
}

type Channels struct {
	FeedChannel      chan *gofeed.Feed
	SentimentChannel chan *sentiment.Sentiment
	BroadCastChannel chan BroadCast
}

func (b *Analyzer) GetSentimentTable() string {
	header := "```| Symbol | Sentiment |\n|--------|-------|\n"
	for coin, compiler := range b.SentimentCompiler {
		if len(compiler.Items) > 0 {
			header = header + fmt.Sprintf("| %s | %f |\n", coin, compiler.Avg)
		}
	}
	header = header + " ```"
	return header
}

func (b *Analyzer) categorize(sentiment *sentiment.Sentiment) error {
	if !b.Db.Has(sentiment.Key()) {
		b.Mutex.Lock()
		log.WithFields(log.Fields{"module": "[ANALYZER]", "title": sentiment.FeedItem.Title}).Info("Categorizing new item")
		b.categorizeFeedItem(sentiment)
		b.Mutex.Unlock()
		return nil
	}
	return fmt.Errorf("sentiment already in sotrage")
}
func (b *Analyzer) categorizeFeedItem(s *sentiment.Sentiment) {
	itemHash := fmt.Sprintf("%x", s.Hash)
	for _, words := range KeyWords {
		coin := words[0]
		compiler := &sentiment.Compiler{Items: make(map[string]*sentiment.Sentiment, 0)}
		if b.SentimentCompiler[coin] == nil {
			b.SentimentCompiler[coin] = compiler
		}
		if b.SentimentCompiler[coin].Items[itemHash] == nil {
			if contains(s.FeedItem.Title, words) {
				s.Sentiment = b.SentimentAnalyzer.PolarityScores(s.FeedItem.Title)
				s.Coin = coin
				compiler.Items[itemHash] = s
				b.SentimentCompiler[coin].Items[itemHash] = compiler.Items[itemHash]
				log.WithFields(log.Fields{"module": "[ANALYZER]", "title": s.FeedItem.Title, "link": s.FeedItem.Link}).Info("successfully ran sentiment analysis")
			}

		}
	}
}

func (b *Analyzer) categorizeFeedItemFromStorage(hashBytes []byte) error {
	item, err := b.Db.Get(hashBytes)
	if err != nil {
		return err
	}
	s := &sentiment.Sentiment{}
	err = json.Unmarshal(item, s)
	if err != nil {
		return err
	}
	if s.Sentiment != nil {
		if b.SentimentCompiler[s.Coin] == nil {
			b.SentimentCompiler[s.Coin] = &sentiment.Compiler{Items: make(map[string]*sentiment.Sentiment, 0)}
		}
		b.SentimentCompiler[s.Coin].Items[fmt.Sprintf("%x", hashBytes)] = s
	}

	return nil
}

func (b *Analyzer) categorizeFeed(feed *gofeed.Feed) {
	for _, feedItem := range feed.Items {
		s := &sentiment.Sentiment{FeedItem: feedItem, Feed: feed.FeedLink}
		err := b.categorize(s)
		if err != nil {
			continue
		}
		if s.Sentiment != nil {
			broadCastSentiment(s, b.Channels.BroadCastChannel)
		}
		sentiment.Save(s, b.Db)
	}
	for _, compiler := range b.SentimentCompiler {
		compiler.Compile()
	}

}

func (b *Analyzer) loadPersistedFeedItems() {
	b.Db.Scan([]byte("sentiment_"), func(key []byte) error {
		return b.categorizeFeedItemFromStorage(key)
	})
	for _, compiler := range b.SentimentCompiler {
		compiler.Compile()
	}
}

func (b *Analyzer) Start() {
	b.loadPersistedFeedItems()
	b.startFeedDownloadTicker()
}

func contains(title string, words []string) bool {
	found := false
	for _, word := range words {
		found = strings.Contains(title, word)
		if found {
			return true
		}
	}
	return found
}
