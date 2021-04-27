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
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		now := time.Now()
		feed := &feeds.Feed{
			Title:       "BTC feed",
			Link:        &feeds.Link{Href: "http://localhost:8080/"},
			Description: "discussion about tech, footie, photos",
			Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		}

		feed.Items = []*feeds.Item{
			&feeds.Item{
				Title:       "Bitcoin is nice! now",
				Link:        &feeds.Link{Href: "https://cointelegraph.com/news/bitcoin-bulls-respond-with-a-150m-short-squeeze-above-53k-can-btc-go-higher"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "Bitcoin and XRP Dawg!!!Lsdasd1",
				Link:        &feeds.Link{Href: "https://cointelegraph.com/news/bitcoin-bulls-respond-with-a-150m-short-squeeze-above-53k-can-btc-go-higher"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "BTC generating money!",
				Link:        &feeds.Link{Href: "https://cointelegraph.com/news/bitcoin-bulls-respond-with-a-150m-short-squeeze-above-53k-can-btc-go-higher"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			}, &feeds.Item{
				Title:       "BTC performance increased by 25%!",
				Link:        &feeds.Link{Href: "https://cointelegraph.com/news/bitcoin-bulls-respond-with-a-150m-short-squeeze-above-53k-can-btc-go-higher"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "BTC now reaching new goal at moonstation!!!",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
				Description: "More thoughts on logicless templates",
				Created:     now,
			},
			&feeds.Item{
				Title:       "Why Bitcoin will fail in near future",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now,
			},
			&feeds.Item{
				Title:       "Why Bitcoin must fail in near future!!",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now,
			},
			&feeds.Item{
				Title:       "XRP reaching new all time high",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now.Add(-(time.Hour * 32)),
			},
			&feeds.Item{
				Title:       "XRP surging to new all time high",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now.Add(-(time.Hour * 32)),
			},
			&feeds.Item{
				Title:       "XLM get fined for being manipulative",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now,
			},
			&feeds.Item{
				Title:       "XLM next level",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now,
			},
			&feeds.Item{
				Title:       "XLM new ATH",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
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
			Description: "discussion about tech, footie, photos",
			Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		}

		feed.Items = []*feeds.Item{
			&feeds.Item{
				Title:       "!YES XRP so bad that SEC lawsuit is fucked up",
				Link:        &feeds.Link{Href: "https://cointelegraph.com/news"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "Ripple labs under investigation for fraud!!!!!!!!!!!",
				Link:        &feeds.Link{Href: "https://www.newsbtc.com/news/bitcoin/bitcoin-dominance-dropped-below-50-what-this-could-mean-for-the-market/"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "Ripple labs SEC investigation rollup!!!!",
				Link:        &feeds.Link{Href: "https://www.newsbtc.com/news/company/okex-insights-catallact-bitcoin-market-witness-the-growth-of-retail-participation-as-institutional-investors-continue-to-lead/"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "Ripple labs SEC investigation rollup!! act now",
				Link:        &feeds.Link{Href: "https://www.newsbtc.com/news/company/okex-insights-catallact-bitcoin-market-witness-the-growth-of-retail-participation-as-institutional-investors-continue-to-lead/"},
				Description: "A discussion on controlled parallelism in golang",
				Author:      &feeds.Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
				Created:     now,
			},
			&feeds.Item{
				Title:       "XLM get fined for being manipulative",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
				Created:     now,
			},&feeds.Item{
				Title:       "XLM to the top ",
				Link:        &feeds.Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
				Description: "How to use interfaces <em>effectively</em>",
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
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
