package telegram

import (
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/prologic/bitcask"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

// create a new telegram bot, initialize handlers and start broadcaster for broadcasting news
func New(db *bitcask.Bitcask, analyzer *news.Analyzer, token string) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		panic(err)
	}
	InitHandler(bot, db, analyzer)
	StartBroadCaster(analyzer, bot, analyzer.Channels.BroadCastChannel)
	return bot
}
