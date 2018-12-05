# Orion-Builder [![Build Status](https://travis-ci.com/carousell/Orion-Builder.svg?token=kSVweyyqayUyyfutjTqD&branch=master)](https://travis-ci.com/carousell/Orion-Builder)

## Setup Instructions
Orion-Builder is written in golang, please follow instructions on [https://golang.org/doc/install](https://golang.org/doc/install) to install, or you can also run
```
brew install golang
```
add the following lines to your `~/.profile`
```
export GOPATH="$HOME/code/go"
export GOBIN="$GOPATH/bin"
export PATH="$GOBIN:$PATH"
export PATH="$HOME/.gotools:$PATH"
```

source your `~/.profile`
```
source ~/.profile
```

then create the code dir
```
mkdir -p $GOPATH
```

we use `govendor` to vendor package in Orion-Builder, install it by running
```
go get -u github.com/kardianos/govendor
```
another helpful tool to check for unupdated packages is `Go-Package-Store`, install it by running
```
go get -u github.com/shurcooL/Go-Package-Store/cmd/Go-Package-Store
```
now clone this repo
```
mkdir -p $GOPATH/src/github.com/carousell/
git clone git@github.com:carousell/Orion-Builder.git $GOPATH/src/github.com/carousell/Orion-Builder
```

You need the following tools to better develop for go
```
go get -u github.com/golang/lint/golint
```

now you can build the package by using `make build`

## gRPC
for gRPC, you need to follow the following steps

get gRPC codebase
```
go get -u google.golang.org/grpc
```

install protobuf
```
brew install protobuf
```

install the protoc plugin for go
```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

install the protoc plugin for orion
```
go get -u github.com/carousell/Orion/protoc-gen-orion
```

## Development
Please install docker-compose from https://docs.docker.com/compose/ and execute `./run.sh`

### Links
Once docker-compose is started
* Hystix Dashboard is available on [http://192.168.99.100:9001/monitor/monitor.html?streams=%5B%7B%22name%22%3A%22%22%2C%22stream%22%3A%22http%3A%2F%2Fibus%3A9283%2Fhystrix.stream%22%2C%22auth%22%3A%22%22%2C%22delay%22%3A%22%22%7D%5D](http://192.168.99.100:9001/monitor/monitor.html?streams=%5B%7B%22name%22%3A%22%22%2C%22stream%22%3A%22http%3A%2F%2Fibus%3A9283%2Fhystrix.stream%22%2C%22auth%22%3A%22%22%2C%22delay%22%3A%22%22%7D%5D)
* Metrics are available on [http://192.168.99.100:9284/metrics](http://192.168.99.100:9284/metrics)
* Pprof is available on [http://192.168.99.100:9284/debug/pprof/](http://192.168.99.100:9284/debug/pprof/)
    * pprof documentation [https://golang.org/pkg/net/http/pprof/](https://golang.org/pkg/net/http/pprof/)
* Zipkin spans are available on [http://192.168.99.100:9411/zipkin/](http://192.168.99.100:9411/zipkin/)
