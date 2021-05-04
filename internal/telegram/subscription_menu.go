package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	subscriptionSelector   = &tb.ReplyMarkup{Selective: true}
	SubscriptionButtons    = make([]tb.Btn, 0)
	SubscriptionButtonsMap = make(map[string]tb.Btn, 0)

	SubscriptionMenu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
)

func initSubscriptionHandler(bot *tb.Bot, db *storage.DB) {
	SubscriptionButtons, SubscriptionButtonsMap = getKeywordButtons("sub_", SubscriptionMenu)
	subscriptionSelector.Inline(buttonWrapper(SubscriptionButtons, SubscriptionMenu, 4)...)
	// ### Subscribe Handler ###
	bot.Handle(&btnSubscribe, func(m *tb.Message) {
		if user, err := storage.UserRequired(m.Sender, db, bot); err == nil {
			config.IgnoreErrorMultiReturn(bot.Send(m.Sender, "manage your news subscriptions", getSubscriptionButtons(user)))
		}
	})
	// ### Inline Keyboard Subscription Handler ###
	for _, btn := range SubscriptionButtonsMap {
		bot.Handle(&btn, func(c *tb.Callback) {
			if user, err := storage.UserRequired(c.Sender, db, bot); err == nil {
				user.ToggleSubscription(c.Data)
				log.WithFields(log.Fields{"module": "[TELEGRAM]", "coin": c.Data, "subscribed": user.Settings.Subscriptions[c.Data]}).Infof("updated subscription for %s", c.Sender.Username)
				err = db.Set(user)
				if err != nil {
					fmt.Println(err)
				}
				newKeyboard := getSubscriptionButtons(user).InlineKeyboard
				config.IgnoreErrorMultiReturn(bot.EditReplyMarkup(c.Message, &tb.ReplyMarkup{InlineKeyboard: newKeyboard}))
			}
		})
	}

}
