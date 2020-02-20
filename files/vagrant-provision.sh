#!/usr/bin/env bash

set -euo pipefail

## Load settings if the file exists inside the VM.
envrc=/home/vagrant/.kusanagi-vagrant/.envrc
if [ -e ${envrc} ]; then
  . ${envrc}
fi

function main() {
  disable_swap
  disable_systemd_resolvd
  ensure_netfilter
  disable_ipv4_on_docker0
  ensure_ip4_localhost_in_etc_hosts
  configure_docker
}


function disable_swap() {
  ## Disable swap (required for kubelet)
  sudo sed -i '/swap/d' /etc/fstab
  sudo swapoff -a
}

function disable_systemd_resolvd() {
  ## Disable systemd-resolved and install resolv.conf.
  ##
  ## TODO: Try this instead: https://unix.stackexchange.com/a/358485
  sudo systemctl disable systemd-resolved.service 2>&1
  sudo systemctl stop systemd-resolved.service
  sudo systemctl mask systemd-resolved.service

  cat <<EOF | sudo tee /etc/resolv.conf
# When running on IPV4-only networks, these IPV6 nameservers
#   break connectivity.
# nameserver 2001:4860:4860::8888
# nameserver 2001:4860:4860::8844
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF
}

function ensure_netfilter() {
  ## Ensure br_netfilter kernel module
  ## is loaded on every reboot.
  echo br_netfilter | sudo tee /etc/modules-load.d/br_netfilter.conf
  sudo systemctl daemon-reload
  sudo systemctl restart systemd-modules-load.service
  lsmod | grep br_netfilter

  ## Confirm bridge-nf-call-ip(6)tables proc values.
  ## ref: https://github.com/corneliusweig/kubernetes-lxd
  grep '^1$' /proc/sys/net/bridge/bridge-nf-call-iptables
  grep '^1$' /proc/sys/net/bridge/bridge-nf-call-ip6tables
}

function disable_ipv4_on_docker0() {
  ## Forcefully disable IPv4 on docker0
  ## now and at every boot.
  cat <<EOF | sudo tee /etc/systemd/system/remove-docker0-ipv4.service
[Unit]
Description=Remove IPv4 address from docker0 interface.
After=docker.service network.target
Requires=docker.service network.target
Before=kubelet.service
Wants=kubelet.service
[Service]
ExecStart=/sbin/ip addr del 172.17.0.1/16 dev docker0
[Install]
WantedBy=multi-user.target
EOF
  sudo systemctl daemon-reload
  sudo systemctl enable remove-docker0-ipv4.service
  sudo systemctl start remove-docker0-ipv4.service
  sleep 1
  ip -4 addr show dev docker0
}

function ensure_ip4_localhost_in_etc_hosts() {
  sudo sed -i 's@^127.0.0.1\tlocalhost$@127.0.0.1\tlocalhost ip4-localhost@g' /etc/hosts
}

function configure_docker() {
  ## Configure Docker with an option to enable IPv6.
  cat <<EOF | sudo tee /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "ipv6": false,
  "fixed-cidr-v6": "2001:db8:1::/64"
}
EOF
  if [ -f /etc/docker/config.json ]; then
    # The virtualbox image has this file, but the hyperv image does not.
    #   Apparently, the are minor differences in the origin bento/ubuntu-18.04
    #   images for the two virt platforms.  Permissions need to be corrected
    #   to avoid a warning message upon each invocation of "docker" cli.
    sudo chmod 644 /etc/docker/config.json
  fi
  sudo mkdir -p /etc/systemd/system/docker.service.d
  cat <<EOF | sudo tee /etc/systemd/system/docker.service.d/override.conf
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H unix:// -D
EOF
  sudo systemctl daemon-reload
  sudo systemctl restart docker.service
}

#####################################################################
# Run the main program
main "$@"

## Report general success.
echo OK
