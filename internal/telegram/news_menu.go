package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	NewsMenu       = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	NewsButton     = NewsMenu.Text("/news")
	selector       = &tb.ReplyMarkup{}
	NewsButtonsMap = make(map[string]tb.Btn, 0)
	NewsButtons    = make([]tb.Btn, 0)

	btnBack = NewsMenu.Text("< Back")
)

func initNewsHandler(bot *tb.Bot, db *storage.DB, feed *news.Analyzer) {
	NewsButtons, NewsButtonsMap = getKeywordButtons("news_", NewsMenu)
	selector.Inline(ButtonWrapper(NewsButtons, selector)...)
	bot.Handle(&NewsButton, func(m *tb.Message) {
		bot.Send(m.Sender, "Choose a coin", selector)
	})
	for _, btn := range NewsButtonsMap {
		bot.Handle(&btn, func(c *tb.Callback) {
			if _, err := storage.UserRequired(c.Sender, db, bot); err == nil {
				SendNews(c, bot, feed)
			}
		})
	}
	bot.Handle(&btnBack, func(c *tb.Callback) {
		bot.Send(c.Sender, "Main Menu", menu)
	})
}

func SendNews(c *tb.Callback, bot *tb.Bot, newsFeed *news.Analyzer) {
	if c.Data != "" {
		if newsFeed.SentimentCompiler[c.Data] != nil {
			log.WithFields(log.Fields{"module": "[TELEGRAM]"}).Infof("sending latest news to %s", c.Sender.Username)
			bot.Send(c.Sender, fmt.Sprintf("Hi %s \n*sending all processed news for %s*\n", c.Sender.Username, c.Data), tb.ModeMarkdownV2)
			latestNews := newsFeed.SentimentCompiler[c.Data].GetNews()
			for _, n := range latestNews {
				text := fmt.Sprintf("\n[*_%s_*](%s)\n\n*Published:* %s\n*Sentiment:* %s\n",
					markdownEscape(n.FeedItem.Title),
					markdownEscape(n.FeedItem.Link),
					markdownEscape(n.FeedItem.Published),
					markdownEscape(fmt.Sprintf("%f",
						n.Sentiment["compound"])))
				_, err := bot.Send(c.Sender, text, tb.ModeMarkdownV2)
				if err != nil {
					log.WithFields(log.Fields{"error": err.Error()}).Errorf("could not sent news")
				}
			}
		}
	}
}
