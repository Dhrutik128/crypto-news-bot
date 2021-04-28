package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

func userBroadCastFunc(db *bitcask.Bitcask, cast news.BroadCast, bot *tb.Bot) func(key []byte) error {
	return func(key []byte) error {
		userBytes, err := db.Get(key)
		if err != nil {
			return err
		}
		user := storage.User{}
		err = json.Unmarshal(userBytes, &user)
		if err != nil {
			log.Println(err)
			return err
		}
		if user.Settings.Subscriptions[cast.Sentiment.Coin] && user.Settings.IsFeedSubscribed(cast.Sentiment.Feed) {
			cast.User = user.User
			err := sendBroadCast(bot, cast)
			if err != nil {
				checkUserBlockedBot(err, cast.User, db)
				return err
			}
			log.WithFields(log.Fields{"module": "[TELEGRAM]"}).Infof("BROADCAST \n%s\nto %s\n", cast.Sentiment.FeedItem.Title, cast.User.Username)
		}
		return nil
	}
}

var markdownEscapes = []string{"_", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

func markdownEscape(s string) string {
	for _, esc := range markdownEscapes {
		if strings.Contains(s, esc) {
			s = strings.Replace(s, esc, fmt.Sprintf("\\%s", esc), -1)
		}
	}
	return s
}

func sendBroadCast(bot *tb.Bot, b news.BroadCast) error {
	text := fmt.Sprintf("[*_Broadcasting latest %s News_*](%s)\n\n*Title:* %s\n*Published:* %s\n*Sentiment:* %s\n", b.Sentiment.Coin, markdownEscape(b.Sentiment.FeedItem.Link), markdownEscape(b.Sentiment.FeedItem.Title), markdownEscape(b.Sentiment.FeedItem.Published), markdownEscape(fmt.Sprintf("%f", b.Sentiment.Sentiment["compound"])))
	_, err := bot.Send(b.User, text, tb.ModeMarkdownV2)
	if err != nil {
		return err
	}
	return nil
}

func StartBroadCaster(b *news.Analyzer, bot *tb.Bot, broadCastChannel chan news.BroadCast) {
	broadcaster := func(broadCast news.BroadCast) {
		b.Db.Scan([]byte("user_"), userBroadCastFunc(b.Db, broadCast, bot))
	}
	//broadcaster()
	go func() {
		for {
			select {
			case broadCast := <-broadCastChannel:
				broadcaster(broadCast)
			}
		}
	}()
}
