From golang:1.10

RUN go get github.com/derekparker/delve/cmd/dlv

RUN mkdir -p /go/src/github.com/ankurs/Feed
RUN mkdir -p /opt/config/

COPY . /go/src/github.com/ankurs/Feed
COPY ./Feed/Feed.toml /opt/config/

RUN go install github.com/ankurs/Feed/Feed/cmd/server

EXPOSE 9281 9282 9283 9284
