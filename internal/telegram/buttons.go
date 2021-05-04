package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

// buttonWrapper wrap buttons slice in rows of length i
func buttonWrapper(buttons []tb.Btn, markup *tb.ReplyMarkup, length int) []tb.Row {
	buttonLength := len(buttons)
	rows := make([]tb.Row, 0)

	if buttonLength > length {
		for i := 0; i < buttonLength; i = i + length {
			buttonRow := make([]tb.Btn, 4)
			if i+length < buttonLength {
				buttonRow = buttons[i : i+4]
			} else {
				buttonRow = buttons[i:]
			}
			rows = append(rows, markup.Row(buttonRow...))
		}
		return rows
	}
	rows = append(rows, markup.Row(buttons...))
	return rows
}

// getKeywordButtons returns buttons for every word defined in news.KeyWords
func getKeywordButtons(uniquePrefix string, menu *tb.ReplyMarkup) ([]tb.Btn, map[string]tb.Btn) {
	buttons := make([]tb.Btn, 0)
	buttonMap := make(map[string]tb.Btn, 0)
	for _, keyWordItem := range news.KeyWords {
		item := keyWordItem[0]
		buttonMap[item] = menu.Data(item, fmt.Sprintf("%s%s", uniquePrefix, item), item)
		buttons = append(buttons, buttonMap[item])
	}
	return buttons, buttonMap
}

// getSubscriptionButtons based on the user coin subscriptions.
// coin is prefixed with ✅ if its included in user subscriptions
func getSubscriptionButtons(user *storage.User) *tb.ReplyMarkup {
	var subButtonSlice = make([]tb.Btn, 0)
	var subButtonSelector = &tb.ReplyMarkup{Selective: true}
	var menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	for _, words := range news.KeyWords {
		unique := fmt.Sprintf("%s%s", "sub_", words[0])
		text := SubscriptionButtonsMap[words[0]].Text
		if user.Settings.Subscriptions[words[0]] {
			subButtonSlice = append(subButtonSlice, menu.Data(fmt.Sprintf("%s %s", "✅", text), unique, text))
		} else {
			subButtonSlice = append(subButtonSlice, menu.Data(text, unique, text))
		}
	}
	subButtonSelector.Inline(buttonWrapper(subButtonSlice, menu, 4)...)
	return subButtonSelector
}

// FEED BUTTONS ARE CURRENTLY NOT IN USE
/*
// getButtons from slice of strings
func getButtons(uniquePrefix string, items []string, menu *tb.ReplyMarkup) ([]tb.Btn, map[string]tb.Btn) {
	buttons := make([]tb.Btn, 0)
	buttonMap := make(map[string]tb.Btn, 0)
	for _, item := range items {
		button := menu.Data(item, fmt.Sprintf("%s%s", uniquePrefix, item), item)
		buttonMap[item] = button
		buttons = append(buttons, buttonMap[item])
	}
	return buttons, buttonMap
}

func getDefaultFeedButtons(uniquePrefix string, items []string, menu *tb.ReplyMarkup, user *storage.User) ([]tb.Btn, map[string]tb.Btn) {
	buttons, buttonsMap := getButtons(uniquePrefix, items, menu)
	for i, button := range buttons {
		if button.Data == "top100" {
			if user != nil {
				if user.Settings.IsDefaultFeedSubscribed {
					text := fmt.Sprintf("%s %s", "✅", "top100")
					prefix := fmt.Sprintf("%s%s", uniquePrefix, "top100")
					buttons[i] = menu.Data(text, prefix, "top100")
				}
			}
		}
	}
	menu.Inline(buttonWrapper(buttons, menu)...)
	return buttons, buttonsMap
}
*/
