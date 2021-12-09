#!/usr/bin/env bash

wget https://github.com/containernetworking/plugins/releases/download/v1.0.1/cni-plugins-linux-amd64-v1.0.1.tgz
tar -C /opt/cni/bin -xzf cni-plugins-linux-amd64-v1.0.1.tgz
rm cni-plugins-linux-amd64-v1.0.1.tgz