#!/usr/bin/env bash

wget https://github.com/containerd/containerd/releases/download/v1.6.0-beta.3/containerd-1.6.0-beta.3-linux-amd64.tar.gz
tar -xzf containerd-1.6.0-beta.3-linux-amd64.tar.gz 
mv bin/* /usr/local/bin
rm -rf bin

wget https://github.com/containerd/containerd/archive/v1.6.0-beta.3.zip
unzip v1.6.0-beta.3.zip 
mv ./containerd
mv containerd-1.6.0-beta.3/containerd.service /etc/systemd/system/.
systemctl start containerd