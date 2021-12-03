I have this running on a basic digital ocean box

The docker images are pushed to github's container registry (ghcr.io).
The server uses `nerdctl` to interface with `containerd` which interfaces with `runc`.

Configuring auth for containerd turned out to not work, but running `ctr pull <img> -u 'chuckha:$PAT'> worked fine.


# installing `runc`

get the latest release and install it into /usr/local/bin

# installing containerd

grab the latest release here https://containerd.io/downloads/

install all of the binaries into /usr/local/bin
run a command to generate the default config into /etc/containerd/something...

# installing nerdctl

download the binary

# installing cni

download the tarball and extract into /opt/cni/bin