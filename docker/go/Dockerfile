FROM golang:1.15-alpine as builder


RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN echo -e  "\e[91mSTART BUILDING SERVICE ON \c" && go version && echo -e "\e[0m"

RUN go build -o crypto-news-bot .

#building finished. Now extracting single bin in second stage.
FROM alpine

RUN adduser -S -D -H -h /app crypto-news-bot

USER crypto-news-bot

COPY --from=builder /build/crypto-news-bot /app/
COPY --from=builder /build/config.yaml /app/config.yaml
COPY --from=builder /build/vader_lexicon.txt /app/vader_lexicon.txt
COPY --from=builder /build/emoji_utf8_lexicon.txt /app/emoji_utf8_lexicon.txt
COPY --from=builder /build/feeds.csv /app/feeds.csv


WORKDIR /app

CMD ["./crypto-news-bot"]
