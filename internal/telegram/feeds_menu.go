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
	"sync"
)

var (
	FeedsMenu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	FeedsButton   = FeedsMenu.Text("/feed")
	FeedsSelector = &tb.ReplyMarkup{}
	helpText      = "Please provide a *comma seperated* list of  rss feed urls i should scrape for you. \n" +
		"By default users are *subscribed* to top 100 crypto sites. \n" +
		"\nUsage: \n" +
		"/feed *add {url}* - add a new rss feed to your subscriptions \n" +
		"/feed *remove {url}* - remove a rss feed from your subscriptions \n" +
		"/feed *list* - list all subscribed feeds\n" +
		"/feed *reset* - reset to default feed list\n"
)

func feedsCommandHandler(bot *tb.Bot, db *storage.DB, analyzer *news.Analyzer) func(m *tb.Message) {
	return func(m *tb.Message) {
		if user, err := storage.UserRequired(m.Sender, db, bot); err == nil {
			if m.Text == "/feed" {
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
					wg := &sync.WaitGroup{}
					for _, feedUrl := range urlSlice {

						u, err := url.Parse(feedUrl)
						if err != nil {
							config.IgnoreErrorMultiReturn(bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not parse %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2))
							return
						}
						log.Print("adding ", u.String())
						log.Print("uri ", u.RequestURI())
						log.Print("path ", u.EscapedPath())
						err = analyzer.AddFeed(u, user, wg, false)
						wg.Wait()
						if err != nil {
							config.IgnoreErrorMultiReturn(bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not add feed %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2))
						}
					}
				case "remove":
					urls := s[1]
					urlSlice := strings.Split(urls, ",")
					for _, feedUrl := range urlSlice {
						log.Print("removing ", feedUrl)
						u, err := url.Parse(feedUrl)
						if err != nil {
							config.IgnoreErrorMultiReturn(bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not parse %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2))
							return
						}
						err = analyzer.RemoveFeed(u, user)
						if err != nil {
							config.IgnoreErrorMultiReturn(bot.Send(m.Sender, markdownEscape(fmt.Sprintf("could not remove feed %s\n%s", feedUrl, err.Error())), FeedsSelector, tb.ModeMarkdownV2))
						}
					}
				case "list":
					feeds := user.GetFeedsString(analyzer.Feeds)
					if len(feeds) > 0 {
						config.IgnoreErrorMultiReturn(
							bot.Send(user.User,
								markdownEscape(fmt.Sprintf("%s", strings.Join(unique(feeds), ", "))),
								tb.ModeMarkdownV2))
					}
				case "reset":
					for _, f := range user.GetFeedsString(analyzer.Feeds) {
						analyzer.Feeds[f].RemoveUser(user)
					}
					analyzer.AddUserToDefaultFeeds(user)
					err := db.Set(user)
					if err != nil {
						return
					}
					config.IgnoreErrorMultiReturn(bot.Send(user.User, "feed list set to default", tb.ModeMarkdownV2))
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

func initFeedsHandler(bot *tb.Bot, db *storage.DB, analyzer *news.Analyzer) {
	bot.Handle(&FeedsButton, feedsCommandHandler(bot, db, analyzer))
}
