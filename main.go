package main

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/storage"
	"github.com/gohumble/crypto-news-bot/internal/telegram"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"io"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type NewsBot struct {
	NewsFeed *news.Analyzer
	Telegram *tb.Bot
	Db       *storage.DB
}

func initLogger() {
	path := "log/crypto-news-bot"
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", path, "%Y-%m-%d.%H:%M:%S"),
		rotatelogs.WithMaxAge(time.Hour*10),
		rotatelogs.WithRotationTime(time.Second*10),
	)
	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, writer))
	log.SetFormatter(&log.JSONFormatter{})
	return
}
func main() {
	initLogger()
	db, err := buntdb.Open("./data/data.db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.CreateIndex("user", "user_*", buntdb.IndexJSON("user.id"))
	if err != nil {
		panic(err)
	}
	err = db.CreateIndex("feed", "feed_*", buntdb.IndexJSON("subscribers"))
	if err != nil {
		panic(err)
	}
	//err = db.CreateIndex("user_feed", "user_*", buntdb.IndexJSON("user.id"))
	err = db.CreateIndex("item", "item_*", buntdb.IndexJSON("feed_item.publishedParsed"))
	if err != nil {
		panic(err)
	}
	log.Infoln("started database")
	database := &storage.DB{DB: db}
	bot := NewsBot{NewsFeed: news.NewAnalyzer(database, Config.RefreshPeriodDuration, Config.NewsStorageDuration), Db: database}
	log.Infoln("initialized bot")
	bot.Telegram = telegram.New(bot.Db, bot.NewsFeed, Config.BotToken)
	bot.Start()
}

func (b *NewsBot) Start() {
	log.Infoln("starting news feed")
	b.NewsFeed.Start()
	log.Infoln("starting telegram")
	b.Telegram.Start()
}
