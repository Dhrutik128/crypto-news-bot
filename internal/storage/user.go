package storage

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type User struct {
	// the telegram user
	User *tb.User `json:"user"`
	// custom user settings
	Settings UserSettings `json:"settings"`
	// todo -- registration timestamp
	Started time.Time `json:"started"`
	// todo -- last message sent to user timestamp
	LastMessageSent time.Time `json:"last_message_sent"`
	// todo -- last message received from user timestamp
	LastMessageReceived time.Time `json:"last_message_received"`
}

// used to store user settings
type UserSettings struct {
	Subscriptions map[string]bool `json:"subscriptions"`
}

func (u User) GetFeeds(analyzerFeeds map[string]*Feed) []*Feed {
	feeds := make([]*Feed, 0)
	for _, feed := range analyzerFeeds {
		for _, subscriber := range feed.Subscribers {
			if subscriber == u.User.ID {
				feeds = append(feeds, feed)
				break
			}
		}
	}
	return feeds
}
func (u User) GetFeedsString(analyzerFeeds map[string]*Feed) []string {
	feeds := make([]string, 0)
	for f, feed := range analyzerFeeds {
		for _, subscriber := range feed.Subscribers {
			if subscriber == u.User.ID {
				feeds = append(feeds, f)
				break
			}
		}
	}
	return feeds
}

// generate users database key
func (u User) Key() []byte {
	return []byte(fmt.Sprintf("user_%d", u.User.ID))
}

// toggle coin subscription
func (u *User) ToggleSubscription(subscription string) {
	u.Settings.Subscriptions[subscription] = !u.Settings.Subscriptions[subscription]
}

// add coin subscription
func (u *User) AddSubscription(subscription string) {
	u.Settings.Subscriptions[subscription] = true
}

// remove coin subscription
func (u *User) RemoveSubscription(subscription string) {
	u.Settings.Subscriptions[subscription] = false
}

// get user based ob telegram user
func GetUser(u *tb.User, db *DB) (*User, error) {
	user := &User{User: u}
	err := db.Get(user)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// checks if user is registered. If user is registered, function will return user.
func UserRequired(user *tb.User, db *DB, bot *tb.Bot) (*User, error) {
	u := User{User: user}
	if ok, _ := db.Exists(u); !ok {
		bot.Send(user, "please run the command /start before using this bot")
		return nil, fmt.Errorf("user not found")
	}
	return GetUser(user, db)
}

func DeleteUser(user *tb.User, db *DB) error {
	u := User{User: user}
	return db.Delete(u)
}
