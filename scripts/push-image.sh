#!/usr/bin/env bash

docker build . -t ghcr.io/chuckha/geogame.xyz:latest
docker push ghcr.io/chuckha/geogame.xyz:latest