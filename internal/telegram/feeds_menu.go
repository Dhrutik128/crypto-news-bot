package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"strings"
)

var (
	FeedsMenu       = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	FeedsButton     = FeedsMenu.Text("/feeds")
	FeedsSelector   = &tb.ReplyMarkup{}
	FeedsButtonsMap = make(map[string]tb.Btn, 0)
	FeedsButtons    = make([]tb.Btn, 0)
	menuItems       = []string{"reset", "list", "top100"}
	helpText        = "Please provide a *comma seperated* list of  rss feed urls i should scrape for you. \n" +
		"By default users are *subscribed* to top 100 crypto sites. \n" +
		"\nUsage: \n" +
		"/feeds *add {url}* - add a new rss feed to your subscriptions \n" +
		"/feeds *remove {url}* - remove a rss feed from your subscriptions \n" +
		"/feeds *list* - list all subscribed feeds\n" +
		"/feeds *reset* - reset to default feed list\n" +
		"/feeds *top100* - returns a list of top 100 crypto sites"
)

func feedsButtonHandler(bot *tb.Bot, db *bitcask.Bitcask, analyzer *news.Analyzer) func(m *tb.Message) {
	return func(m *tb.Message) {

		if user, err := storage.UserRequired(m.Sender, db, bot); err == nil {

			if m.Text == "/feeds" {
				bot.Send(m.Sender, markdownEscape(helpText), FeedsSelector, tb.ModeMarkdownV2)
				return
			}
			s := strings.Split(m.Payload, " ")
			if len(s) >= 1 && len(s) <= 2 {
				switch s[0] {
				case "add":
					urls := s[1]
					urlSlice := strings.Split(urls, ",")
					for _, feedUrl := range urlSlice {
						log.Print("adding ", feedUrl)
						u, err := url.Parse(feedUrl)
						if err != nil {
							bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not parse %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2)
							return
						}
						err = analyzer.AddFeed(u, user)
						if err != nil {
							//bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not add feed %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2)
						}

					}
				case "remove":
				case "list":
					inlineButtonsHandler(bot, db, analyzer, user, "list")
				case "reset":
					inlineButtonsHandler(bot, db, analyzer, user, "reset")
				}
			}
		}
	}
}
func inlineButtonsHandler(bot *tb.Bot, db *bitcask.Bitcask, analyzer *news.Analyzer, user *storage.User, command string) {
	switch command {
	case "reset":
		user.Settings.Feeds = news.DefaultFeed
		storage.StoreUser(user, db)
		bot.Send(user.User, "feed list set to default", tb.ModeMarkdownV2)
	case "top100":
		bot.Send(user.User, markdownEscape(fmt.Sprintf("%s", strings.Join(news.DefaultFeed, ","))), tb.ModeMarkdownV2)
	case "list":
		bot.Send(user.User, markdownEscape(fmt.Sprintf("%s", strings.Join(user.Settings.Feeds, ","))), tb.ModeMarkdownV2)
	}
}
func initFeedsHandler(bot *tb.Bot, db *bitcask.Bitcask, analyzer *news.Analyzer) {
	FeedsButtons, FeedsButtonsMap = getButtons("feeds_", menuItems, FeedsMenu)
	FeedsSelector.Inline(ButtonWrapper(FeedsButtons, FeedsSelector)...)
	bot.Handle(&FeedsButton, feedsButtonHandler(bot, db, analyzer))
	for _, btn := range FeedsButtonsMap {
		bot.Handle(&btn, func(c *tb.Callback) {
			if user, err := storage.UserRequired(c.Sender, db, bot); err == nil {
				if len(c.Data) > 0 {
					inlineButtonsHandler(bot, db, analyzer, user, c.Data)
				}
			}
		})
	}
}
