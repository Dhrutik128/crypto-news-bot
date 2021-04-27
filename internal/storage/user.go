package storage

import (
	"encoding/json"
	"fmt"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type User struct {
	// the telegram user
	User *tb.User
	// custom user settings
	Settings UserSettings `json:"settings"`
	// todo -- registration timestamp
	Started time.Time
	// todo -- last message sent to user timestamp
	LastMessageSent time.Time
	// todo -- last message received from user timestamp
	LastMessageReceived time.Time
}

// used to store user settings
type UserSettings struct {
	Subscriptions map[string]bool `json:"subscriptions"`
	Feeds         []string        `json:"feeds"`
}

// check if user subscribed to feed url
func (s UserSettings) IsFeedSubscribed(feed string) bool {

	for _, userFeed := range s.Feeds {
		if userFeed == feed {
			return true
		}
	}
	return false
}

// add feed url to users feed subscription
func (u *User) AddFeed(feed string, db *bitcask.Bitcask) error {
	u.Settings.Feeds = append(u.Settings.Feeds, feed)
	return StoreUser(u, db)
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
func GetUser(u *tb.User, db *bitcask.Bitcask) (*User, error) {
	user := &User{User: u}
	userBytes, err := db.Get(user.Key())
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("user not found")
	}
	err = json.Unmarshal(userBytes, user)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed unmarshaling")
	}
	return user, nil
}

// store or update user
func StoreUser(user *User, db *bitcask.Bitcask) error {
	userByte, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return nil
	}
	return db.Put(user.Key(), userByte)
}

// checks if user is registered. If user is registered, function will return user.
func UserRequired(user *tb.User, db *bitcask.Bitcask, bot *tb.Bot) (*User, error) {
	u := User{User: user}
	if !db.Has(u.Key()) {
		bot.Send(user, "please run the command /start before using this bot")
		return nil, fmt.Errorf("user not found")
	}
	return GetUser(user, db)
}

func DeleteUser(user *tb.User, db *bitcask.Bitcask) error {
	u := User{User: user}
	return db.Delete(u.Key())
}
