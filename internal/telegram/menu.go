package telegram

import (
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func init() {
	menu.Reply(
		menu.Row(btnHelp, btnSentiments, btnSubscribe),
		menu.Row(NewsButton, FeedsButton, btnGit),
	)
}

var (

	// Universal markup builders.
	menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	// Reply buttons.
	btnHelp       = menu.Text("/help")
	btnSubscribe  = menu.Text("/subscribe")
	btnGit        = menu.URL("github", "https://github.com/gohumble/crypto-news-bot")
	btnSentiments = menu.Text("/sentiments")
)

func InitHandler(bot *tb.Bot, db *bitcask.Bitcask, newsfeed *news.Analyzer) {
	// Command: /start <PAYLOAD>
	bot.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		_, err := bot.Send(m.Sender, markdownEscape("HI! I am a *Crypto News Stream Bot*. \n\n"+
			"I will fetch the latest news from RSS feeds and categorize and analyze them by coin.\n"+
			"You just  need to subscribe to coins using the /subscribe command.\n"+
			"Please not that news updates will be broadcast to subscribers every hour.\n\n"+
			"By default i will scrape the top 100 crypto sites and send you the latest news on subscribed coins.\n"+
			"If you want to manage your personal rss feeds, you can use the /feeds command."), menu, tb.ModeMarkdownV2)
		if err != nil {
			return
		}
		user := storage.User{User: m.Sender}
		if !db.Has(user.Key()) {
			err := storage.StoreUser(&storage.User{
				User:     m.Sender,
				Settings: storage.UserSettings{Subscriptions: make(map[string]bool, 0), Feeds: news.DefaultFeed}, Started: time.Now()},
				db)
			if err != nil {
				log.Println(err)
			}
		}

	})

	// On reply button pressed (message)
	bot.Handle(&btnHelp, func(m *tb.Message) {
		bot.Send(m.Sender, "Usage:\n\n"+
			"/subscribe - subscribe to coins\n"+
			"/feed - manage your RSS feeds \n"+
			"/news - get latest news for a single coin\n"+
			"/sentiments - get sentiment analysis based on all news for all subscribable coins.\n")
	})

	bot.Handle(&btnSentiments, func(m *tb.Message) {
		bot.Send(m.Sender, newsfeed.GetSentimentTable(), &tb.SendOptions{
			ParseMode: tb.ModeMarkdown,
		})
	})
	initNewsHandler(bot, db, newsfeed)
	initSubscriptionHandler(bot, db)
	initFeedsHandler(bot, db, newsfeed)
}
