#!/usr/bin/env bash
set -euo pipefail

# This script configures `dnsmasq` so that Kubernetes resources within the KDK
#   are accessible via a socks proxy, thus enabling a web browser on the host
#   to access kubernetes ingress resources via real DNS names.
#
# When running the KDK docker container in conjunction with a Docker-for-Mac or
#   Docker-for-Windows kubernetes cluster on a local host machine, it is
#   difficult to access via a host web browser the urls of the ingress
#   resources from the host machine.  Doing so typically requires modifying the
#   host machine's /etc/hosts file to match the DNS names of services served by
#   the ingress by pointing locally to the host machine where the ingress ports
#   are exposed.
#
# We solve this problem differently, and without modifying host machine network
#   configurations.  This makes the KDK more portable across Mac, Windows, and
#   Linux.
#
# The approach is:
#
# * Run DNS mask within the KDK to resolve the wildcard domain served by the
#   ingress resource '*.docker-for-desktop.example.org'.
#
# * The wildcard domain '*.docker-for-desktop.example.org' within the KDK will
#   resolve to the address pointed to by special name 'host.docker.internal',
#   which is the host machine where the ingress container is actually running
#   and serving on port 80 and 443.  The KDK docker container may access the
#   host from this IP.
#
# * Upon "kdk ssh", configure the SSH command line to open a local socks proxy
#   port that is run through the KDK docker container.  This uses the "-D
#   <port-number>" option.
#
# * Configure a web browser on the host machine to use the local SOCKS port for
#   all requests.
#
# * Configure a host web browser proxy settings to use localhost:<port-number>
#   and forward all DNS through socks5 tunnel.
#     If using OSX, be sure to create a new Firefox profile with:
#       /Applications/Firefox.app/Contents/MacOS/firefox-bin -P
#     If using Windows , be sure to create a new Firefox profile with:
#       firefox.exe -P
#     Availale ingress services may be found with this command
#       kubectl get ing |grep -v HOSTS | awk '{print "  http://"$2}'


# Repair a potentialy broken /etc/resolv.conf on Windows.  Docker for Windows
#   will produce a broken /etc/resolv.conf if the machine does not have a
#   search domain set.  The resolv.conf will have a line that looks like
#   `search `, which is corrupt because it is missing an argument to search
#   such as `search domain.com`.  This will prevent tools like `nslookup` and
#   `host` from working.  Detect and fix this situation by removing the
#   offending line.  Note, 'sed -i' doesn't work because the file is a mount.
sed '/^search $/d' /etc/resolv.conf > /tmp/resolv.conf
cat /tmp/resolv.conf | tee /etc/resolv.conf && rm -f /tmp/resolv.conf

# Find the host IP from the perspective of *this* container
export HOST_ACCESS_IP=$(host host.docker.internal | grep "has address" | cut -d' ' -f 4)
# List of wildcard domains, space separated
export DOMAINS="kdk kube docker docker-for-desktop docker-for-desktop.example.org"

if [ -z "$HOST_ACCESS_IP" ]; then
    echo "Unable to find IP of ingress controller service"
    exit 1
fi

# Configure Dnsmasq to forward the docker-for-desktop domain
echo 'listen-address=127.0.0.1' | tee /etc/dnsmasq.d/docker-for-desktop
for DOMAIN in $DOMAINS; do
  echo "address=/${DOMAIN}/${HOST_ACCESS_IP}" | tee -a /etc/dnsmasq.d/docker-for-desktop
done

# Rewrite /etc/resolv.conf to use dnsmasq first
if [[ ! -f /etc/resolv.conf.bak ]] && ! (grep 'nameserver 127.0.0.1' /etc/resolv.conf &>/dev/null); then
    cp -a /etc/resolv.conf /etc/resolv.conf.bak
    echo 'nameserver 127.0.0.1' | tee /etc/resolv.conf
    cat /etc/resolv.conf.bak | tee -a /etc/resolv.conf
fi
