package main

import (
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	StartFeedMocker()
}
func StartFeedMocker() {
	r := mux.NewRouter()
	r.HandleFunc("/btc", func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "BTC feed",
			Link:        &feeds.Link{Href: "http://localhost:8080/btc"},
			Description: "discussion about Bitcoin",
			Author:      &feeds.Author{Name: "Author", Email: "author@news.net"},
			Created:     now,
		}

		feed.Items = []*feeds.Item{
			{
				Title:       "Bitcoin is very nice!",
				Link:        &feeds.Link{Href: "https://cointelegraph.com/news/bitcoin-bulls-respond-with-a-150m-short-squeeze-above-53k-can-btc-go-higher"},
				Description: "A discussion on bitcoin",
				Author:      &feeds.Author{Name: "My Name", Email: "author@news.net"},
				Created:     now,
			},
		}

		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		writer.Write([]byte(rss))

	})

	r.HandleFunc("/xrp", func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "XRP feed",
			Link:        &feeds.Link{Href: "http://localhost:8080/xrp"},
			Description: "discussion about tech and ripple",
			Author:      &feeds.Author{Name: "Some Author", Email: "author@somemail.net"},
			Created:     now,
		}
		feed.Items = []*feeds.Item{
			{
				Title:       "Giant Video Game provider Nexon Buys $100M Worth of XRP",
				Link:        &feeds.Link{Href: "https://www.livebitcoinnews.com/bitcoin-price-analysis-btc-faces-major-hurdle-dips-limited/"},
				Description: "A discussion on xrp",
				Author:      &feeds.Author{Name: "Some Author", Email: "author@nexon.com"},
				Created:     now,
			},
		}
		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		writer.Write([]byte(rss))

	})
	r.HandleFunc("/news", func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Crypto feed",
			Link:        &feeds.Link{Href: "http://localhost:8080/news"},
			Description: "discussion about crypto",
			Author:      &feeds.Author{Name: "Crypto Maniac", Email: "crypto@maniac.net"},
			Created:     now,
		}
		feed.Items = []*feeds.Item{
			{
				Title:       "$DOGE up 600% this month. ",
				Link:        &feeds.Link{Href: "https://cryptopotato.com/first-1200-us-stimulus-check-put-in-dogecoin-worth-over-400000-now/"},
				Description: "A discussion on another shitcoin",
				Author:      &feeds.Author{Name: "Some Author", Email: "author@nexon.com"},
				Created:     now,
			},
		}
		rss, err := feed.ToRss()
		if err != nil {
			log.Fatal(err)
		}
		writer.Write([]byte(rss))

	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
