#!/bin/bash
echo "generating proto"
protoc -I . Feed.proto --go_out=plugins=grpc:. --orion_out=.
