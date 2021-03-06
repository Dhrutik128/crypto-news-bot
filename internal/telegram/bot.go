package telegram

import (
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

// create a new telegram bot, initialize handlers and start broadcaster for broadcasting news
func New(db *storage.DB, analyzer *news.Analyzer, token string) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		panic(err)
	}
	InitHandler(bot, db, analyzer)
	StartUserBroadCaster(analyzer, bot, analyzer.Channels.BroadCastChannel)
	return bot
}
