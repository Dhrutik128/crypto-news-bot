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
	Subscriptions           map[string]bool `json:"subscriptions"`
	Feeds                   []string        `json:"feeds"`
	IsDefaultFeedSubscribed bool            `json:"is_default_feed_subscribed"`
}

// check if user subscribed to feed url
func (s UserSettings) IsFeedSubscribed(feed string) bool {
	// todo -- check if feed is one of default
	for _, userFeed := range s.Feeds {
		if userFeed == feed {
			return true
		}
	}
	return false
}

// add feed url to users feed subscription
func (u *User) AddFeed(feed string, db *DB) error {
	alreadyAdded := false
	for _, userFeed := range u.Settings.Feeds {
		if feed == userFeed {
			alreadyAdded = true
			break
		}
	}
	if !alreadyAdded {
		u.Settings.Feeds = append(u.Settings.Feeds, feed)
		return db.Set(u)
	}
	return fmt.Errorf("feed is already included in users feeds")
}
func (u *User) ToggleDefaultFeed(db *DB) error {
	u.Settings.IsDefaultFeedSubscribed = !u.Settings.IsDefaultFeedSubscribed
	return db.Set(u)
}

func (u *User) RemoveFeed(feed string, db *DB) error {
	removed := false
	for i, userFeed := range u.Settings.Feeds {
		if feed == userFeed {
			removed = true
			u.Settings.Feeds = remove(u.Settings.Feeds, i)
			break
		}
	}
	if removed {
		return db.Set(u)
	}
	return fmt.Errorf("no feed removed")
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
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
