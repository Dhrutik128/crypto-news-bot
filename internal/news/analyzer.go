package news

import (
	"encoding/json"
	"fmt"
	"github.com/drankou/go-vader/vader"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/mmcdole/gofeed"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"gopkg.in/tucnak/telebot.v2"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Analyzer will hold all feeds and their analyzed items (sentiments)
type Analyzer struct {
	RefreshRate       time.Duration
	Feeds             map[string]*storage.Feed
	SentimentCompiler map[string]*storage.Compiler
	Mutex             sync.Mutex
	SentimentAnalyzer vader.SentimentIntensityAnalyzer
	Channels          Channels
	Sources           [][]string
	Db                *storage.DB
}

// Channels BroadCastChannel accepts Broadcasts that will be sent to users if relevant.
type Channels struct {
	BroadCastChannel chan BroadCast
}

// BroadCast struct contains the user and the analyzed feed item
type BroadCast struct {
	User     *telebot.User
	FeedItem *storage.FeedItem
}

// NewAnalyzer returns a new Analyzer struct. It will not contain any feed information.
func NewAnalyzer(db *storage.DB, refreshRate time.Duration) *Analyzer {
	// initialize the vader sentiment analyzer first
	sia := vader.SentimentIntensityAnalyzer{}
	err := sia.Init("vader_lexicon.txt", "emoji_utf8_lexicon.txt")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// currently there is no more need for this channels, but they are still included.
	c := Channels{BroadCastChannel: make(chan BroadCast, 10)}

	return &Analyzer{
		RefreshRate:       refreshRate,
		Db:                db,
		Mutex:             sync.Mutex{},
		Feeds:             make(map[string]*storage.Feed, 0),
		SentimentCompiler: make(map[string]*storage.Compiler, 0),
		SentimentAnalyzer: sia,
		Channels:          c,
	}

}

// KeyWords used by analyzer to search for relevant information within a feed item.
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

// GetSentimentTable returns the processed sentiment table
func (b *Analyzer) GetSentimentTable() string {
	sb := &strings.Builder{}
	table := tablewriter.NewWriter(sb)
	table.SetHeader([]string{"Symbol", "FeedItem"})
	for coin, compiler := range b.SentimentCompiler {
		if len(compiler.Items) > 0 {
			table.Append([]string{coin, fmt.Sprintf("%f", compiler.Avg)})
		}
	}
	table.Render()
	return sb.String()
}

// shouldCategorize checks if the item was processed before
func (b *Analyzer) shouldCategorize(feedItem *storage.FeedItem) error {
	if ok, _ := b.Db.Exists(feedItem); !ok {
		b.Mutex.Lock()
		log.WithFields(log.Fields{"module": "[ANALYZER]", "title": feedItem.Item.Title}).Info("Categorizing new item")
		b.categorizeFeedItem(feedItem)
		b.Mutex.Unlock()
		return nil
	}
	return fmt.Errorf("feedItem already in sotrage")
}

// categorizeFeedItem by searching keywords and applying sentiment analysis when keywords match.
// the categorized feed item is then updated within the sentiment compiler
func (b *Analyzer) categorizeFeedItem(s *storage.FeedItem) {
	itemHash := fmt.Sprintf("%x", s.HashKey)
	// get all coin keywords
	for _, words := range KeyWords {
		coin := words[0]
		compiler := storage.NewCompiler()
		if b.SentimentCompiler[coin] == nil {
			b.SentimentCompiler[coin] = compiler
		}
		if b.SentimentCompiler[coin].Items[itemHash] == nil {
			if contains(s.Item.Title, words) {
				// feed title contains a coin keyword so we add sentiment analysis
				s.Sentiment = b.SentimentAnalyzer.PolarityScores(s.Item.Title)
				s.Coin = coin
				compiler.Items[itemHash] = s
				b.SentimentCompiler[coin].Items[itemHash] = compiler.Items[itemHash]
				log.WithFields(log.Fields{"module": "[ANALYZER]", "title": s.Item.Title, "link": s.Item.Link}).Info("successfully ran sentiment analysis")
			}

		}
	}
}

// categorizeFeedItemFromStorage updates the feed item in the sentiment compiler
func (b *Analyzer) categorizeFeedItemFromStorage(object storage.Storable) error {
	s := object.(*storage.FeedItem)
	if s.Sentiment != nil {
		if b.SentimentCompiler[s.Coin] == nil {
			b.SentimentCompiler[s.Coin] = storage.NewCompiler()
		}
		b.SentimentCompiler[s.Coin].Items[fmt.Sprintf("%x", object.Key())] = s
	}

	return nil
}

// categorizeFeed categorizes all feed items
func (b *Analyzer) categorizeFeed(feed *gofeed.Feed) {
	feedUrl, err := url.Parse(feed.FeedLink)
	if err != nil {
		return
	}
	for _, feedItem := range feed.Items {
		item := &storage.FeedItem{Item: feedItem, Feed: feedUrl}
		err := b.shouldCategorize(item)
		if err != nil {
			continue
		}
		if item.Sentiment != nil {
			trySendBroadCast(item, b.Channels.BroadCastChannel)
		}
		storage.SaveFeedItem(item, b.Db)
	}
	for _, compiler := range b.SentimentCompiler {
		b.Mutex.Lock()
		compiler.Compile()
		b.Mutex.Unlock()
	}

}

// trySendBroadCast trough broadcastChannel. This broadcast does not contain a user yet.
// users will be set by the broadcast receiver
func trySendBroadCast(feedItem *storage.FeedItem, broadcastChannel chan BroadCast) {
	if !feedItem.WasBroadcast {
		if feedItem.Item.PublishedParsed != nil {
			if feedItem.Item.PublishedParsed.After(time.Now().Add(-(time.Hour * 24))) {
				broadcastChannel <- BroadCast{FeedItem: feedItem}
				// prevents sending same feed item in broadcast for another coin subscription
				feedItem.WasBroadcast = true
			}
		}
	}
}

// loadPersistedItems from storage and add them to news analyzer (when loading the application)
func (b *Analyzer) loadPersistedItems() {
	// loading all processed feed items from past 3 days
	err := b.Db.View(func(tx *buntdb.Tx) error {
		err := tx.Descend("sentiment", func(key, value string) bool {
			item := &storage.FeedItem{HashKey: []byte(key)}
			err := json.Unmarshal([]byte(value), item)
			if err != nil {
				return true
			}
			// do not load news older than 3 days
			if item.Item.PublishedParsed.Before(time.Now().Add(-(time.Hour * 72))) {
				return false
			}
			config.IgnoreError(b.categorizeFeedItemFromStorage(item))
			return true
		})
		return err
	})
	if err == nil {
		// note -- persist compiler instead of sentiment ?!
		// compiling processed feed items again...
		for _, compiler := range b.SentimentCompiler {
			compiler.Compile()
		}
	}
	// load all feeds that users are subscribed to
	config.IgnoreError(b.Db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("feed", func(key, value string) bool {
			feed := &storage.Feed{}
			err := json.Unmarshal([]byte(value), feed)
			if err != nil {
				return true
			}
			if len(feed.Subscribers) > 0 {
				feedUrl, err := url.Parse(feed.Source.FeedLink)
				if err != nil {
					return true
				}
				b.Feeds[feedUrl.String()] = feed
			}
			return true
		})
		return err
	}))
	b.AddUserToDefaultFeeds(nil)

}

// AddUserToDefaultFeeds adds the user to all default feed items of the analyzer
func (b *Analyzer) AddUserToDefaultFeeds(user *storage.User) {
	if len(b.Feeds) == 0 {
		for _, feed := range DefaultFeed {
			feedUrl, err := url.Parse(feed)
			if err != nil {
				log.WithFields(log.Fields{"feed": feed, "error": err.Error()}).Error("could not parse feed url")
				continue
			}
			err = b.AddFeed(feedUrl, user, true)
			if err != nil {
				log.WithFields(log.Fields{"feed": feed, "error": err.Error()}).Error("could not add user to feed")
			}
		}
	}
}

// Start the analyzer
func (b *Analyzer) Start() {
	// this will load all feeds and processed feed items
	b.loadPersistedItems()
	// start the download ticker for previously loaded feeds
	b.startFeedDownloadTicker()
}

// contains checks if slice of strings contains a certain string
func contains(search string, words []string) bool {
	found := false
	for _, word := range words {
		found = strings.Contains(search, word)
		if found {
			return true
		}
	}
	return found
}
