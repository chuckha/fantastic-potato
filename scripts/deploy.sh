#$!/usr/bin/env bash

nerdctl pull ghcr.io/chuckha/geogame.xyz:latest
nerdctl rm -f geogame.xyz
nerdctl run --name geogame.xyz -d -p 8888:8888 ghcr.io/chuckha/geogame.xyz:latest