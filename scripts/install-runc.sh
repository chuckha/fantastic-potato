#!/usr/bin/env bash

wget https://github.com/opencontainers/runc/releases/download/v1.0.3/runc.amd64
chmod +x runc.amd64
mv runc.amd64 /usr/local/bin/runc
