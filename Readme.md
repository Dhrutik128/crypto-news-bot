# <p align="center"> Crypto News Telegram Bot

<p align="center"> A simple telegram bot that will help you stay updated on your latest crypto news 

* [Getting started.](#getting-started)
* [Usage.](#usage)
* [Deployment.](#deployment)
* [Planned.](#planned)



## Getting started
The latest version is always running at [@crypto-news-bot](https://t.me/crypstream_bot). You can subscribe to this bot or deploy your own bot using. [@BotFather](https://t.me/BotFather)

## Usage
Send ``/start`` to start the bot. By default, you rss feed list is set to the top 100 crypto rss feeds.
You will receive broadcasts from this feed list, once you subscribe to a coin. 
* ``/subscribe``  subscribe to news based on coins. The bot will then send you the latest news based on your subscription. 
* ``/feed`` manage your RSS feeds. The bot will fetch and broadcast news from your personal feeds, based on your coin subscription.
You can add and remove multiple feeds by providing a comma separated list of rss feed urls.
* ``/news`` get the latest news for any coin. personal feeds will also be included (currently for every user).
* ``/sentiments`` get an average sentiment from all news, grouped by coin. 

   
## Deployment 
If you want to deploy your own instance, create a new bot using [@BotFather](https://t.me/BotFather). 
1. Rename `config-example.yaml` to `config.yaml`
2. Paste your Telegram API Token to the `config.yaml`
3. Manage your default RSS feeds using `feeds.csv` 
4. Run `docker-compose up -d` to start the bot

## Planned 
* Huge refactor to improve performance - remove redundant code
* Making keywords configurable by user. Currently, news categorization by coin uses a static list of keywords.This list should be configurable.
* Improve news presentation for better readability.
* Improve sentiment analysis. 

## Contribution 
If you have found bugs or got some improvements, feel free to open a issue or pull request. 