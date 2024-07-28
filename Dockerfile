FROM golang:1.22-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/kish1n/usdt_listening
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/usdt_listening /go/src/github.com/kish1n/usdt_listening


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/usdt_listening /usr/local/bin/usdt_listening
COPY config.yaml /usr/local/bin/config.yaml
COPY contractABI.json /usr/local/bin/contractABI.json
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["usdt_listening"]
