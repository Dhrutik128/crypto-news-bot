package storage

import (
	"encoding/json"
	"github.com/gohumble/crypto-news-bot/internal/config"
	"github.com/tidwall/buntdb"
	"time"
)

// Storable items must provide a function to retrieve the database key
type Storable interface {
	Key() []byte
}

type DB struct {
	*buntdb.DB
}

// GetFeedLastDownloadTime from storage
func (db DB) GetFeedLastDownloadTime() time.Time {
	var t time.Time
	config.IgnoreError(db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("lastDownloadTime")
		if err != nil {
			return err
		}
		t, err = time.Parse(time.RFC3339, val)
		return err
	}))
	return t
}

// SetFeedLastDownloadTime in storage
func (db DB) SetFeedLastDownloadTime(t time.Time) {
	config.IgnoreError(db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("lastDownloadTime", t.Format(time.RFC3339), nil)
		return err
	}))
}

// Exists checks is storable item exists
func (db DB) Exists(storable Storable) (ok bool, err error) {
	ok = false
	err = db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(string(storable.Key()))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err == buntdb.ErrNotFound {
			return
		}
		return
	}
	ok = true
	return

}

// Get a storable item
func (db DB) Get(object Storable) error {
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(string(object.Key()))
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(val), object)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Set a storable item.
func (db DB) Set(object Storable) error {
	err := db.Update(func(tx *buntdb.Tx) error {
		b, err := json.Marshal(object)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(string(object.Key()), string(b), nil)

		return err
	})
	return err
}

// Delete a storable item.
func (db DB) Delete(object Storable) error {
	err := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(string(object.Key()))
		return err
	})
	return err
}
