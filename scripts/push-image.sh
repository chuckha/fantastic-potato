#!/usr/bin/env bash

docker build --platform linux/amd64 -t ghcr.io/chuckha/geogame.xyz:latest . 
docker push ghcr.io/chuckha/geogame.xyz:latest