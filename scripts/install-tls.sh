#!/usr/bin/env bash

apt-get install python3-certbot-nginx -y
certbot --nginx -d matomo.geogame.xyz

