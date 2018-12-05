From golang:1.11

RUN apt-get update
RUN apt-get install supervisor -y

RUN mkdir -p /go/src/github.com/ankurs/Feed
RUN mkdir -p /opt/config/

COPY . /go/src/github.com/ankurs/Feed
COPY ./Feed/feed.toml /opt/config/

RUN go install github.com/ankurs/Feed/Feed/cmd/server

EXPOSE 9281 9282 9283 9284
