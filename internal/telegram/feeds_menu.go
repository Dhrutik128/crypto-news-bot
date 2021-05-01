package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"strings"
)

var (
	FeedsMenu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	FeedsButton   = FeedsMenu.Text("/feeds")
	FeedsSelector = &tb.ReplyMarkup{}
	//FeedsButtonsMap = make(map[string]tb.Btn, 0)
	//FeedsButtons    = make([]tb.Btn, 0)
	menuItems = []string{"reset", "list", "top100"}
	helpText  = "Please provide a *comma seperated* list of  rss feed urls i should scrape for you. \n" +
		"By default users are *subscribed* to top 100 crypto sites. \n" +
		"\nUsage: \n" +
		"/feeds *add {url}* - add a new rss feed to your subscriptions \n" +
		"/feeds *remove {url}* - remove a rss feed from your subscriptions \n" +
		"/feeds *list* - list all subscribed feeds\n" +
		"/feeds *reset* - reset to default feed list\n"
)

func feedsCommandHandler(bot *tb.Bot, db *storage.DB, analyzer *news.Analyzer, callback *tb.Callback) func(m *tb.Message) {
	return func(m *tb.Message) {
		if user, err := storage.UserRequired(m.Sender, db, bot); err == nil {
			if m.Text == "/feeds" {
				//markup := &tb.ReplyMarkup{}
				//getDefaultFeedButtons("feeds_", menuItems, markup, user)
				_, err := bot.Send(m.Sender, markdownEscape(helpText))
				if err != nil {
					fmt.Println(err)
				}
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
							bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not add feed %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2)
						}
					}
				case "remove":
					urls := s[1]
					urlSlice := strings.Split(urls, ",")
					for _, feedUrl := range urlSlice {
						log.Print("removing ", feedUrl)
						u, err := url.Parse(feedUrl)
						if err != nil {
							bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not parse %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2)
							return
						}
						err = analyzer.RemoveFeed(u, user)
						if err != nil {
							//bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not add feed %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2)
						}
					}
				case "list":

					/*if user.Settings.IsDefaultFeedSubscribed {
						feeds = append(feeds, news.DefaultFeed...)
					}*/
					if len(user.Settings.Feeds) > 0 {
						config.IgnoreErrorMultiReturn(
							bot.Send(user.User,
								markdownEscape(fmt.Sprintf("%s", strings.Join(unique(user.Settings.Feeds), ", "))),
								tb.ModeMarkdownV2))
					}
					//inlineButtonsHandler(bot, db, callback, user, analyzer,"list")
				case "reset":
					for _, f := range user.Settings.Feeds {
						analyzer.Feeds[f].RemoveUser(user)
					}
					user.Settings.Feeds = news.DefaultFeed
					analyzer.AddUserToDefaultFeeds(user)
					db.Set(user)
					bot.Send(user.User, "feed list set to default", tb.ModeMarkdownV2)
				}
			}
		}
	}
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

/*
func inlineButtonsHandler(bot *tb.Bot, db *storage.DB, c *tb.Callback, user *storage.User, analyzer *news.Analyzer,command string) {
	switch command {
	case "reset":
		for _, f := range user.Settings.Feeds {
			analyzer.Feeds[f].RemoveUser(user)
		}
		user.Settings.Feeds = news.DefaultFeed
		db.Set(user)
		bot.Send(user.User, "feed list set to default", tb.ModeMarkdownV2)
	case "top100":
		err := user.ToggleDefaultFeed(db)
		if err != nil {
			return
		}
		markup := &tb.ReplyMarkup{}
		getDefaultFeedButtons("feeds_", menuItems, markup, user)
		config.IgnoreErrorMultiReturn(bot.EditReplyMarkup(c.Message, markup))
	case "list":
		feeds := make([]string, 0)

		for _, feed := range user.Settings.Feeds {
			feeds = append(feeds, feed)
		}
		config.IgnoreErrorMultiReturn(
			bot.Send(user.User,
				markdownEscape(fmt.Sprintf("%s", strings.Join(unique(feeds), ", "))),
				tb.ModeMarkdownV2))
	}
}
*/
func initFeedsHandler(bot *tb.Bot, db *storage.DB, analyzer *news.Analyzer) {
	//FeedsButtons, FeedsButtonsMap = getDefaultFeedButtons("feeds_", menuItems, FeedsMenu, nil)
	//FeedsSelector.Inline(ButtonWrapper(FeedsButtons, FeedsSelector)...)
	bot.Handle(&FeedsButton, feedsCommandHandler(bot, db, analyzer, nil))

	/*for _, btn := range FeedsButtonsMap {
		bot.Handle(&btn, func(c *tb.Callback) {
			if user, err := storage.UserRequired(c.Sender, db, bot); err == nil {
				if len(c.Data) > 0 {
					inlineButtonsHandler(bot, db, c, user,analyzer, c.Data)
				}
			}
		})
	}*/
}
