package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	subscriptionSelector   = &tb.ReplyMarkup{Selective: true}
	SubscriptionButtons    = make([]tb.Btn, 0)
	SubscriptionButtonsMap = make(map[string]tb.Btn, 0)

	SubscriptionMenu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
)

func initSubscriptionHandler(bot *tb.Bot, db *bitcask.Bitcask) {
	SubscriptionButtons, SubscriptionButtonsMap = getKeywordButtons("sub_", SubscriptionMenu)
	subscriptionSelector.Inline(ButtonWrapper(SubscriptionButtons, SubscriptionMenu)...)
	// ### Subscribe Handler ###
	bot.Handle(&btnSubscribe, func(m *tb.Message) {
		if user, err := storage.UserRequired(m.Sender, db, bot); err == nil {
			bot.Send(m.Sender, "manage your news subscriptions", getButtonsForUser(user))
		}
	})
	// ### Inline Keyboard Subscription Handler ###
	for _, btn := range SubscriptionButtonsMap {
		bot.Handle(&btn, func(c *tb.Callback) {
			if user, err := storage.UserRequired(c.Sender, db, bot); err == nil {
				user.ToggleSubscription(c.Data)
				log.WithFields(log.Fields{"module": "[TELEGRAM]", "coin": c.Data, "subscribed": user.Settings.Subscriptions[c.Data]}).Infof("updated subscription for %s", c.Sender.Username)
				err = storage.StoreUser(user, db)
				if err != nil {
					fmt.Println(err)
				}
				newKeyboard := getButtonsForUser(user).InlineKeyboard
				//c.Message.ReplyMarkup.InlineKeyboard = newKeyboard
				bot.EditReplyMarkup(c.Message, &tb.ReplyMarkup{InlineKeyboard: newKeyboard})
			}
		})
	}

}
