#!/bin/bash
docker build . -t feed|| exit
docker stop feed
docker rm feed
docker run --privileged --name feed -p 9281:9281 -p 9282:9282 -p 9283:9283 -p 9284:9284 -h `whoami`-`hostname` feed server
