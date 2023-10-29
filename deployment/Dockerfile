FROM alpine:latest

ARG BOT_TOKEN


ENV BOT_TOKEN=${BOT_TOKEN}


#INSTALL CA-CERTS
RUN apk update && apk add ca-certificates
RUN update-ca-certificates
#COPY BINARY
COPY ./cmd/main .
CMD ["./main"]