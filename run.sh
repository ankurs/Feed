#!/bin/bash -e
mkdir -p local-data/cassandra
mkdir -p local-data/redis
docker-compose up --build -d
docker logs -f feed
