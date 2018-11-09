Kubernetes Development Kit (KDK)
===

## Quickstart (`TL;DR`)

This Quickstart assumes that you have installed all of the
[dependencies](https://github.com/cisco-sso/kdk#installation-instructions).


### Mac and Linux

```console
curl -sSL https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install | bash
kdk init && kdk ssh
```

### Windows

Please use Windows10 powershell for installation.

```console
Set-ExecutionPolicy Bypass -Scope Process -Force
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cisco-sso/kdk/master/files/install.ps1'))
kdk init ; kdk ssh
```

NOTE: After installation, Windows CMD prompt will work. The KDK has not been
tested with Cygwin, Mingw, or Windows Subsystem for Linux.

## Installation Instructions

Detailed installation instructions of the KDK along with all of its dependencies are found here:

* [Mac](https://kdf.csco.cloud/getting-started/mac/)
* [Windows](https://kdf.csco.cloud/getting-started/windows/)
* [Linux](https://kdf.csco.cloud/getting-started/linux/)

## Background

The kdk repository may be used to create a docker container with all of the
tools that one would typically use in order to develop and operate kubernetes
clusters.

Getting setup to create and operate a Kubernetes cluster in AWS, Openstack, or
even locally may be painful because a user may be running Windows10 or OSX, and
one must configure 20+ tools for cluster automation to work effectively. We've
created a Docker Image to enable every one of us to work in the same
environment, with the same tools, at the same versions.

Tools include: docker, kubectl, helm, multihelm, kops, terraform, ansible,
minio-cli, aws-cli, direnv, golang, git, vi/vim, emacs, python 2/3, jq, zsh,
helm-s3, kafkacat, dig, ssh-keygen, gitslave, dep, gomplate, minikube, awscli,
docker-compose, neutronclient, openstackclient, supernova, virtualenv, yq,
colordiff, nmap, screen, tmux, yadm, and many others.

* Some example use cases include:
  * Operating Kubernetes clusters.
  * Deploying Kubernetes clusters to AWS using `kops`.
  * Developing and applying Helm Charts and mh Apps.
  * Developing docker containers.


## Basic Usage

1. Create or re-create the config

```console
kdk init
```

2. Connect or reconnect to the KDK (will pull and start container if necessary)

```console
kdk ssh
```

3. Destroy the KDK

```console
kdk destroy
```

4. Update the KDK (binary, config, and container)

```console
kdk update
```

## Running Multiple KDK Containers

You might have a need to run multiple KDK containers.  The KDK CLI can do that!

1. Create a new KDK config

  - **NOTE:** name parameter must be unique (no other container can have this name)
```console
kdk init --name kdk1
```

2. Connect to `kdk1` container

```console
kdk ssh --name kdk1
```

**NOTE:** There are many configuration options available in `kdk init`.See `kdk init --help` for details 
