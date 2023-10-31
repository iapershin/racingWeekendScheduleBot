FROM alpine:latest

ARG BOT_TOKEN


ENV BOT_TOKEN=${BOT_TOKEN}


RUN apk update && apk upgrade
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

WORKDIR /app

COPY ./config/config.yaml .
COPY ./cmd/race-weekend-bot/main .

CMD ["./main", "--config", "config.yaml"]
