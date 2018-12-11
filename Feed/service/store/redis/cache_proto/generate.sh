#!/bin/bash
echo "generating cache proto"
protoc -I . cache.proto --go_out=.
