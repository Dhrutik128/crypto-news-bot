package main

import (
	"fmt"
	"github.com/gohumble/crypto-news-bot/internal/news"
	"github.com/gohumble/crypto-news-bot/internal/telegram"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/prologic/bitcask"
	log "github.com/sirupsen/logrus"
	"io"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type NewsBot struct {
	NewsFeed *news.Analyzer
	Telegram *tb.Bot
	Db       *bitcask.Bitcask
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
	db, err := bitcask.Open("data")
	if err != nil {
		panic(err)
	}
	bot := NewsBot{NewsFeed: news.NewAnalyzer(db, Config.RefreshRate), Db: db}
	bot.Telegram = telegram.New(bot.Db, bot.NewsFeed, Config.BotToken)
	bot.Start()
}

func (b *NewsBot) Start() {
	go b.NewsFeed.Start()
	b.Telegram.Start()
}
