package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	tb "gopkg.in/tucnak/telebot.v2"
)

var markdownEscapes = []string{"_", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

// StartUserBroadCaster that checks every incoming broadcast from broadCastChannel and sends them to users
// that are subscribed to the feed and coin
func StartUserBroadCaster(b *news.Analyzer, bot *tb.Bot, broadCastChannel chan news.BroadCast) {
	broadcaster := func(broadCast news.BroadCast) {
		err := b.Db.View(func(tx *buntdb.Tx) error {
			err := tx.Ascend("user", func(key, value string) bool {
				user := storage.User{}
				err := json.Unmarshal([]byte(value), &user)
				if err != nil {
					log.Println(err)
					return false
				}
				feed := b.Feeds[broadCast.FeedItem.Feed.String()]
				if feed != nil && feed.HasUser(user) {
					if user.Settings.Subscriptions[broadCast.FeedItem.Coin] {
						broadCast.User = user.User
						err := sendBroadCast(bot, broadCast)
						if err != nil {
							checkUserBlockedBot(err, broadCast.User, b.Db)
							return false
						}
						log.WithFields(log.Fields{"module": "[TELEGRAM]"}).Infof("BROADCAST \n%s\nto %s\n", broadCast.FeedItem.Item.Title, broadCast.User.Username)
						return true
					}
				}
				return true
			})
			return err
		})
		if err != nil {
			log.WithFields(log.Fields{"error": err.Error()}).Error("error while broadcasting")
		}
	}
	// start the broadcast channel
	go func() {
		for {
			select {
			case broadCast := <-broadCastChannel:
				broadcaster(broadCast)
			}
		}
	}()
}

// sendBroadCast sends the actual broadcast to the user
func sendBroadCast(bot *tb.Bot, b news.BroadCast) error {
	if b.User != nil {

		text := fmt.Sprintf("[*_Broadcasting latest %s News_*](%s)\n\n*Title:* %s\n*Published:* %s\n*Item:* %s\n",
			b.FeedItem.Coin,
			markdownEscape(b.FeedItem.Item.Link),
			markdownEscape(b.FeedItem.Item.Title),
			markdownEscape(b.FeedItem.Item.Published),
			markdownEscape(fmt.Sprintf("%f", b.FeedItem.Sentiment["compound"])))
		_, err := bot.Send(b.User, text, tb.ModeMarkdownV2)
		if err != nil {
			return err
		}

		return nil
	}
	return fmt.Errorf("broadcasting user is not set")
}
