package telegram

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

// check if telegram api returns a user blocked
func checkUserBlockedBot(err error, user *tb.User, db *storage.DB) {
	switch err.(type) {
	case *tb.APIError:
		apiError := err.(*tb.APIError)
		if apiError.Code == 401 && strings.Contains(apiError.Description, "blocked") {
			log.WithFields(log.Fields{"module": "[TELEGRAM]"}).Infof("user %d blocked bot. Deleting user data.", user.ID)
			u := storage.User{User: user}
			err := db.Delete(&u)
			if err != nil {
				log.WithFields(log.Fields{"error": err.Error(), "module": "[TELEGRAM]"}).Error("could not delete user")
			}
		}
	}
}

func markdownEscape(s string) string {
	for _, esc := range markdownEscapes {
		if strings.Contains(s, esc) {
			s = strings.Replace(s, esc, fmt.Sprintf("\\%s", esc), -1)
		}
	}
	return s
}
