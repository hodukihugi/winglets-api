FROM golang:1.22.1-alpine

# Required because go requires gcc to build
RUN apk add build-base
RUN apk add inotify-tools
RUN apk add git
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN echo $GOPATH

COPY . /winglets_web
WORKDIR /winglets_web

RUN go mod download

CMD /winglets_web/docker/run.sh